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
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
)

type MountData struct {
	Fstype  string
	Source  string
	Target  string
	Options []string
}

const (
	sysvfses string = `[{
      "Fstype": "proc",
      "Source": "proc",
      "Target": "/proc",
      "Options": ["rw", "nosuid", "nodev", "noexec", "relatime"]
    },
    {
      "Fstype": "sysfs",
      "Source": "sysfs",
      "Target": "/sys",
      "Options": ["rw", "nosuid", "nodev", "noexec", "relatime"]
    },
    {
      "Fstype": "devtmpfs",
      "Source": "devtmpfs",
      "Target": "/dev",
      "Options": ["rw", "nosuid", "relatime", "size=10m", "nr_inodes=500444", "mode=755"]
    },
    {
      "Fstype": "devpts",
      "Source": "devpts",
      "Target": "/dev/pts",
      "Options": ["rw", "nosuid", "noexec", "relatime", "gid=5", "mode=620", "ptmxmode=000"]
    },
    {
      "Fstype": "tmpfs",
      "Source": "shm",
      "Target": "/dev/shm",
      "Options": ["rw", "nosuid", "nodev", "noexec", "relatime"]
    }
  ]`
)

func mount(fstype, source, target string, options []string) error {
	mo := strings.Join(options, ",")
	err := os.MkdirAll(target, 0755)
	if err != nil {
		return errors.Wrapf(err, "cannot create target dir %v", target)
	}
	cmd := exec.Command("/bin/mount", "-t", fstype, "-o", mo, source, target)
	var out bytes.Buffer
	cmd.Stdin = strings.NewReader("")
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func MountSysVFS() error {
	var mds []MountData
	err := json.Unmarshal([]byte(sysvfses), &mds)
	if err != nil {
		return errors.Wrapf(err, "cannot unmarshall sys vfs data")
	}
	for _, md := range mds {
		if err = mount(md.Fstype, md.Source, md.Target, md.Options); err != nil {
			return errors.Wrapf(err, "cannot mount %v", md)
		}
	}
	err = os.RemoveAll("/etc/mtab")
	if err != nil {
		return errors.Wrapf(err, "cannot remove /etc/mtab")
	}
	err = os.Symlink("/proc/mounts", "/etc/mtab")
	if err != nil {
		return errors.Wrapf(err, "cannot symlink /etc/mtab to /proc/mounts")
	}
	return mountCgroups()
}
