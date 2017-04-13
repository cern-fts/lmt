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

package bolt_test

import (
	"io/ioutil"
	"os"
	"time"

	"gitlab.cern.ch/fts/lmt/bolt"
)

// Now is the mocked current time for testing.
var Now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// Client is a test wrapper for bolt.Client.
type Client struct {
	*bolt.Client
}

// NewClient returns a new instance of Client pointing at a temporary file.
func NewClient() *Client {
	// Generate temporary filename.
	f, err := ioutil.TempFile("", "lmt-bolt-client-")
	if err != nil {
		panic(err)
	}
	f.Close()

	// Create client wrapper.
	c := &Client{
		Client: bolt.NewClient(),
	}
	c.Path = f.Name()
	c.Now = func() time.Time { return Now }

	return c
}

// MustOpenClient returns a new, open instance of Client.
func MustOpenClient() *Client {
	c := NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}
	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	defer os.Remove(c.Path)
	return c.Client.Close()
}
