SRC_PATH=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

clean:
	rm -rf\
	 $(SRC_PATH)/dist\
	 $(SRC_PATH)/debug\
	 $(SRC_PATH)/godepgraph.png\
	 $(SRC_PATH)/*/cover.out\
	 $(SRC_PATH)/*/cover.html

generate_random_rsa:
	openssl req -x509 -newkey rsa:2048 -days 3650 -keyout rsa_private.pem -nodes -out rsa_cert.pem -subj "/CN=pablo-test-common-name"
	openssl ecparam -genkey -name prime256v1 -noout -out ec_private.pem
	openssl ec -in ec_private.pem -pubout -out ec_public.pem
	openssl rsa -in rsa_private.pem -pubout -out rsa_public.pem

