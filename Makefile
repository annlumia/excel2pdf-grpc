protoc:
	protoc --proto_path=. --go-grpc_out=paths=source_relative:. --go_out=:. proto/*.proto