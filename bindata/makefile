# Requires go-bindata. Install with command:
# go get github.com/jteeuwen/go-bindata/...

# do not generate the cert every time as the result is nondeterministic
all: key.pem bindata

key.pem:
	openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 3650 -out cert.pem -subj "/C=US/ST=New York/L=New York/O=/CN=t9"

bindata:
	go-bindata -pkg bindata -o bindata.go *.txt *.pem frontend/...

clean:
	rm -rf *.pem *.go frontend