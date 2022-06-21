# GoChat Design

## Security
Security was not a requirement for this project, if developed further I would look into using JWTs to secure and validate all requests against the API.

## Architecture
The project will leverage GoLang and SQLite(via the gorm framework). The client will communicate with the API via a RESTful API and optionally a websocket connection for realtime communications. This simple architecture and stack was chosen to allow rapid development and to ensure it is simple to run on any platform.

### Future Improvements
- Switch the database to Postgres or another RDBMS to allow scaling the service with a single database
  - Only a few lines would need to be changed, but left out of this project for simplicity in running.