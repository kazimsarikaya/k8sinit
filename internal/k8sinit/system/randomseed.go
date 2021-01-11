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
	"crypto/rand"
	"fmt"
	"io/ioutil"
	klog "k8s.io/klog/v2"
	"os"
	"strconv"
	"strings"
)

func getEntropyCount() int64 {
	data, _ := ioutil.ReadFile("/proc/sys/kernel/random/entropy_avail")
	ecnt, _ := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 32)
	klog.V(0).Infof("entropy count %v", ecnt)
	return ecnt
}

func getRndWakeupThreshold() int64 {
	data, _ := ioutil.ReadFile("/proc/sys/kernel/random/read_wakeup_threshold")
	ecnt, _ := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 32)
	klog.V(0).Infof("threashold count %v", ecnt)
	return ecnt
}

func writeRandomSeed() {
	ic, _ := ReadConfig()
	if ic != nil {
		rndfile := fmt.Sprintf("/%v/config/rndfile", ic.PoolName)
		data := make([]byte, 1024*1024*16)
		in, _ := os.Open("/dev/urandom")
		in.Read(data)
		ioutil.WriteFile(rndfile, data, 0600)
	}
}

func SeedRandom(rndfile string) {
	thres := getRndWakeupThreshold()
	max_ecnt := getEntropyCount()
	if rndfile != "" {
		data, err := ioutil.ReadFile(rndfile)
		if err == nil {
			ioutil.WriteFile("/dev/random", data, 0666)
			ioutil.WriteFile("/dev/urandom", data, 0666)
			ecnt := getEntropyCount()
			if ecnt > max_ecnt {
				max_ecnt = ecnt
			}
			if ecnt >= thres || ecnt < max_ecnt {
				klog.V(0).Infof("seed random done")
				return
			}
		}
	}
	max_ecnt = getEntropyCount()
	for {
		klog.V(0).Infof("try to seed random, please produce entropy")
		data := make([]byte, 1024*1024*16)
		r, err := rand.Read(data)
		klog.V(0).Infof("read %v random data", r)
		if err == nil {
			ioutil.WriteFile("/dev/random", data, 0666)
			ioutil.WriteFile("/dev/urandom", data, 0666)
			ecnt := getEntropyCount()
			if ecnt > max_ecnt {
				max_ecnt = ecnt
			}
			if ecnt >= thres || ecnt < max_ecnt {
				klog.V(0).Infof("seed random done")
				return
			}
		} else {
			klog.V(0).Error(err, "cannot read random data")
		}
	}
}
