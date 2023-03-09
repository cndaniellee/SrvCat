# Port Forwarding Tool (Google Validator)

#### Introduction
Through Google validator verification, IP whitelist management
Prevents unauthorized IP addresses from accessing sensitive ports, and ensures simple verification

#### Installation tutorial

1. Build the Project \
   `set GOOS=linux`  \
   `go build -o run main.go`
2. Required documents: \
   `run  config.yml  web/`

#### Instructions for use

1. Edit configuration file \
   `Note the contents` \
   `Secret suggestion >= 32-bit random string` \
   `(can use your password's md5)`
2. Run \
   `nohup ./run &`
3. Visit \
   ` http://127.0.0.1:8888/ `
4. Follow the instructions

#### Final Steps

1. Set the program to boot
2. Firewall sensitive ports (better to restrict access to a fixed IP segment instead)
3. Enjoy it!