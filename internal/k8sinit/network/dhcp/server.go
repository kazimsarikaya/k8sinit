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
)

type DhcpConf struct {
	LeasesFile      string
	Interface       string
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
	start = (start >> 4) << 4
	size, bits := addrs[0].Mask.Size()
	rembits := bits - size
	max := (1 << (rembits + 1)) - 11

	confbase := fmt.Sprintf("/%v/config/dhcpd.%v.", poolName, ifname)
	conf := DhcpConf{
		LeasesFile:      confbase + "leases.json",
		Interface:       ifname,
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
	laddr := net.UDPAddr{
		IP:   lip,
		Port: dhcpv4.ServerPort,
	}
	server, err := server4.NewServer(ifname, &laddr, s.handler)
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
			s.server.Serve()
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
	klog.V(0).Infof("dhcp packet: %v", m)
}
