# Little and very simple cloud-ish service 
## Motivation
I have a small server for personal use and backups at home. Since it's at home, it's not accessible from outside of the home network. I want to create an access point to the server such that I could upload files for backup from the outside. Obviously it is possible in different, easier and more secure ways, but I wanted to try and build something with the emphasis on security for a long time, and this is the perfect oportunity to do so. 

## How it works
The entire system consists of three components: the client, the reciever and the middleman. The client prepares and uploads the files to the middleman, which can be then downloaded by the reciever. 

First, the client recieves a list of files to upload. The client reads all of the files, compresses them using the built-in zip library, and then encrypts them using AES in the counter mode. The package is then authenticated using HMAC and sent as a HTTP POST request to the middleman. 

The middleman accepts an incoming package if and only if:
1) The package with such contents has not been uploaded previously. 
2) No package is currently in the buffer. 
3) It is authenticated. 

After the middleman has accepted the package, it is stored in memory until it is downloaded by the reciever. The reciever accepts the package from the middleman if and only if the package is authenticated.

## Usage 
In case you would like to use this overengineered system, you will need a dedicated server that will host the middleman server. Simply clone this repository and deploy it such that the `./middleman/cmd/main.go` is the "main" file. 

After the server is running, you will need to configure the client and the reciever (see config.json in both client/cmd and reciever/cmd) and compile both main.go files. Then to upload the files simply type `<client-executable> <path-to-client-config-file> <path-to-directory-to-upload>`. To download the recieved files you will need to type `<reciever-executable> <path-to-reciever-config-file> <path-to-destination-directory>`.

USE AT YOUR OWN RISK. 
