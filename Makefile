BUILDER=./.builder
RULES=go
include $(BUILDER)/rules.mk
$(BUILDER)/rules.mk:
	-go run endobit.io/builder@latest init

format::
	buf format -w

lint::
	buf lint
	go tool github.com/sqlc-dev/sqlc/cmd/sqlc compile

generate::
	buf generate
	go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate

build::
	CGO_ENABLED=0 $(GO_BUILD) -o metald ./services/metal


.PHONY: docker-volumes
docker-volumes: ## creates the cert/data docker volumes (do this once)
	docker volume create metal-certs || true
	docker volume create metal-data || true

clean::
	rm -rf stackd

nuke::
	rm -rf gen
	rm -rf internal/data/db



