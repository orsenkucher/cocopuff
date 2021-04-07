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
* graphql
  > This service is an edge gateway. It resolves graphql queries by being a shim to other services via grpc.  
  > Go code is generated with `99designs/gqlgen`.
  * main.go
  * pb
    > `protoc` compiled Protobuf.
  * gql
    > `gqlgen` generated code.  
    >  `graphql-autobind` handwritten models.
  * gqlgen.yml
  * resolver
    * schema.go
    * resolvers.go
  * app.dockerfile
    > - stage 1: build `go.mod` dependencies.
    >   ```go
    >   go.mod
    >   go.sum
    >   go download
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
  * main.go
  * pb
  * app.dockerfile
  * db.dockerfile
  > ... and other services
* go.mod
* docker-compose.yaml
* .github
* .gitignore

