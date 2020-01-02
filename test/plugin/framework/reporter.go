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
	"encoding/json"
	"github.com/SkyAPM/go2sky"
	"log"
	"os"
	"sync/atomic"
)

var (
	instanceSequence int32 = 0
	endpointSequence int32 = 1
	serviceMapping   map[string]int32
)

func instanceId() int32 {
	return atomic.AddInt32(&instanceSequence, 1)
}

func serviceId(service string) (serviceId int32) {
	if serviceId, ok := serviceMapping[service]; !ok {
		serviceId = atomic.AddInt32(&endpointSequence, 1)
		serviceMapping[service] = serviceId
	}
	return
}

func NewMockReporter() mockReporter {
	return mockReporter{
		logger: log.New(os.Stderr, "go2sky-test", log.LstdFlags),
		validateData: validateData{
			RegistryItem: registryItem{
				Applications:    make(map[string]int32),
				InstanceMapping: make(map[string][]int32),
			},
		},
	}
}

type mockReporter struct {
	logger       *log.Logger
	validateData validateData
}

func (mr *mockReporter) Register(applicationCode string, instance string) (int32, int32, error) {
	mr.logger.Println("Register mock reporter")
	return serviceId(applicationCode), instanceId(), nil
}

func (mr *mockReporter) Send(spans []go2sky.ReportedSpan) {
	if spans == nil {
		return
	}
	b, err := json.Marshal(spans)
	if err != nil {
		mr.logger.Printf("Error: %s", err)
		return
	}
	root := spans[len(spans)-1]
	mr.logger.Printf("Segment-%v: %s \n", root.Context().SegmentID, b)
}

func (mr *mockReporter) Close() {
	mr.logger.Println("Close mock reporter")
}

func (mr *mockReporter) ReceiveData() {

}
