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
	"github.com/pkg/errors"
	"io/ioutil"
	klog "k8s.io/klog/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func findModaliases() ([]string, error) {
	var result []string
	err := filepath.Walk("/sys", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if filepath.Base(path) == "modalias" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			klog.V(5).Infof("modalias path: %v data: %v", path, data)
			result = append(result, strings.TrimSpace(string(data)))
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find modalias files")
	}
	return result, nil
}

func listModules() ([]string, error) {
	cmd := exec.Command("/sbin/lsmod")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrapf(err, "cannot list modules")
	}
	mods := strings.Split(out.String(), "\n")
	return mods[1:], nil
}

func Modprobe(moddata string) error {
	cmd := exec.Command("/sbin/modprobe", "-a", "-b", moddata)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "cannot load module %v", moddata)
	}
	return nil
}

func ModprobeNoerr(moddata string) {
	err := Modprobe(moddata)
	if err != nil {
		klog.V(5).Error(err, "soft error occured")
	}
}

func LoadBaseModules() error {
	for {
		mods, err := listModules()
		if err != nil {
			return err
		}
		old_mod_cnt := len(mods)
		mas, err := findModaliases()
		if err != nil {
			return err
		}
		for _, ma := range mas {
			ModprobeNoerr(ma)
		}
		mods, err = listModules()
		if err != nil {
			return err
		}
		new_mod_cnt := len(mods)
		klog.V(5).Infof("old mod count: %d new mod count: %d", old_mod_cnt, new_mod_cnt)
		if old_mod_cnt == new_mod_cnt {
			break
		}
	}
	if err := Modprobe("zfs"); err != nil {
		return errors.Wrapf(err, "cannot load zfs module")
	}
	if err := LoadZpools(); err != nil {
		return err
	}
	return nil
}
