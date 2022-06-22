# Realtime Communications
This package implements the functionality needed to support realtime chat operations. This happens over a websocket connection that is initiated on the `/live` endpoint. An example client that connects and displays chat messages in realtime is provided in the `/demoLiveClient directory.

## Known Limitations
- Sending messages from the live client WIP
- Restricted to a single connection per user. In the real world this would need to support multiple connections from the same user.

## TODO
- Unit testing