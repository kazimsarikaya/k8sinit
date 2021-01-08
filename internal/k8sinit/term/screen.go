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

package term

import (
	"fmt"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/network"
	"github.com/pkg/errors"
	klog "k8s.io/klog/v2"
	"os"
	"strings"
)

func ClearScreen() {
	os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
	os.Stdout.WriteString(`
    ██╗  ██╗ █████╗ ███████╗
    ██║ ██╔╝██╔══██╗██╔════╝
    █████╔╝ ╚█████╔╝███████╗
    ██╔═██╗ ██╔══██╗╚════██║
    ██║  ██╗╚█████╔╝███████║
    ╚═╝  ╚═╝ ╚════╝ ╚══════╝
    ██╗███╗   ██╗██╗████████╗
    ██║████╗  ██║██║╚══██╔══╝
    ██║██╔██╗ ██║██║   ██║
    ██║██║╚██╗██║██║   ██║
    ██║██║ ╚████║██║   ██║
    ╚═╝╚═╝  ╚═══╝╚═╝   ╚═╝
    ███████╗██╗   ██╗███████╗████████╗███████╗███╗   ███╗
    ██╔════╝╚██╗ ██╔╝██╔════╝╚══██╔══╝██╔════╝████╗ ████║
    ███████╗ ╚████╔╝ ███████╗   ██║   █████╗  ██╔████╔██║
    ╚════██║  ╚██╔╝  ╚════██║   ██║   ██╔══╝  ██║╚██╔╝██║
    ███████║   ██║   ███████║   ██║   ███████╗██║ ╚═╝ ██║
    ╚══════╝   ╚═╝   ╚══════╝   ╚═╝   ╚══════╝╚═╝     ╚═╝

`)

	addrs, err := network.ListIpAddresses()
	if err != nil {
		klog.V(5).Error(err, "cannot get ip addresses")
	} else {
		list := strings.Join(addrs, ",")
		os.Stdout.WriteString(fmt.Sprintf("Ip Adresses: %s\n\n", list))
	}
	os.Stdout.WriteString(`For Console press   C
For Poweroff press  P
For Reboot press    R
`)
	os.Stdout.Sync()
}

func ReadKeyPress() (byte, error) {

	var char [10]byte
	_, err := os.Stdin.Read(char[:])
	if err != nil {
		return 0, errors.Wrapf(err, "cannot read command")
	}

	return char[0], nil
}
