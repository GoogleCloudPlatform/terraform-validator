// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfgcv

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// protoViaJSON uses JSON as an intermediary serialization to convert a value into
// a protobuf message.
func protoViaJSON(from interface{}, to proto.Message) error {
	jsn, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("marshaling to json: %s", err)
	}

	if err := jsonpb.UnmarshalString(string(jsn), to); err != nil {
		return fmt.Errorf("unmarshaling to proto: %s", err)
	}

	return nil
}
