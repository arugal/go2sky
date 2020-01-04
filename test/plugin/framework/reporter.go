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
	"sync/atomic"

	"github.com/SkyAPM/go2sky"
)

var (
	instanceSequence int32 = 0
	endpointSequence int32 = 1
	serviceMapping         = make(map[string]int32)
)

func NewTracer(applicationCode string) (tracer *go2sky.Tracer, err error) {
	return go2sky.NewTracer(applicationCode, go2sky.WithReporter(newMockReporter()))
}

func newMockReporter() go2sky.Reporter {
	return &mockReporter{}
}

type mockReporter struct {
	applicationID int32
	instanceID    int32
}

func (mr *mockReporter) Register(service string, instance string) (int32, int32, error) {
	mr.applicationID = serviceID(service)
	registryApplication(service, mr.applicationID)
	mr.instanceID = instanceID()
	if err := registryInstance(mr.applicationID, mr.instanceID); err != nil {
		return 0, 0, err
	}
	return mr.applicationID, mr.instanceID, nil
}

func (mr *mockReporter) Send(spans []go2sky.ReportedSpan) {
	if spans == nil {
		return
	}
	_ = addSpans(mr.applicationID, mr.instanceID, spans)
}

func (mr *mockReporter) Close() {
}

func instanceID() int32 {
	return atomic.AddInt32(&instanceSequence, 1)
}

func serviceID(service string) int32 {
	var serviceID int32
	var ok bool
	if serviceID, ok = serviceMapping[service]; !ok {
		serviceID = atomic.AddInt32(&endpointSequence, 1)
		serviceMapping[service] = serviceID
	}
	return serviceID
}
