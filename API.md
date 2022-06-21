# API Documentation
By default, the server will bind on port 8080 on all interfaces. All RESTful API endpoints currently allow both GET or POST requests to allow easy use both via cURL or web browser. For GET the parameters should be passed as query parameters, and for POST as standard `application/x-www-form-urlencoded` data in the body.

## Time Format
All timestamps returns are RFC3339 compliant and include nanoseconds

## Response Content Type
All endpoints except for the `/live` endpoint return `application/json` 

## Users
### Create
Creates a new user, currently a user is just a `username` that is unique.
#### Endpoint
```
/user/create
```
#### Parameters
| Name     | Type   | Description                        | Required | Default Value |
|----------|--------|------------------------------------|----------|---------------|
| username | String | The name of the user to be created | Yes      | N/A           |
#### Response Object
Returns a JSON object containing the user UUID and their username
```json
{
  "uuid": "fb7fe917-07aa-4a87-8921-182d8bcfa2dd",
  "username": "example"
}
```
#### Response codes
- 500 If there was any error creating the user
- 201 On success
#### Example curl
```shell
curl -X POST -F 'username=example' http://localhost:8080/user/create
```
```json
{"uuid":"a0157386-4d85-43cd-b9f8-74feb5d88fdd","username":"example"}
```

## Messages
### Send
Sends a message from a specified user to a specified recipient
#### Endpoint
```
/message/send
```
#### Parameters
| Name      | Type   | Description                                         | Required | Default Value |
|-----------|--------|-----------------------------------------------------|----------|---------------|
| sender    | String | The name of the user that is sending the message    | Yes      | N/A           |
| recipient | String | The name of the user that is to receive the message | Yes      | N/A           |
| message   | String | The contents of the message to be sent              | Yes      | N/A           |
#### Response Object
Returns a JSON object containing the message that was sent and if the message was sent in realtime to the sender and/or recipient
```json
{
  "message": {
    "uuid": "a01a5462-a263-4540-a50f-a4a449edf992",
    "createdAt": "2022-06-21T16:28:07.778464-05:00",
    "content": "Hello World",
    "sender": {
      "uuid": "115728f0-e97c-47f8-b38a-c6eb313be3d9",
      "username": "sender"
    },
    "recipient": {
      "uuid": "4b395a8e-f335-48dc-a251-fe56412f5ac4",
      "username": "recipient"
    }
  },
  "sentLiveToSender": false,
  "sentLiveToRecipient": true
}
```
#### Response codes
- 400 If the message is empty
- 404 If either of the users are not found
- 500 If there was any error sending the message
- 201 On success
#### Example curl
```shell
# If you need to add users
curl -X POST -F 'username=sender' http://localhost:8080/user/create
curl -X POST -F 'username=recipient' http://localhost:8080/user/create
# Send the message
curl -X POST -F 'sender=sender' -F 'recipient=recipient' -F 'message=Hello World' http://localhost:8080/message/send
```
```json
{"message":{"uuid":"d6103304-ba02-4b92-8868-0fc218f16283","createdAt":"2022-06-21T16:39:00.743348-05:00","content":"Hello World","sender":{"uuid":"d19ead85-158b-4852-8c78-b9b6529b70a8","username":"sender"},"recipient":{"uuid":"9939cdd4-be17-45a6-b3c9-d95aebe6fefb","username":"recipient"}},"sentLiveToSender":false,"sentLiveToRecipient":false}
```
### Search
Searches for messages that match either the provided recipient or sender/recipient pair. Messages are in ascending order(the oldest first).
#### Endpoint
```
/message/search
```
#### Parameters
| Name      | Type              | Description                                                  | Required | Default Value       |
|-----------|-------------------|--------------------------------------------------------------|----------|---------------------|
| recipient | String            | The name of the user whose messages we are searching through | Yes      | N/A                 |
| sender    | String            | If provided search for messages only from this sender        | No       | N/A                 |
| from      | RFC3339 Timestamp | Limits how far back to search for messages                   | No       | 30 Days in the past |
| offset    | Unsigned Integer  | The offset to start from when returning data                 | No       | 0                   |
| limit     | Unsigned Integer  | The maximum number of results to be returned (max 1000)      | No       | 100                 |
#### Response Object
Returns a JSON object that contains the list of messages found(if any) and pagination metadata
```json
{
  "messages": [
    {
      "uuid": "fe44f845-1cc9-4573-bf7f-e77e81c56694",
      "createdAt": "2022-06-21T16:50:56.889605-05:00",
      "content": "Hello World 2",
      "sender":{
        "uuid": "d19ead85-158b-4852-8c78-b9b6529b70a8",
        "username": "sender"
      },
      "recipient":{
        "uuid": "9939cdd4-be17-45a6-b3c9-d95aebe6fefb",
        "username": "recipient"
      }
    }
  ],
  "limit": 100,
  "offset": 0,
  "count": 1
}
```
#### Response Codes
- 400 On one or any of the following conditions:
  - An invalid formatted timestamp in `from`
  - An invalid value passed in `offset`
  - An invalid value passed in `limit`
  - Missing `recipient`
- 500 If an error occurred while searching for the messages
- 200 If the search was performed successfully
#### Example curl
##### Create example data
```shell
# Adding users
curl -X POST -F 'username=sender1' http://localhost:8080/user/create
curl -X POST -F 'username=sender2' http://localhost:8080/user/create
curl -X POST -F 'username=recipient' http://localhost:8080/user/create
# Send some messages message
curl -X POST -F 'sender=sender1' -F 'recipient=recipient' -F 'message=Hello World From 1' http://localhost:8080/message/send
curl -X POST -F 'sender=sender2' -F 'recipient=recipient' -F 'message=Hello World From 2' http://localhost:8080/message/send
```
##### All messages for recipient
```shell
curl -X POST -F 'recipient=recipient' http://localhost:8080/message/search
```
##### All messages for recipient from sender1
```shell
curl -X POST -F 'recipient=recipient' -F 'sender=sender1' http://localhost:8080/message/search
```