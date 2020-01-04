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

import "fmt"

type ValidatorError interface {
	Error() string
}

type ValueAssertError struct {
	Desc          string
	ExceptedValue string
	ActualValue   string
}

func (e *ValueAssertError) Error() string {
	return fmt.Sprintf("[%s]: excepted=>{%s}, actual=>{%s}", e.Desc, e.ExceptedValue, e.ActualValue)
}

type RegistryApplicationNotFoundError struct {
	ValidatorError
	ApplicationCode string
}

func (r *RegistryApplicationNotFoundError) Error() string {
	return fmt.Sprintf("RegistryApplicationNotFoundError\nexpected: %s\nactual: Not Found\n", r.ApplicationCode)
}

type RegistryApplicationSizeNotEqualsError struct {
	ValidatorError
	ApplicationCode  string
	ValueAssertError ValueAssertError
}

func (r *RegistryApplicationSizeNotEqualsError) Error() string {
	return fmt.Sprintf("RegistryApplicationSizeNotEqualsError %s\nexpected: %s\nactual: %s\n",
		r.ApplicationCode, r.ValueAssertError.ExceptedValue, r.ValueAssertError.ActualValue)
}

type RegistryInstanceOfApplicationNotFoundError struct {
	ValidatorError
	ApplicationCode string
}

func (r *RegistryInstanceOfApplicationNotFoundError) Error() string {
	return fmt.Sprintf("RegistryInstanceOfApplicationNotFoundError\nexpected: Instance of Service(%s)\nactual: Not Found\n", r.ApplicationCode)
}

type RegistryInstanceSizeNotEqualsError struct {
	ValidatorError
	ApplicationCode  string
	ValueAssertError ValueAssertError
}

func (r *RegistryInstanceSizeNotEqualsError) Error() string {
	return fmt.Sprintf("RegistryInstanceSizeNotEqualsError %s\nexpected: %s\nactual: %s\n",
		r.ApplicationCode, r.ValueAssertError.ExceptedValue, r.ValueAssertError.ActualValue)
}
