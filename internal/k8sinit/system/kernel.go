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
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

func GetKernelParameterValue(paramKey string) (bool, interface{}, error) {
	data, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		return false, nil, errors.Wrapf(err, "cannot read kernel parameters")
	}
	params := strings.Split(string(data), " ")
	var result []interface{}

	found := false

	for _, param := range params {
		keyval := strings.SplitN(param, "=", 2)
		if len(keyval) == 2 {
			if keyval[0] == paramKey {
				found = true
				result = append(result, strings.TrimSpace(keyval[1]))
			}
		} else if len(keyval) == 1 {
			if keyval[0] == paramKey {
				found = true
			}
		}
	}

	if len(result) == 0 {
		return found, nil, nil
	} else if len(result) == 1 {
		return found, result[0], nil
	}
	return found, result, nil
}
