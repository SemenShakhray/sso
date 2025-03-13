YAML_FILE=config/local.yaml

### Run app ###
run-app:
	go run cmd/sso/main.go --config=config/local.yaml

MIGRATION_DIR:=$(shell yq '.migration_dir' $(YAML_FILE))
STORAGE_PATH:=$(shell yq '.storage_path' $(YAML_FILE))

migration-add:
	goose -dir ${MIGRATION_DIR} create $(name) sql

migration-up:
	goose -dir ${MIGRATION_DIR} sqlite3 ${STORAGE_PATH} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} sqlite3 ${STORAGE_PATH} down -v	

migration-status:
	goose -dir ${MIGRATION_DIR} sqlite3 ${STORAGE_PATH} status