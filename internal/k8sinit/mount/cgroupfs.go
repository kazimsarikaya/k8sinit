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

package mount

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	klog "k8s.io/klog/v2"
	"os"
	"strings"
)

func activateControllers(path string) error {
	// bdata, err := ioutil.ReadFile(path + "/cgroup.controllers")
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		return nil
	// 	}
	// 	return errors.Wrapf(err, "cannot read controllers of %v", path)
	// }
	// controllers := strings.Split(string(bdata), " ")
	// klog.V(0).Infof("controllers: %v", controllers)
	// for _, controller := range controllers {
	// 	err = ioutil.WriteFile(path+"/cgroup.subtree_control", []byte(fmt.Sprintf("+%s", controller)), 0644)
	// 	if err != nil {
	// 		return errors.Wrapf(err, "cannot activate controller %v of %v", controller, path)
	// 	}
	// }
	return nil
}

func getEnabledCgroups() ([]string, error) {
	var result []string
	bdata, err := ioutil.ReadFile("/proc/cgroups")
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read /proc/cgroups")
	}
	lines := strings.Split(string(bdata), "\n")
	for _, line := range lines {
		if len(line) <= 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		if len(parts) != 4 {
			return nil, fmt.Errorf("malformed /proc/mounts line: %v", line)
		}
		if parts[3] == "1" {
			result = append(result, parts[0])
		}
	}
	return result, nil
}

func mountCgroup(cg, cgtype, src string, opts []string) error {
	klog.V(5).Infof("mounting cgroup %v", cg)
	if err := os.MkdirAll("/sys/fs/cgroup/"+cg, 0755); err != nil {
		return errors.Wrapf(err, "cannot create cgroup %v dir", cg)
	}
	if err := mount(cgtype, src, "/sys/fs/cgroup/"+cg, opts); err != nil {
		return errors.Wrapf(err, "cannot mount %v as %v with %v", cg, cgtype, opts)
	} else if err := activateControllers("/sys/fs/cgroup/" + cg); err != nil {
		return err
	}
	return nil
}

func mountCgroups() error {
	defaults := []string{"rw", "nosuid", "nodev", "noexec", "relatime"}
	cgroup_root_opts := []string{"size=10m", "mode=755"}
	if err := mount("tmpfs", "cgroup_root", "/sys/fs/cgroup", append(defaults, cgroup_root_opts...)); err != nil {
		return errors.Wrapf(err, "cannot mount base dir of cgroups")
	}
	if err := mountCgroup("unified", "cgroup2", "none", append(defaults, "nsdelegate")); err != nil {
		return err
	}
	if cgs, err := getEnabledCgroups(); err == nil {
		for _, cg := range cgs {
			if err := mountCgroup(cg, "cgroup", cg, append(defaults, cg)); err != nil {
				return err
			}
		}
	} else {
		return err
	}
	return nil
}
