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

var inputElement = document.getElementById("uploadFile");
inputElement.addEventListener("change", handleFiles, false);

function handleFiles() {
  var selectedFile = this.files[0];
  var params = new Object();
  params.url = "wss://lmt.cern.ch:8080/socket";
  params.file = selectedFile;
  AttachProxy(params);
}
