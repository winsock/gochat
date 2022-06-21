# gochat
A simple messaging server created for Guild Education.

## Requirements
- go - https://go.dev/doc/install
  - Optional if only running via docker and do not want to run the demo live client
  - For macOS with Homebrew `brew install go` will set up everything
- docker
  - Optional, only if you want to run the server in a Docker container

## Running
All commands assumed to be run from the root of the project
### Server
#### Manually
```shell
go run
```
#### Docker
```shell
docker build -t gochat-server .
docker run -dp 8080:8080 gochat-server
```
### Demo Live Client
You will need to create a user via the API before the client will work
```shell
go run ./demoLiveClient -username=<your username>
```