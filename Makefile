bin:
	go build -o dist/mta cmd/*
config:
	./dist/mta configservice
hosting:
	./dist/mta hostingservice
help:
	./dist/mta help
helphosting:
	./dist/mta help hostingservice
helpconfig:
	./dist/mta help configservice