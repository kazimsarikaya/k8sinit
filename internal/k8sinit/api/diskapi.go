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
	"github.com/gorilla/mux"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/system"
	"net/http"
)

func DiskApiListBlockDevices(w http.ResponseWriter, r *http.Request) {
	bds, err := system.ListDisks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": bds})
}

func DiskApiListZpools(w http.ResponseWriter, r *http.Request) {
	zps, err := system.ListZpools()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": zps})
}

func DiskApiGetZpool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pool, ok := vars["pool"]
	if !ok || pool == "" {
		http.Error(w, "no pool param", http.StatusBadRequest)
		return
	}
	zps, err := system.ListZpools()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, zp := range zps {
		if zp.Name == pool {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": zp})
			return
		}
	}
	http.Error(w, "pool not found", http.StatusNotFound)
}

func DiskApiListDatasets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pool, ok := vars["pool"]
	if !ok || pool == "" {
		http.Error(w, "no pool param", http.StatusBadRequest)
		return
	}
	zps, err := system.ListZpools()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, zp := range zps {
		if zp.Name == pool {
			dses, err := zp.Datasets()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": dses})
			return
		}
	}
	http.Error(w, "pool not found", http.StatusNotFound)
}

func DiskApiGetDataset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pool, ok := vars["pool"]
	if !ok || pool == "" {
		http.Error(w, "no pool param", http.StatusBadRequest)
		return
	}
	dataset, ok := vars["dataset"]
	if !ok || dataset == "" {
		http.Error(w, "no dataset param", http.StatusBadRequest)
		return
	}
	zps, err := system.ListZpools()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, zp := range zps {
		if zp.Name == pool {
			dses, err := zp.Datasets()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			for _, ds := range dses {
				if ds.Name == dataset {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": ds})
					return
				}
			}
			http.Error(w, "dataset not found", http.StatusNotFound)
			return
		}
	}
	http.Error(w, "pool not found", http.StatusNotFound)
}
