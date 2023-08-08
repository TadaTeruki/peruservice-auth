include .env

.PHONY: generate_key_pair
generate_key_pair:
	mkdir -p .secret
	openssl genrsa -out .secret/id_rsa 4096
	openssl rsa -in .secret/id_rsa -pubout -out .secret/id_rsa.pub

.PHONY: create-db
create-db:
	mkdir -p postgres/data

.PHONY: delete-db
delete-db:
	rm -r postgres/data 

.PHONY: login-db
login-db:
	docker exec -it auth-db psql -U ${DB_USER} -d ${DB_NAME}

.PHONY: insert-admin-for-debug
insert-admin-for-debug:
	docker exec -it auth-db psql -U ${DB_USER} -d ${DB_NAME} -c "INSERT INTO admin (id, password) VALUES ('hoge@example.com', 'password')"