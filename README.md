# mTLS Golang Example


## What is mutual TLS (mTLS)?

Mutual TLS, or mTLS for short, is a method for [mutual authentication](https://en.wikipedia.org/wiki/Mutual_authentication). mTLS using TLS do both side authentication & authorization.

mTLS helps ensure that traffic is secure and trusted in both directions between a client and server. 


## How does mTLS work? 

Normally in TLS, the server has a TLS certificate and a public/private key pair, while the client does not. The typical TLS process works like this:

1. Client connects to server
2. Server presents its TLS certificate
3. Client verifies the server's certificate
4. Client and server exchange information over encrypted TLS connection

In mTLS, however, both the client and server have a certificate, and both sides authenticate using their public/private key pair. Compared to regular TLS, there are additional steps in mTLS to verify both parties (additional steps in **bold**):

1. Client connects to server
2. Server presents its TLS certificate
3. Client verifies the server's certificate
4. **Client presents its TLS certificate**
5. **Server verifies the client's certificate**
6. **Server grants access**
7. Client and server exchange information over encrypted TLS connection

## What does This Example do?

**The example to generate to certs and key**


The basic idea as below:

1. **Generate CA Root**. The first thing we need to do to add mTLS to the connection is to generate a self-signed rootCA file that would be used to sign both the server and client cert. 

    ```bash
    openssl req -newkey rsa:2048 \
        -new -nodes -x509 \
        -days ${DAYS} \
        -out ca.crt \
        -keyout ca.key \
        -subj "/C=US/ST=Earth/L=Mountain View/O=MegaEase/OU=MegaCloud/CN=localhost" 
    ```

2. **Generate the Server Certificate**

    ```bash
    #create a key for server
    openssl genrsa -out server.key 2048
    
    #generate the Certificate Signing Request 
    openssl req -new -key server.key -days ${DAYS} -out server.csr \
        -subj "/C=US/ST=Earth/L=Mountain View/O=MegaEase/OU=MegaCloud/CN=localhost" 
    
    #sign it with Root CA
    openssl x509  -req -in server.csr \
        -extfile <(printf "subjectAltName=DNS:localhost") \ 
        -CA ca.crt -CAkey ca.key  \
        -days ${DAYS} -sha256 -CAcreateserial \
        -out server.crt 
    ```
    > Note:  after golang 1.15, we could have the following errors:
    > 
    > x509: certificate relies on legacy Common Name field, use SANs or temporarily enable     Common Name matching with GODEBUG=x509ignoreCN=0"
    > 
    > https://stackoverflow.com/questions/64814173/    how-do-i-use-sans-with-openssl-instead-of-common-name

3. **Generate the Client certification**

    It's similar to server-side 

    ```bash
    openssl genrsa -out ${CLIENT}.key 2048

    openssl req -new -key ${CLIENT}.key -days ${DAYS} -out ${CLIENT}.csr \
        -subj "/C=US/ST=Earth/L=Mountain View/O=$O/OU=$OU/CN=localhost"

    openssl x509  -req -in ${CLIENT}.csr \
        -extfile <(printf "subjectAltName=DNS:localhost") \ 
        -CA ca.crt -CAkey ca.key -out ${CLIENT}.crt -days ${DAYS} -sha256 -CAcreateserial
    ```


The completed script is [certs/key.sh](certs/key.sh)


**The example of how to write the Golang program for both server and client-side**

1. [server.go](server.go) shows how to listen on HTTPS with requiring client-side verification. 

	Note: service side need the three files -  `CA.cert`, `servier.cert`, `server.key`

	```go
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert, //<-- this is the key
		MinVersion: tls.VersionTLS12,
	}
	```

	You can run like this
	
	```bash
	go run server.go
	```

	When the client successfully sent the request. It would output the header and TLS connection state which includes the client's subjects.

	```log
	(HTTP) Listen on :8080
	(HTTPS) Listen on :8443
	2021/12/31 14:47:13 >>>>>>>>>>>>>>>> Header <<<<<<<<<<<<<<<<
	2021/12/31 14:47:13 User-Agent:curl/7.77.0
	2021/12/31 14:47:13 Accept:*/*
	2021/12/31 14:47:13 >>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<
	2021/12/31 14:47:13 Version: 303
	2021/12/31 14:47:13 HandshakeComplete: true
	2021/12/31 14:47:13 DidResume: false
	2021/12/31 14:47:13 CipherSuite: c02f
	2021/12/31 14:47:13 NegotiatedProtocol: h2
	2021/12/31 14:47:13 NegotiatedProtocolIsMutual: true
	2021/12/31 14:47:13 Certificate chain:
	2021/12/31 14:47:13  0 s:/C=[SO]/ST=[Earth]/L=[Mountain]/O=[Client-B]/OU=[Client-B-OU]/CN=localhost
	2021/12/31 14:47:13    i:/C=[SO]/ST=[Earth]/L=[Mountain]/O=[MegaEase]/OU=[MegaCloud]/CN=localhost
	2021/12/31 14:47:13  1 s:/C=[SO]/ST=[Earth]/L=[Mountain]/O=[MegaEase]/OU=[MegaCloud]/CN=localhost
	2021/12/31 14:47:13    i:/C=[SO]/ST=[Earth]/L=[Mountain]/O=[MegaEase]/OU=[MegaCloud]/CN=localhost
	2021/12/31 14:47:13 >>>>>>>>>>>>>>>>> End <<<<<<<<<<<<<<<<<<
	``` 

2. [client.go](client.go) shows how the client connects to the server by using its certification.

	Note: client side need the three files -  `CA.cert`, `client.cert`, `client.key`

	You can indicate `-c=a` for client a, and `-c=b` for client b. such as:
	```bash
	go run client.go -c=a
	go run client.go -c=b
	```

3. You also can use the `curl` to connect to the server.

	```bash
	curl --trace trace.log -k \
		--cacert ./certs/ca.crt \
		--cert ./certs/client.b.crt \
		--key ./certs/client.b.key \
		https://localhost:8443/hello
	```

	- `--trace trace.log` would record the network details of how the client communicates to the server.
	- `-k` because we use a self-signed certificate, so we need to add this.



