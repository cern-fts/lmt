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

package internal

import (
	"github.com/golang/protobuf/proto"
	"gitlab.cern.ch/fts/lmt"
	"time"
)

//go:generate protoc --gogo_out=. internal.proto

// MarshalTransfer encodes a TransferRequest to binary format.
func MarshalTransfer(t *lmt.Transfer) ([]byte, error) {
	return proto.Marshal(&Transfer{
		UserCN:       t.UserCN,
		Token:        t.Token,
		Name:         t.Name,
		Filepath:     t.Filepath,
		SubmitTime:   t.SubmitTime.UnixNano(),
		Origin:       t.Origin,
		ListenerAddr: t.ListenerAddr,
		Destination:  t.Destination,
	})
}

// UnmarshalTransfer decodes a Trnsfer from binary data.
func UnmarshalTransfer(data []byte, t *lmt.Transfer) error {
	var pb Transfer
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	t.UserCN = pb.UserCN
	t.Token = pb.Token
	t.Name = pb.Name
	t.Filepath = pb.Filepath
	t.SubmitTime = time.Unix(0, pb.SubmitTime).UTC()
	t.Origin = pb.Origin
	t.ListenerAddr = pb.ListenerAddr
	t.Destination = pb.Destination

	return nil
}
