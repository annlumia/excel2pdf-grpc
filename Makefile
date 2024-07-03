protoc:
	protoc --proto_path=. --go-grpc_out=paths=source_relative:. --go_out=:. proto/*.proto


build-server:
	# Buil server with no console
	GOOS=windows go build -ldflags -H=windowsgui .