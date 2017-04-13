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

package bolt

import (
	"gitlab.cern.ch/fts/lmt"
	"gitlab.cern.ch/fts/lmt/bolt/internal"
)

// Ensure ProxyService implements lmt.ProxyService.
var _ lmt.ProxyService = &ProxyService{}

// ProxyService represents a service for managing transfer requests to the proxy.
type ProxyService struct {
	client *Client
}

// RegisterTransfer registers a new transfer.
func (s *ProxyService) RegisterTransfer(t *lmt.Transfer) error {
	// Require UserCN.
	if t.UserCN == "" {
		return lmt.ErrTransferUserCNRequired
	}
	// Require Token
	if t.Token == "" {
		return lmt.ErrTransferTokenRequired
	}

	// Start read-write transaction.
	tx, err := s.client.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify transfer does not already exist.
	b := tx.Bucket([]byte("LastMile"))
	if v := b.Get([]byte(t.UserCN)); v != nil {
		return lmt.ErrTransferExists
	}

	// Marshal and insert record into the database.
	if v, err := internal.MarshalTransfer(t); err != nil {
		return err
	} else if err := b.Put([]byte(t.UserCN), v); err != nil {
		return err
	}

	return tx.Commit()
}
