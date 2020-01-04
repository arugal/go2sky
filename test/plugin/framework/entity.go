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
	"github.com/SkyAPM/go2sky/reporter/grpc/common"
	"github.com/SkyAPM/go2sky/test/plugin/util/convert"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	validateDataInstance = newValidateData()
)

func newValidateData() validateData {
	return validateData{
		RegistryItem: registryItem{
			Applications:    make(map[string]string),
			InstanceMapping: make(map[string][]string),
		},
		SegmentItems: segmentItems{
			SegmentItems: map[string]*segmentItem{},
		},
	}
}

type validateData struct {
	RegistryItem registryItem `yaml:"registryItem"`
	SegmentItems segmentItems `yaml:"segmentItems"`
}

type registryItem struct {
	Applications    map[string]string   `yaml:"applications"`
	InstanceMapping map[string][]string `yaml:"instances"`
}

func (r *registryItem) findApplicationCode(applicationID int32) (string, error) {
	applicationIDStr := convert.Int32ConvertString(applicationID)
	for k, v := range r.Applications {
		if v == applicationIDStr {
			return k, nil
		}
	}
	return "", errors.Errorf("Cannot found the code of applicationID[ %d ].", applicationID)
}

type segmentItems struct {
	SegmentItems map[string]*segmentItem `yaml:"segmentItems"`
}

func (s *segmentItems) addSegmentItem(applicationID int32, segment segment) error {
	applicationCode, err := validateDataInstance.RegistryItem.findApplicationCode(applicationID)
	if err != nil {
		return err
	}
	item, ok := s.SegmentItems[applicationCode]
	if !ok {
		item = &segmentItem{
			ApplicationCode: applicationCode,
		}
		s.SegmentItems[applicationCode] = item
	}
	item.Segments = append(item.Segments, segment)
	return nil
}

type segmentItem struct {
	ApplicationCode string    `yaml:"applicationCode"`
	Segments        []segment `yaml:"segments"`
}

type segment struct {
	SegmentID string `yaml:"segmentId"`
	Spans     []span `yaml:"spans"`
}

type span struct {
	OperationName string         `yaml:"operationName"`
	OperationID   string         `yaml:"operationId"`
	ParentSpanID  string         `yaml:"parentSpanId"`
	SpanID        string         `yaml:"spanId"`
	ComponentID   string         `yaml:"componentId"`
	PeerID        string         `yaml:"peerId"`
	IsError       string         `yaml:"isError"`
	SpanLayer     string         `yaml:"spanLayer"`
	StartTime     string         `yaml:"startTime"`
	EndTime       string         `yaml:"endTime"`
	ComponentName string         `yaml:"componentName"`
	SpanType      string         `yaml:"spanType"`
	Peer          string         `yaml:"Peer"`
	Tags          []keyValuePair `yaml:"tags"`
	Logs          []logEvent     `yaml:"logs"`
	Refs          []segmentRef   `yaml:"refs"`
}

type segmentRef struct {
	ParentEndpointID        string `yaml:"parentEndpointId"`
	ParentSpanID            string `yaml:"parentSpanId"`
	NetworkAddressID        string `yaml:"networkAddressId"`
	EntryEndpointID         string `yaml:"entryEndpointId"`
	ParentServiceInstanceID string `yaml:"parentServiceInstanceId"`
	EntryServiceInstanceID  string `yaml:"entryServiceInstanceId"`
	ParentEndpoint          string `yaml:"parentEndpoint"`
	RefType                 string `yaml:"refType"`
	ParentTraceSegmentID    string `yaml:"parentTraceSegmentId"`
	NetworkAddress          string `yaml:"networkAddress"`
	EntryEndpoint           string `yaml:"entryEndpoint"`
}

type keyValuePair struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type logEvent struct {
	LogEvent []keyValuePair `yaml:"logEvent"`
}

func registryApplication(applicationCode string, applicationID int32) {
	if _, ok := validateDataInstance.RegistryItem.Applications[applicationCode]; !ok {
		validateDataInstance.RegistryItem.Applications[applicationCode] = convert.Int32ConvertString(applicationID)
	}
}

func registryInstance(applicationID int32, instanceID int32) error {
	applicationCode, err := validateDataInstance.RegistryItem.findApplicationCode(applicationID)
	if err != nil {
		return err
	}

	instances := validateDataInstance.RegistryItem.InstanceMapping[applicationCode]
	validateDataInstance.RegistryItem.InstanceMapping[applicationCode] = append(instances, convert.Int32ConvertString(instanceID))
	return nil
}

func addSpans(applicationID int32, instanceID int32, spans []go2sky.ReportedSpan) error {
	return validateDataInstance.SegmentItems.addSegmentItem(applicationID, newSegment(instanceID, spans))
}

func newSegment(instanceID int32, spans []go2sky.ReportedSpan) (segment segment) {
	spanSize := len(spans)
	if spanSize < 1 {
		return
	}
	rootSpan := spans[spanSize-1]
	segment.SegmentID = convert.GlobalIDConvertString(rootSpan.Context().TraceID)
	segment.Spans = make([]span, spanSize)

	for i, s := range spans {
		spanCtx := s.Context()
		segment.Spans[i] = span{
			OperationName: s.OperationName(),
			ParentSpanID:  convert.Int32ConvertString(spanCtx.ParentSpanID),
			SpanID:        convert.Int32ConvertString(spanCtx.SpanID),
			SpanLayer:     convert.SpanLayerConvertString(s.SpanLayer()),
			StartTime:     convert.Int64ConvertString(s.StartTime()),
			EndTime:       convert.Int64ConvertString(s.EndTime()),
			ComponentID:   convert.Int32ConvertString(s.ComponentID()),
			IsError:       convert.BoolConvertString(s.IsError()),
			SpanType:      convert.SpanTypeConvertString(s.SpanType()),
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
			var tags []keyValuePair
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
				RefType:                 convert.RefTypeConvertString(common.RefType_CrossThread),
				ParentSpanID:            convert.Int32ConvertString(spanCtx.ParentSpanID),
				ParentTraceSegmentID:    convert.GlobalIDConvertString(spanCtx.ParentSegmentID),
				ParentServiceInstanceID: convert.Int32ConvertString(instanceID),
			})
		}
		if len(s.Refs()) > 0 {
			for _, tc := range s.Refs() {
				srr = append(srr, segmentRef{
					ParentSpanID:            convert.Int32ConvertString(tc.ParentSpanID),
					ParentTraceSegmentID:    convert.GlobalIDConvertString(tc.ParentSegmentID),
					ParentServiceInstanceID: convert.Int32ConvertString(tc.ParentServiceInstanceID),
					EntryEndpoint:           tc.EntryEndpoint,
					EntryEndpointID:         convert.Int32ConvertString(tc.EntryEndpointID),
					EntryServiceInstanceID:  convert.Int32ConvertString(tc.EntryServiceInstanceID),
					NetworkAddress:          tc.NetworkAddress,
					NetworkAddressID:        convert.Int32ConvertString(tc.NetworkAddressID),
					ParentEndpoint:          tc.ParentEndpoint,
					ParentEndpointID:        convert.Int32ConvertString(tc.ParentEndpointID),
					RefType:                 convert.RefTypeConvertString(common.RefType_CrossProcess),
				})
			}
			segment.Spans[i].Refs = srr
		}
	}
	return
}

func receiveData() string {
	data, _ := yaml.Marshal(validateDataInstance)
	return string(data)
}
