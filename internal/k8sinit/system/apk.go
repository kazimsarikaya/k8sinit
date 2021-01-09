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
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
)

func SetupDefaultApkRepos() error {
	cmd := exec.Command("apk", "add", "--initdb")
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "cannot init apk db")
	}

	repos := `http://dl-cdn.alpinelinux.org/alpine/v3.12/main
http://dl-cdn.alpinelinux.org/alpine/v3.12/community
`
	err := ioutil.WriteFile("/etc/apk/repositories", []byte(repos), 0644)
	if err != nil {
		return errors.Wrapf(err, "cannot setup apk repositories")
	}
	return nil
}

func apkInstallPacket(packetname string) error {
	cmd := exec.Command("apk", "add", packetname)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "cannot add packet: %v", packetname)
	}
	return nil
}
