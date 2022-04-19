server:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags -static" -o bin/api-linux-amd64 shield/cmd/api
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -extldflags -static" -o bin/api-linux-arm64 shield/cmd/api
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -extldflags -static" -o bin/api-windows-amd64.exe shield/cmd/api
	GOOS=windows GOARCH=386 go build -ldflags "-s -w -extldflags -static" -o bin/api-windows-386.exe shield/cmd/api
test:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags -static" -o bin/balancer-linux-amd64 shield/cmd/balancer
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -extldflags -static" -o bin/balancer-linux-arm64 shield/cmd/balancer
	GOOS=windows GOARCH=386 go build -ldflags "-s -w -extldflags -static" -o bin/gateway-windows-386.exe shield/cmd/gateway
	GOOS=windows GOARCH=386 go build -ldflags "-s -w -extldflags -static" -o bin/client-windows-386.exe shield/cmd/client
