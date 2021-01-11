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

package system

import (
	"golang.org/x/sys/unix"
	"io/ioutil"
	klog "k8s.io/klog/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type StopAller interface {
	StopAll()
}

var managementServicesStopper StopAller

func SetManagementServicesStopper(ms StopAller) {
	managementServicesStopper = ms
}

func reapProcs() {
	cmdlines, err := filepath.Glob("/proc/*/cmdline")
	if err != nil {
		klog.V(0).Error(err, "error during fetching cmdlines")
	}
	for _, cmdline := range cmdlines {
		if cmdline == "/proc/1/cmdline" {
			continue
		}
		data, err := ioutil.ReadFile(cmdline)
		if err != nil || len(data) == 0 {
			continue
		}
		pidstr := strings.TrimSuffix(cmdline[6:], "/cmdline")
		pid, err := strconv.ParseInt(pidstr, 10, 32)
		if err != nil {
			continue
		}
		p, err := os.FindProcess(int(pid))
		p.Kill()
	}
}

func stopSystem() {
	managementServicesStopper.StopAll()
	writeRandomSeed()
	CloseZpools()
	reapProcs()
	klog.V(0).Infof("System stopped")
	klog.Flush()
}

func Poweroff() {
	klog.V(0).Infof("System will be powered off")
	stopSystem()
	unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
	os.Exit(0)
}

func Reboot() {
	klog.V(0).Infof("System will be rebooted")
	stopSystem()
	unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
	os.Exit(0)
}

func SendPoweroff() {
	p, err := os.FindProcess(1)
	if err != nil {
		return
	}
	p.Signal(unix.SIGUSR2)
}
func SendReboot() {
	p, err := os.FindProcess(1)
	if err != nil {
		return
	}
	p.Signal(unix.SIGUSR1)
}
