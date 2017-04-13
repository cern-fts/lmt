/*
 * Copyright (c) CERN 2017
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package internal_test

import (
	"reflect"
	"testing"
	"time"

	"gitlab.cern.ch/fts/lmt"
	"gitlab.cern.ch/fts/lmt/bolt/internal"
)

// Ensure dial can be marshaled and unmarshaled.
func TestMarshalTransfer(t *testing.T) {
	v := lmt.Transfer{
		UserCN:     "ID",
		Token:      "TOKEN",
		Name:       "Some Name",
		Filepath:   "/home/user/somefile",
		SubmitTime: time.Now().UTC(),
	}

	var other lmt.Transfer
	if buf, err := internal.MarshalTransfer(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.UnmarshalTransfer(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: %#v", other)
	}
}
