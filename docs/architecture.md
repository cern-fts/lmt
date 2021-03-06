# Architecture

### On one end...
`LMT` listens to requests from `WebFTS` at `wss://hostname:port/socket`.

When a client (`WebFTS`) connects to `LMT` via a WSS (WebSocket Secure) connection, it sends the metadata for all the files it wishes to transfer via `FTS`.

For each file, `LMT` would then create an endpoint of the form `/transfer/delegationID/filename` and map it to that particular client.
It then informs the client (`WebFTS`) of the endpoint it created for each file via a JSON-based protocol message sent over the same WebSocket connection.

A transfer message might look something like the following:
```json
{
    "action":"transfer",
    "data":"https://lmt.cern.ch:8080/transfer/4d7dfd5d-f67a-461b-bc4e-20bf4a24c638/someFile.tar.gz"
}
```
The proxy maintains the WebSocket connection open, waiting for the File Transfer Service (`FTS`) to ask for the files.

`WebFTS` would receive those endpoint URLs, and proceed to submit a transfer job via `FTS`' REST API, with the source being the endpoint URLs it received from the proxy service.


### On the other end...
`LMT` will also be listening to incoming TCP connections at `https://hostname:port/transfer`.

Depending on what type of storage system the destination endpoint has, one of two things will happen:

a. `gsiftp` endpoint: `FTS` will be in charge of getting the data from the source, and streaming it to the destination.

b. `davs` or `dCache` endpoint: the destination will initiate a 3rd party pull request.

In both cases, the source for the transfer request will be an endpoint belonging to `LMT`, that is:


`https://lmt.cern.ch:8080/transfer/4d7dfd5d-f67a-461b-bc4e-20bf4a24c638/someFile.tar.gz`

When the `LMT` proxy receives a request asking for the files, it checks if:
1. The origin which sent the HTTP request has permissions to access the file. That is, if the identity of the X509 delegation certificate the request has is the same as the one the client (`WebFTS`) had when it registered the files to be transferred.
2. The files exist, and the client has not closed the WebSocket connection.


If those conditions are met, `LMT` would then tell `WebFTS` to start streaming the files via the open WebSocket connection, and it pipes the files contents to the response body of the GET request it received.


![alt text](diagram.png "LMT's Architecture")
