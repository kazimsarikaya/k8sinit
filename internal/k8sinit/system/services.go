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
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/network/http"
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit/network/tftp"
	klog "k8s.io/klog/v2"
)

type ManagementServices struct {
	tftpServer *tftp.NonBlockingTftpSever
	httpServer *http.NonBlockingHttpServer
}

func NewManagementServices(tftpDir, htdocsDir string) (*ManagementServices, error) {
	ms := &ManagementServices{}
	var err error
	ms.tftpServer, err = tftp.NewNonBlockingTftpSever(tftpDir)
	if err != nil {
		return nil, err
	}
	ms.httpServer, err = http.NewNonBlockingHttpSever(htdocsDir)
	if err != nil {
		return nil, err
	}
	return ms, nil
}

func (ms *ManagementServices) StopAll() {
	if ms.tftpServer != nil {
		klog.Infof("stopping tftp server")
		klog.Flush()
		ms.tftpServer.Stop()
	}
	if ms.httpServer != nil {
		klog.Infof("stopping http server")
		klog.Flush()
		ms.httpServer.Stop()
	}
}

func (ms *ManagementServices) StartTftp(ipaddr string) {
	ms.tftpServer.Start(ipaddr)
}

func (ms *ManagementServices) StartHttp() {
	ms.httpServer.Start()
}
