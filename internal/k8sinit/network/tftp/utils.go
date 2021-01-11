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
	"github.com/kazimsarikaya/k8sinit/internal/k8sinit"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
)

func DownloadIpxeUndi(tftproot string) (string, error) {
	resp, err := http.Get(k8sinit.UndiUrl)
	if err != nil {
		return "", errors.Wrapf(err, "cannot get undi pxe")
	}
	defer resp.Body.Close()

	filepath := tftproot + "/" + k8sinit.UndiFilename
	out, err := os.Create(filepath)
	if err != nil {
		return "", errors.Wrapf(err, "cannot create local undi pxe")
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "cannot download undi pxe")
	}
	return filepath, nil
}
