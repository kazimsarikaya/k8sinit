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
	"encoding/json"
	"fmt"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"io"
	klog "k8s.io/klog/v2"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

var (
	singletonIC *k8sinit.InstallConfig = nil
)

func FirstStep() error {
	unix.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	os.Symlink("/init", "/sbin/reboot")
	os.Symlink("/init", "/sbin/poweroff")

	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGUSR1, unix.SIGUSR2)
	go func() {
		sig := <-c
		if sig == unix.SIGUSR1 {
			Reboot()
		}
		if sig == unix.SIGUSR2 {
			Poweroff()
		}
	}()

	cmd := exec.Command("/bin/busybox", "--install", "-s")
	var out bytes.Buffer
	cmd.Stdin = strings.NewReader("")
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "cannot setup busybox: %v", out.String())
	}
	return nil
}

func InstallSystem(config k8sinit.InstallConfig, output io.WriteCloser) error {
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
	if err = copyOsFilesToDisk(config.PoolName, output); err != nil {
		klog.V(0).Error(err, "copying os files failed")
		return err
	}
	if err = grubInstall(config.Disk, config.PoolName, output); err != nil {
		klog.V(0).Error(err, "cannot install grub")
		return err
	}
	if err = WriteConfig(config); err != nil {
		klog.V(0).Error(err, "config write failed")
		return errors.Wrapf(err, "config write failed")
	}
	klog.V(0).Infof("installtion ended")
	output.Write([]byte("installation ended\neject cdrom and reboot\n"))
	return nil
}

func WriteConfig(config k8sinit.InstallConfig) error {
	singletonIC = &config
	writeRandomSeed()
	out, err := os.OpenFile("/"+config.PoolName+"/config/config.json", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrapf(err, "cannot create config file")
	}
	defer out.Close()
	enc := json.NewEncoder(out)
	if err := enc.Encode(config); err != nil {
		return errors.Wrapf(err, "cannot write config")
	}
	return nil
}

func ReadConfig() (*k8sinit.InstallConfig, error) {
	if singletonIC != nil {
		return singletonIC, nil
	}
	found, poolName, err := GetKernelParameterValue("k8sinit.pool")
	if !found {
		poolName = "zp_k8s"
	}
	if err != nil {
		return nil, err
	}
	in, err := os.Open("/" + poolName.(string) + "/config/config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "cannot open config")
	}
	defer in.Close()
	var config k8sinit.InstallConfig
	if err := json.NewDecoder(in).Decode(&config); err != nil {
		return nil, errors.Wrapf(err, "cannot decode config")
	}
	singletonIC = &config
	return singletonIC, nil
}

func GetRole() string {
	found, role, _ := GetKernelParameterValue("k8sinit.role")
	if !found {
		role = "node"
	}
	return role.(string)
}
