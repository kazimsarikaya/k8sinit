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
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/system"
	"io"
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
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	messageType, sr, err := conn.NextReader()
	if err != nil {
		return
	}
	var ic system.InstallConfig
	err = json.NewDecoder(sr).Decode(&ic)
	if err != nil {
		conn.WriteMessage(messageType, []byte(fmt.Sprintf("error: cannot decode json data err: %v mt: %v", err, messageType)))
		conn.Close()
		return
	}
	pr, pw := io.Pipe()
	go func() {
		scanner := bufio.NewScanner(pr)
		for {
			for scanner.Scan() {
				conn.WriteMessage(messageType, scanner.Bytes())
			}
			if err := scanner.Err(); err != nil {
				conn.WriteMessage(messageType, []byte("error occured "+err.Error()))
				break
			}
		}
		conn.Close()
	}()
	err = system.InstallSystem(ic, pw)
	if err != nil {
		conn.WriteMessage(messageType, []byte("error: installaction error"))
		conn.Close()
	}
}
