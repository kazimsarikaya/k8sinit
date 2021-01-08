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
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/modules"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/mount"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/network"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/system"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/term"
	"github.com/pkg/errors"
	klog "k8s.io/klog/v2"
	"time"
)

var (
	version   = ""
	buildTime = ""
	goVersion = ""
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
	err = modules.LoadBaseModules()
	if err != nil {
		return errors.Wrapf(err, "error at mounting sys vfses")
	}
	err = network.StartNetworking()
	if err != nil {
		return errors.Wrapf(err, "cannot start networking")
	}
	klog.V(0).Infof("feeding random")
	system.SeedRandom()
	err = system.SetupDefaultApkRepos()
	if err != nil {
		return errors.Wrapf(err, "cannot setup apk")
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
