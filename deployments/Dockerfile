FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go test ./...
RUN go build -o /out/auto-recruiter ./cmd/app

FROM alpine:3.20
RUN adduser -D -H app
USER app
COPY --from=build /out/auto-recruiter /usr/local/bin/auto-recruiter
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/auto-recruiter"]