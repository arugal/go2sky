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
	"log"
)

func verify(excepted validateData, actual validateData) {
	log.Printf("excepted data: %v\n", excepted)
	log.Printf("actual data: %v\n", actual)
	// RegistryItemsAssert

	// SegmentItemsAssert
}

func registryItemsAssert(excepted registryItem, actual registryItem) ValidatorError {
	err := applicationAssert(excepted.Applications, actual.Applications)
	if err != nil {
		return err
	}
	log.Printf("registry applications assert successful.")
	err = instanceAssert(excepted.InstanceMapping, actual.InstanceMapping)
	if err != nil {
		return err
	}
	log.Printf("registry instances assert successful.")
	return nil
}

func applicationAssert(excepted map[string]string, actual map[string]string) ValidatorError {
	if excepted == nil {
		return nil
	}
	for exceptedCode, exceptedID := range excepted {
		actualID, err := getMatchApplication(actual, exceptedCode)
		if err != nil {
			return err
		}
		assert := AssertParse(exceptedID)
		assertErr := assert.AssertValue("registry application", actualID)
		if assertErr != nil {
			return &RegistryApplicationSizeNotEqualsError{
				ApplicationCode:  exceptedCode,
				ValueAssertError: *assertErr,
			}
		}
	}
	return nil
}

func getMatchApplication(actual map[string]string, exceptedCode string) (string, error) {
	for actualCode, actualID := range actual {
		if actualCode == exceptedCode {
			return actualID, nil
		}
	}
	return "", &RegistryApplicationNotFoundError{ApplicationCode: exceptedCode}
}

func instanceAssert(excepted map[string][]string, actual map[string][]string) ValidatorError {
	if excepted == nil {
		return nil
	}

	for exceptedCode, exceptedInstances := range excepted {
		actualInstances, err := getMatchApplicationInstance(actual, exceptedCode)
		if err != nil {
			return err
		}
		assert := AssertParse(fmt.Sprintf("%v", exceptedInstances))
		assertErr := assert.AssertValue(fmt.Sprintf("The registry instance of %s", exceptedCode), fmt.Sprintf("%v", actualInstances))
		if assertErr != nil {
			return &RegistryInstanceSizeNotEqualsError{
				ApplicationCode:  exceptedCode,
				ValueAssertError: *assertErr,
			}
		}
	}
	return nil
}

func getMatchApplicationInstance(actual map[string][]string, exceptedCode string) ([]string, error) {
	for actualCode, actualInstances := range actual {
		if actualCode == exceptedCode {
			return actualInstances, nil
		}
	}
	return nil, &RegistryInstanceOfApplicationNotFoundError{ApplicationCode: exceptedCode}
}
