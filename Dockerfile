FROM golang:1.24 AS build

ENV GOPATH=/
WORKDIR /src/
COPY . .

RUN go mod download; CGO_ENABLED=0 go build -o /food_delivery_api ./cmd/main.go

FROM alpine:3.17
COPY --from=build /food_delivery_api /food_delivery_api

EXPOSE 8080
CMD [ "./food_delivery_api" ]