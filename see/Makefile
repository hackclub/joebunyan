PI_HOSTNAME=pi@joebunyan.local

build:
	GOOS=linux GOARCH=arm go build -v -o bin/see

send:
	scp bin/see $(PI_HOSTNAME):/home/pi/see
	scp see.service $(PI_HOSTNAME):/home/pi/see.service
	echo "Please move see.service to /lib/systemd/system/ on $(PI_HOSTNAME) (it's already on the pi in hte home dir). Then just run \"sudo systemctl start see\" to boot it up!"

deploy: build send
