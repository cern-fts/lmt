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

package main

import (
	"time"
)

// Transfer represents a transfer request made by an authorized CERN user.
// A user/device token is used to authenticate requests.
type Transfer struct {
	UserCN       string    `json:"userCN"`
	Token        string    `json:"token"`
	Name         string    `json:"name"`
	Filepath     string    `json:"filepath"`
	SubmitTime   time.Time `json:"submitTime"`
	Origin       string
	ListenerAddr string
	Destination  string
}

// Client creates a connection to the services.
type Client interface {
	ProxyService() ProxyService
}

// ProxyService represents a service for managing transfer requests to the
// proxy.
type ProxyService interface {
	RegisterTransfer(transfer *Transfer) error
}
