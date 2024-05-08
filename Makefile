scheduler:
	go run ./services/scheduler/service/

worker:
	PORT=8081 go run ./services/worker/service/ 

start:
	make worker & make scheduler

seed:
	go run ./scripts/seeder/

srcWorker=/home/razvan/Desktop/facultate/master/a1/sem2/da/service-provider/services/worker/proto
dstWorker=/home/razvan/Desktop/facultate/master/a1/sem2/da/service-provider/services/worker/proto

genWorkerRPC:
	protoc --proto_path=$(srcWorker) \
	--go_out=$(dstWorker) --go_opt=paths=source_relative api.proto \
	--go-grpc_out=$(dstWorker) --go-grpc_opt=paths=source_relative \
	api.proto

dockerWorker:
	docker build -t worker -f services/worker/Dockerfile --no-cache  --progress=plain . 

firstWorker:
	docker run --name first -p 8081:8081 -e PORT="8081" -d  worker 

secondWorker:
	docker run --name second -p 8082:8082 -e PORT="8082" -d  worker 

thirdWorker:
	docker run --name third -p 8083:8083 -e PORT="8083" -d  worker 

buildWorker:
	go build -C services/worker/service -o ../out/app

