# Golang's Context package examples

To run the server:
`go run . api-up`

The server will be started in port 3001

In the file "router.go" of the folder "api", you will find the endpoints.

```
GET /token: Endpoint to create a token

GET /value: Endpoint to run the context.WithValue example (decode the authorization token and return it as response)

GET /timeout: Endpoint to run the context.WithTimeout example (it will stop after 3 seconds, the duration variable is directly declared in the code)

GET /cancel: Endpoint to run the context.WithCancel example (it will no stop until you run the /close endpoint)

GET /close: Endpoint to stop the /cancel endpoint.
```

None of those endpoints receive any parameter or payload and only the /token endpoint doesn't have a Bearer Token protection.
