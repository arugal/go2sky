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
	"strconv"
	"strings"
)

type Assert interface {
	AssertValue(desc string, actualValue string) *ValueAssertError
}

type elementAssert struct {
	exceptedValue string
}

type EqualsAssert struct {
	elementAssert
}

func (e *EqualsAssert) AssertValue(desc string, actualValue string) *ValueAssertError {
	if e.exceptedValue != actualValue {
		return valueAssertError(e.exceptedValue, actualValue, desc)
	}
	return nil
}

type GreatThanAssert struct {
	elementAssert
}

func (g *GreatThanAssert) AssertValue(desc string, actualValue string) *ValueAssertError {
	excepted := parseInt(g.exceptedValue)
	actual := parseInt(actualValue)

	if actual <= excepted {
		return valueAssertError("gt "+g.exceptedValue, actualValue, desc)
	}
	return nil
}

type GreetEqualAssert struct {
	elementAssert
}

func (g *GreetEqualAssert) AssertValue(desc string, actualValue string) *ValueAssertError {
	excepted := parseInt(g.exceptedValue)
	actual := parseInt(actualValue)

	if actual < excepted {
		return valueAssertError("ge "+g.exceptedValue, actualValue, desc)
	}
	return nil
}

type NoopAssert struct {
}

func (n *NoopAssert) AssertValue(desc string, actualValue string) *ValueAssertError {
	return nil
}

type NotEqualsAssert struct {
	elementAssert
}

func (n *NotEqualsAssert) AssertValue(desc string, actualValue string) *ValueAssertError {
	if n.exceptedValue == actualValue {
		return valueAssertError("not eq "+n.exceptedValue, actualValue, desc)
	}
	return nil
}

type NotNullAssert struct {
}

func (n *NotNullAssert) AssertValue(desc string, actualValue string) *ValueAssertError {
	if actualValue == "" {
		return valueAssertError("not null", actualValue, desc)
	}
	return nil
}

type NullAssert struct {
}

func (n *NullAssert) AssertValue(desc string, actualValue string) *ValueAssertError {
	if actualValue != "" {
		return valueAssertError("null", actualValue, desc)
	}
	return nil
}

func parseInt(str string) int64 {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return value
}

func valueAssertError(exceptedValue string, actualValue string, desc string) *ValueAssertError {
	return &ValueAssertError{
		Desc:          desc,
		ExceptedValue: exceptedValue,
		ActualValue:   actualValue,
	}
}

func AssertParse(express string) Assert {
	if express == "" {
		return &NoopAssert{}
	}

	if express == "not null" {
		return &NotNullAssert{}
	}

	if express == "null" {
		return &NullAssert{}
	}

	expressSegment := strings.Split(express, " ")
	if len(expressSegment) == 1 {
		return &EqualsAssert{elementAssert{exceptedValue: expressSegment[0]}}
	} else if len(expressSegment) == 2 {
		switch expressSegment[0] {
		case "nq":
			return &NotEqualsAssert{elementAssert{exceptedValue: expressSegment[1]}}
		case "eq":
			return &EqualsAssert{elementAssert{exceptedValue: expressSegment[1]}}
		case "gt":
			return &GreatThanAssert{elementAssert{exceptedValue: expressSegment[1]}}
		case "ge":
			return &GreetEqualAssert{elementAssert{exceptedValue: expressSegment[1]}}
		}
	}

	return &EqualsAssert{elementAssert{exceptedValue: express}}
}
