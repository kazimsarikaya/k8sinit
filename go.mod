// Copyright 2020 Kazım SARIKAYA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

module github.com/kazimsarikaya/k8sinit

go 1.15

require (
	github.com/creack/pty v1.1.11
	github.com/google/uuid v1.1.4 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible
	github.com/pin/tftp v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/vishvananda/netlink v1.1.0
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf
	k8s.io/klog/v2 v2.4.0
)
