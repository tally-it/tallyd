start-test-environment:
	@docker-compose up

stop-test-environment:
	@docker-compose down

generate-models:
	@echo "checking for dependencies"
	@which sqlboiler > /dev/null
	@echo "generating database model"
	@-rm repository/models/*
	@sqlboiler -o repository/models --no-tests --no-hooks mysql
	@rm repository/models/gorp_migrations.go
	@echo "done"

embed-migration-scripts:
	@echo "checking for dependencies"
	@which go-bindata > /dev/null
	@echo "embedding migration scripts"
	@go-bindata -prefix repository/migration -pkg migration -o repository/migration/migration.go repository/migration/resources
	@echo "done"

migrate-db:
	@go run repository/migration/cmd/main.go

update-db: embed-migration-scripts migrate-db generate-models

get-deps:
	go get -u -v github.com/jteeuwen/go-bindata/...
	go get -u -v github.com/vattle/sqlboiler/...
