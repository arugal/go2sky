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
	RegistryItem registryItem `json:"registry_item"`
	SegmentItems segmentItems `json:"segment_items"`
}

type registryItem struct {
	Applications    map[string]int32   `json:"applications"`
	InstanceMapping map[string][]int32 `json:"instance_mapping"`
}

func (r *registryItem) RegistryApplication(applicationCode string, applicationId int32) {
	if _, ok := r.Applications[applicationCode]; !ok {
		r.Applications[applicationCode] = applicationId
	}
}

func (r *registryItem) RegistryInstance(applicationId int32, instanceId int32) error {
	_, err := r.findApplicationCode(applicationId)
	if err != nil {
		return err
	}

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

type Span struct {
	OperationName string `json:"operation_name"`
	OperationId   int    `json:"operation_id"`
	ParentSpanId  int    `json:"parent_span_id"`
	SpanId        int    `json:"span_id"`
	SpanLayer     int    `json:"span_layer"`
	StartTime     int    `json:"start_time"`
	EndTime       int    `json:"end_time"`
}
