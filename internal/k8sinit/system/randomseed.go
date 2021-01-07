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
	"io/ioutil"
	"math/rand"
	"time"
)

func SeedRandom() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var data [512]byte
	rng.Read(data[:])
	ioutil.WriteFile("/dev/random", data[:], 0666)
	rng.Read(data[:])
	ioutil.WriteFile("/dev/urandom", data[:], 0666)
}
