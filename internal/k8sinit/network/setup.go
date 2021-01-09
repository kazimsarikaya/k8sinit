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

package network

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	klog "k8s.io/klog/v2"
	"os/exec"
)

func GetInterfaces() ([]string, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get ip links")
	}
	var result []string
	for _, link := range links {
		result = append(result, link.Attrs().Name)
	}
	return result, nil
}

func GetInterfacesWithMacs() (map[string]string, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get ip links")
	}
	result := make(map[string]string)
	for _, link := range links {
		name := link.Attrs().Name
		if name != "lo" {
			mac := link.Attrs().HardwareAddr.String()
			result[name] = mac
		}
	}
	return result, nil
}

func InterfaceUp(ifname string) error {
	link, err := netlink.LinkByName(ifname)
	if err != nil {
		return errors.Wrapf(err, "cannot find ip link %v", ifname)
	}
	err = netlink.LinkSetUp(link)
	if err != nil {
		return errors.Wrapf(err, "cannot set ip link up %v", ifname)
	}
	return nil
}

func InterfaceDhcp(ifname string) error {
	cmd := exec.Command("/sbin/udhcpc", "-i", ifname, "-b", "-p", fmt.Sprintf("/run/udhcpc.%s.pid", ifname))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		klog.V(5).Error(err, "cannot start dhcp client: %v", out.String())
	}
	return nil
}

func SetupLoopback() error {
	addr, err := netlink.ParseAddr("127.0.0.1/8")
	if err != nil {
		return errors.Wrapf(err, "cannot parse loopback address")
	}
	link, err := netlink.LinkByName("lo")
	if err != nil {
		return errors.Wrapf(err, "cannot find ip link lo")
	}
	err = netlink.AddrAdd(link, addr)
	if err != nil && err.Error() != "file exists" {
		return errors.Wrapf(err, "cannot set loopback ip address")
	}
	return nil
}

func StartNetworking() error {
	ifnames, err := GetInterfaces()
	if err != nil {
		return err
	}
	for _, ifname := range ifnames {
		err = InterfaceUp(ifname)
		if err != nil {
			return err
		}
		if ifname == "lo" {
			err = SetupLoopback()
			if err != nil {
				return errors.Wrapf(err, "cannot set loopback")
			}
		} else {
			InterfaceDhcp(ifname)
		}
	}
	return nil
}

func ListIpAddresses() ([]string, error) {
	addrs, err := netlink.AddrList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get ip addresses")
	}
	var result []string
	for _, addr := range addrs {
		result = append(result, addr.String())
	}
	return result, nil
}
