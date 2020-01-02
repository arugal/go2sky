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

import "github.com/pkg/errors"

type validateData struct {
	RegistryItem registryItem `json:"registryItem"`
	SegmentItems segmentItems `json:"segmentItems"`
}

type registryItem struct {
	Applications    map[string]int32   `json:"applications"`
	InstanceMapping map[string][]int32 `json:"instances"`
}

func (r *registryItem) RegistryApplication(applicationCode string, applicationId int32) {
	if _, ok := r.Applications[applicationCode]; !ok {
		r.Applications[applicationCode] = applicationId
	}
}

func (r *registryItem) RegistryInstance(applicationId int32, instanceId int32) error {
	applicationCode, err := r.findApplicationCode(applicationId)
	if err != nil {
		return err
	}

	instances, _ := r.InstanceMapping[applicationCode]
	r.InstanceMapping[applicationCode] = append(instances, instanceId)
	return nil
}

func (r *registryItem) findApplicationCode(applicationId int32) (string, error) {
	for k, v := range r.Applications {
		if v == applicationId {
			return k, nil
		}
	}
	return "", errors.Errorf("Cannot found the code of applicationID[ %d ].", applicationId)
}

type segmentItems struct {
}

type segmentItem struct {
	applicationCode string    `json:"applicationCode"`
	segments        []segment `json:"segments"`
}

type segment struct {
	segmentId string `json:"segmentId"`
	spans     []span `json:"spans"`
}

type span struct {
	OperationName string         `json:"operationName"`
	OperationId   int            `json:"operationId"`
	ParentSpanId  int            `json:"parentSpanId"`
	SpanId        int            `json:"spanId"`
	SpanLayer     int            `json:"spanLayer"`
	StartTime     int            `json:"startTime"`
	EndTime       int            `json:"endTime"`
	ComponentId   int            `json:"componentId"`
	ComponentName string         `json:"componentName"`
	IsError       bool           `json:"isError"`
	SpanType      string         `json:"spanType"`
	Peer          string         `json:"Peer"`
	PeerId        int            `json:"peerId"`
	Tags          []keyValuePair `json:"tags"`
	Logs          []logEvent     `json:"logs"`
	Refs          []segmentRef   `json:"refs"`
}

type segmentRef struct {
	ParentEndpointId        int    `json:"parentEndpointId"`
	ParentEndpoint          string `json:"parentEndpoint"`
	NetworkAddressId        int    `json:"networkAddressId"`
	EntryEndpointId         int    `json:"entryEndpointId"`
	RefType                 string `json:"refType"`
	ParentSpanId            int    `json:"parentSpanId"`
	ParentTraceSegmentId    string `json:"parentTraceSegmentId"`
	ParentServiceInstanceId int    `json:"parentServiceInstanceId"`
	NetworkAddress          string `json:"networkAddress"`
	EntryEndpoint           string `json:"entryEndpoint"`
	EntryServiceInstanceId  int    `json:"entryServiceInstanceId"`
}

type keyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type logEvent struct {
	LogEvent []keyValuePair `json:"logEvent"`
}
