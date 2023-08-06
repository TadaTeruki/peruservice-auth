
.PHONY: generate_key_pair
generate_key_pair:
	mkdir -p .secret
	openssl genrsa -out .secret/id_rsa 4096
	openssl rsa -in .secret/id_rsa -pubout -out .secret/id_rsa.pub