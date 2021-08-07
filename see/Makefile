build:
	GOOS=linux GOARCH=arm go build -v -o bin/see

send:
	scp bin/see pi@joebunyan.local:/home/pi/see

deploy: build send
