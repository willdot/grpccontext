Playing around with the Go context and seeing what's possible in terms of passing it around using gRPC calls.

* Adding the source to the context as a header so that the server can see who sent the request.

* Adding a timeout to the context so that the client abandons the request. The server will catch this context being cancelled and handles it.
