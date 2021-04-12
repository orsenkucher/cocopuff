## Repository
* api
  > The source of truth for our services and clients.
  * graphql
    > `GraphQL` is a language for our external API.  
    > We're using `.graphql` extension for our schemas as Facebook does.
    * schema.graphql
  * proto
    > `gRPC` is the way our services communicate each other, so `Protobuf` describes contract between them.
    * account.proto
      > ... and other services
* rfc: pkg, pub
  > Shared packages.
  * wg, wait, waitgroup
  * gs, grace, graceful shutdown
  * ec, errch, error channel
* graphql
  > This service is an edge gateway. It resolves graphql queries by being a shim to other services via grpc.  
  > Go code is generated with `99designs/gqlgen`.
  * cmd
    * main.go
  * pb
    > `protoc` compiled Protobuf.
  * gql
    > `gqlgen` generated code.  
    >  `graphql-autobind` handwritten models.
  * gqlgen.yml
  * resolver
    * schema.go
    * resolver.go
  * dataloader
    > batch and deduplicate queries.
  * client.go
    > with `go:generate` comments. As it uses grpc and is on surface.
  * cmd.dockerfile
    > - stage 1: build `go.mod` dependencies.
    >   ```go
    >   go.mod
    >   go.sum
    >   go mod download
    >   ```
    > - stage 2: build service.  
    >   ```dockerfile
    >   FROM deps AS build
    >   COPY . .
    >   ```
    > - stage 3: make little `buster-slim` image.
    >   ```dockerfile
    >   FROM build AS pure
    >   ```
* account
  * cmd
    > As main.go resides in `package main`, however I want to keep service logic inside `package account`.  
    * main.go
  * pb
  * log
  * env
  * migrations
  * server.go
    > with `go:generate` comments, as it closely uses grpc.
  * service.go
  * repository.go
  * cmd.dockerfile
  * db.dockerfile
  > ... and other services
* go.mod
* tools.go
  > `go build` ignores this file.   
  > `go.mod` will keep those imports and maintain versioning.  
  > We want the whole team to consistently generate the same code at any given point in time.
* Makefile
  > on windows as admin
  > ```ps
  > choco install make
  > ```
  >```bash
  > make download
  > make tools
  >```
* docker-compose.yaml
* .github
* .gitignore

