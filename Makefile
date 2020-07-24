build:
	docker build -t ayamaruyama/pcp:latest .
	docker push ayamaruyama/pcp:latest

build-v2:
	docker build -t ayamaruyama/pcp:latest-v2 .
	docker push ayamaruyama/pcp:latest-v2