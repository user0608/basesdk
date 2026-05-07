run:
	@TZ=UTC $(shell xenv) go run example/*.go -c example/config.yml

upgrade:
	@TZ=UTC go run example/*.go -c example/config.yml migrate up

downgrade:
	@TZ=UTC go run example/*.go -c example/config.yml migrate down

status:
	@TZ=UTC go run example/*.go -c example/config.yml migrate status

sql-scripts:
	@TZ=UTC go run example/*.go -c example/config.yml migrate script

new-migrate:
	@goose -dir=migrations create $(filter-out $@,$(MAKECMDGOALS)) sql

fmt-sql:
	@sqruff --config .sqruff fix migrations
