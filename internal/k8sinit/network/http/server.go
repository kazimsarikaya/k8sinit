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

package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	klog "k8s.io/klog/v2"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type NonBlockingHttpServer struct {
	htdocs  string
	wg      *sync.WaitGroup
	server  *http.Server
	started bool
}

func NewNonBlockingHttpSever(htdocs string) (*NonBlockingHttpServer, error) {
	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	var wg sync.WaitGroup

	srv := &NonBlockingHttpServer{
		htdocs:  htdocs,
		wg:      &wg,
		started: false,
	}

	router.PathPrefix("/").HandlerFunc(srv.defaultHandler)

	server := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.server = server

	return srv, nil
}

func (s *NonBlockingHttpServer) defaultHandler(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(s.htdocs, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(s.htdocs, "index.html"))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(s.htdocs)).ServeHTTP(w, r)
}

func (s *NonBlockingHttpServer) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			klog.V(0).Error(err, "cannot start http server")
		}
	}()
	s.started = true
}

func (s *NonBlockingHttpServer) Stop() {
	if s.started {
		s.server.Shutdown(context.Background())
	}
	s.Wait()
	s.started = false
}

func (s *NonBlockingHttpServer) Wait() {
	s.wg.Wait()
}
