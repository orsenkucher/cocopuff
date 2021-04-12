FROM deps AS build
WORKDIR /build
COPY --from=deps /deps .
COPY account account
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/account ./account/cmd

FROM debian:buster-slim AS bin
COPY --from=build /bin/account /bin/account
ENTRYPOINT [ "/bin/account" ]
