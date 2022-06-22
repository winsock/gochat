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

## Testing
Unit tests have been written for both the API and Database access. Unit tests are run automatically via GitHub Actions, but they can be run locally too if you have go installed.
### Manually Run Unit Tests
```shell
go test -v ./...
```
### Manually Testing API
If you want to test the API manually via curl or another method, [API.md](API.md) has the reference for all API calls.
#### Steps
Generic steps, refer to [API.md](API.md) for example curl commands.
- Create two users using `/user/create`
- Start the demo live client for one or both of the users if wanting to test the realtime communications
- Send a message `/message/send`
- Search for the message with `/message/search`