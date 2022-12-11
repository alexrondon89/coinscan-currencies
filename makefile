build:
	docker build . --no-cache -t coinscan-currencies

run:
	docker run --name container-currencies -p 3001:80 -it coinscan-currencies

stop:
	docker stop container-currencies

push: #build
	docker tag coinscan-currencies alexrondon89/coinscan-currencies
	docker push alexrondon89/coinscan-currencies

run-from-hub:
	docker run -p 3001:3001 alexrondon89/coinscan-currencies

docker-compose:
	docker-compose up -d