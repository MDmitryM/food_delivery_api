.PHONY: proto img img_push

IMAGE_NAME = mdmitrym/food_delivery_api
TAG ?= latest

#собрать grpc из /proto/...
proto:
	protoc -I proto --go_out=src --go-grpc_out=src proto/api.proto

img:
	docker build -t ${IMAGE_NAME}:${TAG} .

img_push:
	docker push ${IMAGE_NAME}:${TAG}