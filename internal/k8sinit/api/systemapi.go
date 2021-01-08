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
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/system"
	"net/http"
	"time"
)

func SystemApiReboot(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(time.Second * 15)
		system.Reboot()
	}()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": true, "data": "system will be rebooted in 15 seconds"})
}

func SystemApiPoweroff(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(time.Second * 15)
		system.Poweroff()
	}()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "data": "system will be poweroffed in 15 seconds"})
}

func SystemApiInstall(w http.ResponseWriter, r *http.Request) {

}
