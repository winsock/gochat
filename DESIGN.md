# GoChat Design

## Security
Security was not a requirement for this project, if developed further I would look into using JWTs to secure and validate all requests against the API.

## Architecture
The project will leverage GoLang and SQLite to power the server. The client will communicate with the API via a RESTful API and optionally a websocket connection for realtime communications. This simple architecture and stack was chosen to allow rapid development and to ensure it is simple to run on any platform.

### Future Improvements
- Switch the database to Postgres or another RDBMS to allow scaling the service with a single database
  - Only a few lines would need to be changed, but left out of this project for simplicity in running.
- Allow multiple connections from the same user to the live websocket interface. For simplicity, I am only allowing a single client per user.
- At some point in scaling, this service may outgrow being solely backed by a RDBMS. This could be potentially resolved by having Kafka topics per user and running the messaging through there instead of the database.