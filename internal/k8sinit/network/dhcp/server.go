/*
Copyright 2020 Kazım SARIKAYA

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dhcp

import (
	"fmt"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/network"
	"github.com/pkg/errors"
	klog "k8s.io/klog/v2"
	"net"
	"sync"
	"time"
)

type DhcpConf struct {
	LeasesFile      string
	Interface       string
	ServerIP        net.IP
	Mask            net.IPMask
	PoolStart       uint32
	MaxHost         int
	StaticHostsFile string
}

type NonBlockingDhcpServer struct {
	conf    DhcpConf
	wg      *sync.WaitGroup
	server  *server4.Server
	started bool
}

func NewNonBlockingDhcpSever(poolName, ifname string) (*NonBlockingDhcpServer, error) {
	if poolName == "" {
		return nil, k8sinit.K8SInitNotInstalledError
	}
	addrs, err := network.ListIpAddressesOfIfname(ifname)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, errors.New("cannot find server ip address")
	}
	lip := addrs[0].IP
	start := network.Ip2Int(lip)
	start = ((start >> 4) << 4) + 10
	size, bits := addrs[0].Mask.Size()
	rembits := bits - size
	max := (1 << (rembits + 1)) - 11

	confbase := fmt.Sprintf("/%v/config/dhcpd.%v.", poolName, ifname)
	conf := DhcpConf{
		LeasesFile:      confbase + "leases.json",
		Interface:       ifname,
		ServerIP:        lip,
		StaticHostsFile: confbase + "static.json",
		Mask:            addrs[0].Mask,
		PoolStart:       start,
		MaxHost:         max,
	}

	var wg sync.WaitGroup
	s := &NonBlockingDhcpServer{
		conf:    conf,
		wg:      &wg,
		started: false,
	}
	server, err := server4.NewServer(ifname, nil, s.handler)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create dhcp server")
	}
	s.server = server
	return s, nil
}

func (s *NonBlockingDhcpServer) Start() {
	s.wg.Add(1)
	go func() {
		s.started = true
		for s.started {
			klog.V(0).Infof("dhcpd will be started")
			err := s.server.Serve()
			klog.V(0).Error(err, "dhcpd stopped it will be restarted")
		}
		s.wg.Done()
	}()
}

func (s *NonBlockingDhcpServer) Stop() {
	if s.started {
		s.started = false
		s.server.Close()
	}
	s.wg.Wait()
}

func (s *NonBlockingDhcpServer) Wait() {
	s.wg.Wait()
}

func (s *NonBlockingDhcpServer) handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	if m == nil {
		return
	}
	if m.OpCode != dhcpv4.OpcodeBootRequest {
		return
	}
	klog.V(0).Infof("dhcp packet: %v, user-class: %v", m, m.UserClass())
	reply, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		klog.V(0).Error(err, "cannot create dhcp reply")
		return
	}

	reply.UpdateOption(dhcpv4.OptServerIdentifier(s.conf.ServerIP))
	reply.UpdateOption(dhcpv4.OptSubnetMask(s.conf.Mask))
	reply.UpdateOption(dhcpv4.OptDNS(s.conf.ServerIP))
	reply.UpdateOption(dhcpv4.OptRouter(s.conf.ServerIP))
	reply.UpdateOption(dhcpv4.OptNTPServers(s.conf.ServerIP))
	reply.UpdateOption(dhcpv4.OptIPAddressLeaseTime(time.Minute * 30))

	uses := m.UserClass()
	if len(uses) == 1 && uses[0] == "iPXE" {
		reply.UpdateOption(dhcpv4.OptBootFileName(fmt.Sprintf("http://%v:8000/api/network/tftp", s.conf.ServerIP)))
	} else {
		reply.UpdateOption(dhcpv4.OptBootFileName(k8sinit.UndiFilename))
	}
	reply.UpdateOption(dhcpv4.OptTFTPServerName(s.conf.ServerIP.String()))
	reply.YourIPAddr = network.Int2Ip(s.conf.PoolStart)
	reply.ServerIPAddr = s.conf.ServerIP

	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	case dhcpv4.MessageTypeRequest:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	default:
		klog.V(0).Error(errors.New("unknown dhcp mt"), "cannot select dhcp mt")
		return
	}
	klog.V(0).Infof("dhcp packet: %v, user-class: %v", reply, reply.UserClass())
	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		klog.V(0).Error(err, "cannot send dhcp reply")
	}
}
