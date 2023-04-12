protoc --proto_path=. *.proto --go_out==plugins=grpc,paths=source_relative:.

protoc --go_out=. --go_opt=paths=source_relative --micro_out=. --micro_opt=paths=source_relative *.proto