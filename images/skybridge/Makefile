OSFLAVOR=centos7

all:
		cp $(GOBIN)/etcd  ./bin
		cp $(GOBIN)/etcdctl  ./bin
		cp $(GOBIN)/skydns  ./bin
		cp $(GOBIN)/skybridgeserver  ./bin
		cp $(GOBIN)/skybridgeclient  ./bin
		cp ../../sbin/start-skybridge.sh ./bin
		sudo docker build -t skybridge2 -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f skybridge2:latest crunchydata/skybridge2
