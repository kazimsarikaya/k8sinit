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

func main() {
	flag.Parse()
	klog.V(0).Infof("hello from k3s init")
	err := system.FirstStep()
	if err != nil {
		klog.V(0).Error(err, "cannot execute first step")
	} else {
		err = mount.MountSysVFS()
		if err != nil {
			klog.V(0).Error(err, "error at mounting sys vfses")
		} else {
			err = modules.LoadBaseModules()
			if err != nil {
				klog.V(0).Error(err, "error at mounting sys vfses")
			} else {
				err = network.StartNetworking()
				if err != nil {
					klog.V(0).Error(err, "cannot start networking")
				} else {
					klog.V(0).Infof("feeding random")
					system.SeedRandom()
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
				}
			}
		}
	}

	for {
		klog.V(0).Infof("Sleeping...")
		time.Sleep(time.Hour)
	}
}
