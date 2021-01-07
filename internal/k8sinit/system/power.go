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
	klog "k8s.io/klog/v2"
	"os"
)

func Poweroff() {
	klog.V(0).Infof("System will be powered off")
	unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
	os.Exit(0)
}

func Reboot() {
	klog.V(0).Infof("System will be rebooted")
	unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
	os.Exit(0)
}
