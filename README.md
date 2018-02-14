# abacus
Abacus Registration application

### Setup Following Enviroinment Variables before starting the server.
* MAIL_API_KEY
* DBHOST
* DBUSER
* DBPASS
* DBNAME

## Build Binary

### For host architecture
`go build *.go`

### For different target architecture
Setup enviroinment variables `GOARCH=386,GOOS=linux` and run `go build *.go`

## Starting Server
`./app -port=3000`

port flag can be left and defaults to 8080
