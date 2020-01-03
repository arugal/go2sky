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
	"github.com/SkyAPM/go2sky"
	"github.com/pkg/errors"
)

var (
	ValidateData = newValidateData()
)

func newValidateData() validateData {
	return validateData{
		RegistryItem: registryItem{
			Applications:    nil,
			InstanceMapping: nil,
		},
		SegmentItems: segmentItems{
			SegmentItems: nil,
		},
	}
}

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
	SegmentItems map[string]segmentItem `json:"segmentItems"`
}

func (s *segmentItems) AddSegmentItem(applicationId int32, segment segment) error {
	applicationCode, err := ValidateData.RegistryItem.findApplicationCode(applicationId)
	if err != nil {
		return err
	}
	item, ok := s.SegmentItems[applicationCode]
	if !ok {
		item = segmentItem{
			ApplicationCode: applicationCode,
			Segments:        nil,
		}
	}
	item.Segments = append(item.Segments, segment)
	return nil
}

type segmentItem struct {
	ApplicationCode string    `json:"applicationCode"`
	Segments        []segment `json:"segments"`
}

func NewSegment(spans []go2sky.ReportedSpan) (segment segment) {
	spanSize := len(spans)
	if spanSize < 1 {
		return
	}
	rootSpan := spans[spanSize-1]
	segment.SegmentId = rootSpan.Context().GetReadableGlobalTraceID()
	segment.Spans = make([]span, spanSize)

	for i, s := range spans {
		spanCtx := s.Context()
		segment.Spans[i] = span{
			OperationName: s.OperationName(),
			ParentSpanId:  spanCtx.ParentSpanID,
			SpanId:        spanCtx.SpanID,
			SpanLayer:     s.SpanLayer().String(),
			StartTime:     s.StartTime(),
			EndTime:       s.EndTime(),
			ComponentId:   s.ComponentID(),
			IsError:       s.IsError(),
			SpanType:      s.SpanType().String(),
			Peer:          s.Peer(),
		}
		// logs
		if len(s.Logs()) > 0 {
			var logs = make([]logEvent, len(s.Logs()))
			for _, log := range s.Logs() {
				var keyValues []keyValuePair
				for _, keyStringValue := range log.Data {
					keyValues = append(keyValues, keyValuePair{
						Key:   keyStringValue.Key,
						Value: keyStringValue.Value,
					})
				}
				logs = append(logs, logEvent{LogEvent: keyValues})
			}
			segment.Spans[i].Logs = logs
		}

		// tags
		if len(s.Tags()) > 0 {
			var tags = make([]keyValuePair, len(s.Tags()))
			for _, tag := range s.Tags() {
				tags = append(tags, keyValuePair{
					Key:   tag.Key,
					Value: tag.Value,
				})
			}
			segment.Spans[i].Tags = tags
		}

		// refs
		srr := make([]segmentRef, 0)
		if i == (spanSize-1) && spanCtx.ParentSpanID > -1 {
			srr = append(srr, segmentRef{
				ParentEndpointId:        0,
				ParentEndpoint:          "",
				NetworkAddressId:        0,
				EntryEndpointId:         0,
				RefType:                 "",
				ParentSpanId:            spanCtx.ParentSpanID,
				ParentTraceSegmentId:    spanCtx.ParentSegmentID,
				ParentServiceInstanceId: 0,
				NetworkAddress:          "",
				EntryEndpoint:           "",
				EntryServiceInstanceId:  0,
			})
		}
		if len(s.Refs()) > 0 {

		}
	}
	return
}

type segment struct {
	SegmentId string `json:"segmentId"`
	Spans     []span `json:"spans"`
}

type span struct {
	OperationName string         `json:"operationName"`
	OperationId   int32          `json:"operationId"`
	ParentSpanId  int32          `json:"parentSpanId"`
	SpanId        int32          `json:"spanId"`
	SpanLayer     string         `json:"spanLayer"`
	StartTime     int64          `json:"startTime"`
	EndTime       int64          `json:"endTime"`
	ComponentId   int32          `json:"componentId"`
	ComponentName string         `json:"componentName"`
	IsError       bool           `json:"isError"`
	SpanType      string         `json:"spanType"`
	Peer          string         `json:"Peer"`
	PeerId        int32          `json:"peerId"`
	Tags          []keyValuePair `json:"tags"`
	Logs          []logEvent     `json:"logs"`
	Refs          []segmentRef   `json:"refs"`
}

type segmentRef struct {
	ParentEndpointId        int32  `json:"parentEndpointId"`
	ParentEndpoint          string `json:"parentEndpoint"`
	NetworkAddressId        int32  `json:"networkAddressId"`
	EntryEndpointId         int32  `json:"entryEndpointId"`
	RefType                 string `json:"refType"`
	ParentSpanId            int32  `json:"parentSpanId"`
	ParentTraceSegmentId    string `json:"parentTraceSegmentId"`
	ParentServiceInstanceId int32  `json:"parentServiceInstanceId"`
	NetworkAddress          string `json:"networkAddress"`
	EntryEndpoint           string `json:"entryEndpoint"`
	EntryServiceInstanceId  int32  `json:"entryServiceInstanceId"`
}

type keyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type logEvent struct {
	LogEvent []keyValuePair `json:"logEvent"`
}
