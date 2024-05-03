scheduler:
	go run ./services/scheduler/service/

worker:
	go run ./services/worker/service/

seed:
	go run ./scripts/seeder/

srcWorker=/home/razvan/Desktop/facultate/master/a1/sem2/da/service-provider/services/worker/proto
dstWorker=/home/razvan/Desktop/facultate/master/a1/sem2/da/service-provider/services/worker/proto

genWorkerRPC:
	protoc --proto_path=$(srcWorker) \
	--go_out=$(dstWorker) --go_opt=paths=source_relative api.proto \
	--go-grpc_out=$(dstWorker) --go-grpc_opt=paths=source_relative \
	api.proto
