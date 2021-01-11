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

package main

import (
	"flag"
	"fmt"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/management"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/mount"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/network"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/system"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/term"
	"github.com/pkg/errors"
	"io/ioutil"
	klog "k8s.io/klog/v2"
	"strings"
	"time"
)

var (
	version   = ""
	buildTime = ""
	goVersion = ""
	htdocsDir = ""
)

func init() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
}

func loader() error {
	err := system.FirstStep()
	if err != nil {
		return errors.Wrapf(err, "cannot execute first step")
	}
	err = mount.MountSysVFS()
	if err != nil {
		return errors.Wrapf(err, "error at mounting sys vfses")
	}
	err = system.LoadBaseModules()
	if err != nil {
		return errors.Wrapf(err, "error at mounting sys vfses")
	}

	role := system.GetRole()

	ic, err := system.ReadConfig()
	if err != nil {
		return errors.Wrapf(err, "cannot read config")
	}

	err = network.StartNetworking(ic)
	if err != nil {
		return errors.Wrapf(err, "cannot start networking")
	}
	klog.V(0).Infof("feeding random")
	rndfile := ""
	if ic != nil {
		rndfile = fmt.Sprintf("/%v/config/rndfile", ic.PoolName)
	}
	system.SeedRandom(rndfile)
	err = system.SetupDefaultApkRepos()
	if err != nil {
		return errors.Wrapf(err, "cannot setup apk")
	}

	tftproot, err := ioutil.TempDir("/tmp", "tftp")
	if err != nil {
		return errors.Wrapf(err, "cannot create tftp root dir")
	}

	var poolName, ifname string = "", ""
	if ic != nil {
		poolName = ic.PoolName
		ifname = ic.InternalNetwork
	}
	klog.V(0).Infof("setup management services")
	managementServices, err := management.NewOrGetManagementServices(role, poolName, ifname, tftproot, htdocsDir)
	if err != nil {
		return errors.Wrapf(err, "cannot setup management services")
	}
	system.SetManagementServicesStopper(managementServices)
	managementServices.StartHttp()
	if role == k8sinit.RoleManager {
		managementServices.StartTftp(strings.Split(ic.InternalNetworkIPAndPrefix, "/")[0])
		managementServices.StartDhcp()
	}
	return nil
}

func showUI() error {
	klog.V(0).Infof("entering ui")
	for {
		term.ClearScreen()
		cmd, err := term.ReadKeyPress()
		if err != nil {
			klog.V(0).Error(err, "cannot get command")
		}
		if cmd == 'C' {
			err = term.CreateTerminal()
			if err != nil {
				klog.V(0).Error(err, "error occured")
			}
		} else if cmd == 'P' {
			system.Poweroff()
		} else if cmd == 'R' {
			system.Reboot()
		} else {
			klog.V(0).Infof("Unknown command...")
			time.Sleep(time.Second * 5)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	klog.V(0).Infof("hello from k3s init")
	err := loader()
	if err != nil {
		klog.V(0).Error(err, "cannot load system")
	} else {
		err = showUI()
		if err != nil {
			klog.V(0).Error(err, "error at ui")
		}
	}

	for {
		klog.V(0).Infof("Sleeping...")
		time.Sleep(time.Hour)
	}
}
