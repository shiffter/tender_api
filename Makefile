run: build
	docker run -p 8080:8080 app_avito

build:
	docker build -t app_avito .

clean:
	docker rmi app_avito -f
