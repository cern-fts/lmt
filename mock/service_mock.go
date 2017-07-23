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

package mock

import (
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

// GfalCopy executes a syscall to gfal-copy to transfer a file from src
// to dst.
func GfalCopy(dst string, src string) {
	log.Infof("gfal-copy %s %s", src, dst)
	binary, lookErr := exec.LookPath("gfal-copy")
	if lookErr != nil {
		panic(lookErr)
	}

	gfalCmd := exec.Command(binary, src, dst)
	gfalOut, err := gfalCmd.Output()
	if err != nil {
		panic(err)
	}
	log.Info(string(gfalOut))
}
