v1:
	docker build -f Dockerfile1 -t roads001/httpserver:v1 .	
pushv1:
	docker push roads001/httpserver:v1
v2:
	docker build -f Dockerfile2 -t roads001/httpserver:v2 .	
pushv2:
	docker push roads001/httpserver:v2
run:
	docker container run --rm -d -p8765:8765 roads001/httpserver:v2
