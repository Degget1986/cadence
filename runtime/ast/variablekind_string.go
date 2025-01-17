/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated by "stringer -type=VariableKind"; DO NOT EDIT.

package ast

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[VariableKindNotSpecified-0]
	_ = x[VariableKindVariable-1]
	_ = x[VariableKindConstant-2]
}

const _VariableKind_name = "VariableKindNotSpecifiedVariableKindVariableVariableKindConstant"

var _VariableKind_index = [...]uint8{0, 24, 44, 64}

func (i VariableKind) String() string {
	if i < 0 || i >= VariableKind(len(_VariableKind_index)-1) {
		return "VariableKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _VariableKind_name[_VariableKind_index[i]:_VariableKind_index[i+1]]
}
