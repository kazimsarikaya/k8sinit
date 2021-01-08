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
	zfs "github.com/mistifyio/go-zfs"
	"github.com/pkg/errors"
	klog "k8s.io/klog/v2"
	"os/exec"
	"regexp"
	"strings"
)

type BlockDevice struct {
	Name          string
	Path          string
	PartitionType string `json:"pttype"`
	Size          uint64
}

func LoadZpools() error {
	cmd := exec.Command("zpool", "import")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "cannot list avaliable pools for import")
	}
	re := regexp.MustCompile(`pool: (.*)\b`)
	matches := re.FindAllStringSubmatch(out.String(), -1)
	if matches != nil {
		for _, match := range matches {
			if len(match) != 2 {
				return fmt.Errorf("pool name error :%v", match)
			}
			zpn := strings.TrimSpace(match[1])

			cmd = exec.Command("zpool", "import", zpn)
			if err := cmd.Run(); err != nil {
				return errors.Wrapf(err, "cannot import zpool %v", zpn)
			}
		}
	}
	return nil
}

func ListDisks() ([]*BlockDevice, error) {
	cmd := exec.Command("/bin/lsblk", "-o", "NAME,PATH,SIZE,PTTYPE", "-d", "-J", "-b")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrapf(err, "cannot list block devices")
	}

	dec := json.NewDecoder(&out)

	var data map[string][]*BlockDevice

	if err := dec.Decode(&data); err != nil {
		return nil, errors.Wrapf(err, "cannot decode lsblk output")
	}

	if bds, ok := data["blockdevices"]; ok {
		return bds, nil
	}

	return nil, fmt.Errorf("cannot find block devices")
}

func ListZpools() ([]*zfs.Zpool, error) {
	klog.V(5).Infof("list zpools called")
	zps, err := zfs.ListZpools()
	klog.V(5).Infof("zpool %v, err: %v", zps, err)
	return zps, err
}

func CloseZpools() {
	cmd := exec.Command("zpool", "export", "-a")
	if err := cmd.Run(); err != nil {
		klog.V(0).Error(err, "cannot export zpools,trying with force")
		cmd := exec.Command("zpool", "export", "-af")
		if err := cmd.Run(); err != nil {
			klog.V(0).Error(err, "cannot export zpools with force")
		}
	}
}
