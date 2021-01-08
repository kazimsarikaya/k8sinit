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
)

type BlockDevice struct {
	Name          string
	Path          string
	PartitionType string `json:"pttype"`
	Size          uint64
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
	klog.V(0).Infof("zpool %v, err: %v", zps, err)
	return zps, err
}
