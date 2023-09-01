bin:
	go build -o dist/mta cmd/*
config:
	./dist/mta configservice
hosting:
	./dist/mta hosting