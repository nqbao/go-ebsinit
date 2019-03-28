all: clean linux

linux:
	GOOS=linux go build -o ebsinit_linux_amd64

clean:
	rm -f ebsinit* go-ebsinit*
