// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package framework

import (
	"fmt"
	"testing"
)

var (
	instanceMap = make(map[string][]int32)
)

func TestRegis(t *testing.T) {
	instances, _ := instanceMap["a"]
	instances = append(instances, 1)
	instanceMap["a"] = append(instances, 1)

	fmt.Printf("instance: %v\n", instanceMap["a"])
}
