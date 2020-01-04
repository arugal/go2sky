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

package convert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SkyAPM/go2sky/reporter/grpc/common"
)

func GlobalIDConvertString(id []int64) string {
	ii := make([]string, len(id))
	for i, v := range id {
		ii[i] = fmt.Sprint(v)
	}
	return strings.Join(ii, ".")
}

func RefTypeConvertString(refType common.RefType) string {
	return common.RefType_name[int32(refType)]
}

func SpanLayerConvertString(spanLayer common.SpanLayer) string {
	return common.SpanLayer_name[int32(spanLayer)]
}

func SpanTypeConvertString(spanType common.SpanType) string {
	return common.SpanType_name[int32(spanType)]
}

func Int32ConvertString(value int32) string {
	return Int64ConvertString(int64(value))
}

func Int64ConvertString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func BoolConvertString(bool2 bool) string {
	return strconv.FormatBool(bool2)
}
