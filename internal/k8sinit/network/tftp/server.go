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

package tftp

import (
	"fmt"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit"
	"github.com/pin/tftp"
	"io"
	klog "k8s.io/klog/v2"
	"os"
	"sync"
	"time"
)

type NonBlockingTftpSever struct {
	tftproot string
	undi     string
	server   *tftp.Server
	wg       *sync.WaitGroup
	started  bool
}

func NewNonBlockingTftpSever(tftproot string) (*NonBlockingTftpSever, error) {
	var wg sync.WaitGroup

	s := &NonBlockingTftpSever{
		tftproot: tftproot,
		wg:       &wg,
		started:  false,
	}

	server := tftp.NewServer(s.readHandler, nil)
	server.SetTimeout(5 * time.Second)

	s.server = server

	return s, nil
}

func (s *NonBlockingTftpSever) Start(ipaddr string) {
	undi, err := DownloadIpxeUndi(s.tftproot)
	if err != nil {
		klog.V(0).Error(err, "cannot download "+k8sinit.UndiFilename)
	}
	s.undi = undi
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := s.server.ListenAndServe(ipaddr + ":69")
		if err != nil {
			klog.V(0).Error(err, "cannot start tftp server")
		}
	}()
	s.started = true
}

func (s *NonBlockingTftpSever) Stop() {
	if s.started {
		s.server.Shutdown()
	}
	s.Wait()
	s.started = false
}

func (s *NonBlockingTftpSever) Wait() {
	s.wg.Wait()
}

func (s *NonBlockingTftpSever) readHandler(filename string, rf io.ReaderFrom) error {
	if filename != k8sinit.UndiFilename {
		return fmt.Errorf("only %s supported", k8sinit.UndiFilename)
	}
	file, err := os.Open(s.undi)
	if err != nil {
		klog.V(5).Error(err, "cannot open undi pxe file")
		return err
	}
	_, err = rf.ReadFrom(file)
	if err != nil {
		klog.V(5).Error(err, "cannot send undi pxe file")
		return err
	}
	return nil
}
