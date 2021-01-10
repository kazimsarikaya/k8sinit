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
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"io"
	klog "k8s.io/klog/v2"
	"os/exec"
	"strings"
)

type InstallConfig struct {
	Disk                       string `json:"disk"`
	Force                      bool   `json:"force"`
	PoolName                   string `json:"poolname"`
	ExternalNetwork            string `json:"extnet"`
	IsExternalNetworkStatic    bool   `json:"extnettype"`
	ExternalNetworkIPAndPrefix string `json:"extnetip"`
	ExternalNetworkGateway     string `json:"extnetgw"`
	AdminNetwork               string `json:"adminnet"`
	IsAdminNetworkStatic       bool   `json:"adminnettype"`
	AdminNetworkIPAndPrefix    string `json:"adminnetip"`
	InternalNetwork            string `json:"internalnet"`
	IsInternalNetworkStatic    bool   `json:"internalnettype"`
	InternalNetworkIPAndPrefix string `json:"internalnetip"`
}

func FirstStep() error {
	unix.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	cmd := exec.Command("/bin/busybox", "--install", "-s")
	var out bytes.Buffer
	cmd.Stdin = strings.NewReader("")
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "cannot setup busybox")
	}
	return nil
}

func InstallSystem(config InstallConfig, output io.WriteCloser) error {
	defer output.Close()
	klog.V(0).Infof("starting install")
	output.Write([]byte("starting install\n"))
	err := apkInstallPacketWithOutput("grub-bios", output)
	if err != nil {
		klog.V(0).Error(err, "cannot install apk deps")
		return errors.Wrapf(err, "cannot install grub-bios")
	}
	klog.V(0).Infof("apk deps installed")
	output.Write([]byte("apk deps installed\n"))
	zps, err := ListZpools()
	if err != nil {
		klog.V(0).Error(err, "cannot list zpools")
		output.Write([]byte(err.Error()))
		return errors.Wrapf(err, "cannot get zpool info")
	}
	for _, zp := range zps {
		if zp.Name == config.PoolName {
			klog.V(0).Infof("same zpool found")
			output.Write([]byte("same zpool found\n"))
			if !config.Force {
				errstr := "pool exists with same name and force parameter not given"
				output.Write([]byte(errstr))
				return fmt.Errorf(errstr)
			} else {
				err = zp.Destroy()
				if err != nil {
					output.Write([]byte(err.Error()))
					return errors.Wrapf(err, "cannot destroy zpool")
				}
				output.Write([]byte("zpool destroyed\n"))
				klog.V(0).Infof("zpool destroyed")
			}
			break
		}
	}
	if err = partDisk(config.Disk, output); err != nil {
		klog.V(0).Error(err, "partitioning failed")
		return err
	}
	if err = createZfs(config.Disk+"2", config.PoolName, output); err != nil {
		klog.V(0).Error(err, "create zfs failed")
		return err
	}
	klog.V(0).Infof("installtion ended")
	output.Write([]byte("installation ended\n"))
	return nil
}
