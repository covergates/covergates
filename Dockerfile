FROM golang:alpine AS build

RUN apk --update add musl-dev
RUN apk --update add util-linux-dev
RUN apk --update add gcc g++
WORKDIR /go/src/github.com/covergates/covergates
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -v -o covergates ./cmd/server
FROM alpine
COPY --from=build /go/src/github.com/covergates/covergates/covergates /covergates
ENTRYPOINT [ "/covergates" ]
