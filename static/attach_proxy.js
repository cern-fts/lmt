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

function AttachProxy(param) {
  var ws = new WebSocket(param.url);
  var file = param.file;
  var delegationID = Math.random().toString(36).substring(2, 15) +
    Math.random().toString(36).substring(2, 15);

  var filedata = { name: file.name, size: file.size, delegationID: delegationID };
  console.log(filedata)

  var slice_start = 0;
  var slice_end = filedata.size;
  var finished = false;
  var success = false;
  var error_messages = [];
  var endpoint;
  var shell = $("#shell");

  shell.WriteLine = function (line) {
    var newP;
    if (line[0] == '#') {
      newP = shell.append("<p class='comment'>" + line + "</p>");
    }
    else {
      newP = shell.append("<p>" + line + "</p>");
    }
    shell.scrollTop(shell.scrollTop() + newP.position().top + newP.height());
  }


  ws.onopen = function () {
    ws.send(JSON.stringify(filedata))
  };


  ws.onmessage = function (event) {
    // Parse controlMsg
    var controlMsg = JSON.parse(event.data);
    // Log message to shell.
    shell.WriteLine(JSON.stringify(controlMsg));
    console.log(event.data)


    if (controlMsg.action == "transfer") {
      endpoint = controlMsg.data;
    }
    // Last mile proxy is ready.
    if (controlMsg.action == "ready") {
      // Send file through websocket.
      ws.send(file.slice(slice_start, slice_end));
      ws.close()
      return;
    }
  };

  ws.onclose = function () {
    if (success) {
      return;
    }

    if (error_messages.length == 0) {
      error_messages[0] = { error: 'Unknown upload error' };
    }
    console.log(error_messages);
  }
}
