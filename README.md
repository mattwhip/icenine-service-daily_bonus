# icenine-service-daily_bonus
Daily Bonus service for the IceNine system

This service is part of the IceNine project, a scalable cloud-based multiplayer server. The Daily Bonus service is built using Golang with the Buffalo web framework. It implements a random wheel spin award with increasing daily wheel award values. Here are some important pieces in the repository code structure:
- *actions*: HTTP request handlers (e.g. play daily bonus or get user status)
- *rpc*: RPC server code

Concepts/technologies used:
- Golang with Buffalo web framework for application creation
	- HTTP request handling
	- Object Relational Mapping for MySQL database access using Buffalo's Pop library
- gRPC for inter-service communication (e.g. login service retrieves daily bonus user status)
- Protobuf for all message (de)serialization
