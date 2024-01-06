.PHONY: purify_db

asseble_proto:	/usr/bin/sudo PATH=$GOPATH/bin /usr/local/bin/protoc --proto_path=protoc\
	--go_out=protoc --go_opt=paths=source_relative --go-grpc_out=protoc\
	--go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional protoc/*.proto

start_evans: PATH=$GOPATH/bin evans --proto protoc/chat.proto repl -p 8081

install_package: PATH=$GOPATH/bin /usr/local/go/bin/go install github.com/ktr0731/evans@v0.10.11

new_migration: migrate create -ext sql -dir dir/migrations -seq user

up_migration: migrate -path migration -database "postgresql://test:password@localhost:5432/auth_db?sslmode=disable" -verbose up

down_migration: migrate -path migration -database "postgresql://test:password@localhost:5432/auth_db?sslmode=disable" -verbose down

connect_db: docker exec -it $(docker ps --filter status=running --filter ancestor=postgre:14 --format "{{.ID}}") psql --username=test

purify_db:
	docker exec -it $(docker ps --filter ancestor=postgres:14 --format "{{.ID}}") psql -U test -d auth_db --command="update schema_migrations set dirty=true where 'version' like '%'"