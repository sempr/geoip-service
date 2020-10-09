FROM golang:1.15.2 as build

WORKDIR /app
RUN go get -u github.com/gobuffalo/packr/v2/packr2
COPY cmd /app/cmd
COPY data /app/data
COPY go.mod go.sum /app/

RUN go generate mygeo/cmd/web
RUN go build -o /app1 ./cmd/web/

FROM debian:10-slim
COPY --from=build /app1 /app
CMD /app

