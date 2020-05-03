FROM golang:1.11-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/newsfetcher

FROM scratch
COPY --from=build /bin/newsfetcher /bin/newsfetcher
ENTRYPOINT ["/bin/newsfetcher"]
