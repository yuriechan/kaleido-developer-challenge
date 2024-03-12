.PHONY: docker-compose-up
docker-compose-up:
	docker-compose build && docker compose up

.PHONY: docker-compose-down
docker-compose-down:
	rm -r db_data && docker compose down -v

.PHONY: solc 
solc:
	solc --evm-version paris --bin --abi --optimize --overwrite -o contracts/ contracts/marketplace.sol
