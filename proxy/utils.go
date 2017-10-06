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

package proxy

import (
	"crypto/rand"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
	voms "gitlab.cern.ch/flutter/go-proxy"
)

// NewUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8],
		uuid[8:10], uuid[10:]), nil
}

// X509Identity parses an HTTP request in order to extract the X509 proxy
// certificate's identity.
func X509Identity(req *http.Request) (pkix.Name, error) {
	var identity pkix.Name
	var err error
	x509 := voms.X509Proxy{}

	if req.TLS.PeerCertificates == nil {
		err = errors.New(errProxyCertRequired)
		log.WithFields(logrus.Fields{
			"event": "no_x505_proxy_cert",
		}).Error(err)
	} else {
		if err = x509.InitFromCertificates(req.TLS.PeerCertificates); err != nil {
			log.Error(err)
			log.WithFields(logrus.Fields{
				"event": "x509_proxy_cert_init_error",
			}).Error(err)
		}
		identity = x509.Identity
	}
	return identity, err
}

// X509DelegationID parses an HTTP request in order to extract the X509 proxy
// certificate's delegation ID.
func X509DelegationID(req *http.Request) (string, error) {
	var delegationID string
	var err error
	x509 := voms.X509Proxy{}

	if req.TLS.PeerCertificates == nil {
		err = errors.New(errProxyCertRequired)
		log.WithFields(logrus.Fields{
			"event": "no_x505_proxy_cert",
		}).Error(err)
	} else {
		if err = x509.InitFromCertificates(req.TLS.PeerCertificates); err != nil {
			log.Error(err)
			log.WithFields(logrus.Fields{
				"event": "x509_proxy_cert_init_error",
			}).Error(err)
		}
		delegationID = x509.DelegationID()
	}
	return delegationID, err
}

// CheckIdentity checks if the identity of the FTS job matches the one
// that was provided by the client.
func CheckIdentity(transferID, identity string) bool {
	return Transfers[transferID].identity == identity
}

// TransferID concatenates the user's delegation ID and
// the filename the user's wish to transfer to create
// the transfer endpoint.
func TransferID(delegationID, filename string) string {
	return fmt.Sprintf("%s/%s", delegationID, filename)
}
