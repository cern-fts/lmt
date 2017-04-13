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
	"time"

	"github.com/boltdb/bolt"
	"gitlab.cern.ch/fts/lmt"
)

// Client represents a client to the underlying BoltDB data store.
type Client struct {
	// Path to the BoltDB database.
	Path string
	// Returns the current time.
	Now func() time.Time

	// Services
	proxyService ProxyService

	db *bolt.DB
}

func NewClient() *Client {
	c := &Client{Now: time.Now}
	c.proxyService.client = c
	return c
}

// Initializes the BoltDB database.
func (c *Client) Open() error {
	// Open database file.
	db, err := bolt.Open(c.Path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	c.db = db

	// Initialize top-level buckets.
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte("LastMile")); err != nil {
		return err
	}

	return tx.Commit()
}

// Close closes the underlying BoltDB database.
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// ProxyService returns the proxy service associated with the client.
func (c *Client) ProxyService() lmt.ProxyService { return &c.proxyService }
