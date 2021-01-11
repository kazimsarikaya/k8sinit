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

package api

import (
	"encoding/json"
	"fmt"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/network"
	klog "k8s.io/klog/v2"
	"net/http"
)

func NetworkApiInterfaceList(w http.ResponseWriter, r *http.Request) {
	res, err := network.GetInterfacesWithMacs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": true, "data": res})
}

func NetworkApiTftp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `#!ipxe
echo loading kernel...
kernel http://%s/api/network/tftp/vmlinuz k8sinit.role=node k8sinit.pool=%s
echo loading initrd...
initrd http://%s/api/network/tftp/initrd
boot
`, r.Host, "zp_k8s", r.Host)
}

func NetworkApiTftpVmlinuz(w http.ResponseWriter, r *http.Request) {
	klog.V(0).Infof("start sending vmlinuz")
	http.ServeFile(w, r, "/zp_k8s/boot/vmlinuz") // TODO: get base path from config
	klog.V(0).Infof("sending vmlinuz ended")
}

func NetworkApiTftpInitrd(w http.ResponseWriter, r *http.Request) {
	klog.V(0).Infof("start sending initramfs")
	http.ServeFile(w, r, "/zp_k8s/boot/initramfs") // TODO: get base path from config
	klog.V(0).Infof("sending initramfs ended")
}
