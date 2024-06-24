gen: 
	protoc \
	--go_out=services/common --go_opt=paths=source_relative \
	--go-grpc_out=services/common --go-grpc_opt=paths=source_relative \
	protobuf/*.proto
