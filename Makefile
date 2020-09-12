

run-db:
	docker run -d --network=ac_test -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 --name db postgres:12.4