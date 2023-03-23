swag: swag-fmt
	swag init --parseDependency --parseInternal --parseDepth 1 -d cmd/server/,pkg/api/handlers/deliveries,pkg/api/handlers/useraccount,pkg/api/handlers/payments,pkg/api/handlers/useraccount,pkg/api/handlers/management,pkg/api/handlers/verification

swag-fmt:
	swag fmt -d cmd/server/,pkg/api/handlers/deliveries,pkg/api/handlers/useraccount,pkg/api/handlers/payments,pkg/api/handlers/useraccount,pkg/api/handlers/management,pkg/api/handlers/verification

release: push
	echo "pushing changes to heroku..."
	heroku container:release web -a portaldeliveries-backend

push:
	echo "deploying app in heroku..."
	heroku container:push web -a portaldeliveries-backend

build:
	docker compose up --build

down:
	docker compose down

webhook:
	stripe listen --forward-to localhost:8080/stripe/webhook