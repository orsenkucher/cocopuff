FROM deps AS build
WORKDIR /build
COPY --from=deps /deps .
COPY graphql graphql
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/graphql ./graphql/cmd

FROM debian:buster-slim AS bin
COPY --from=build /bin/graphql /bin/graphql
ENTRYPOINT [ "/bin/graphql" ]
