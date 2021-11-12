# do not generate the cert every time as the result is nondeterministic
all: key.pem

key.pem:
	openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 3650 -out cert.pem -subj "/C=US/ST=New York/L=New York/O=/CN=t9" -addext "subjectAltName = DNS:t9"

clean:
	rm -rf *.pem frontend