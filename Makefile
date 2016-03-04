
gendeps:
		godep save \
		github.com/crunchydata/skybridge2/skybridge 

build:
		cd skybridgeserver && make
		cd skybridgeclient && make

image:
		cd images/skybridge && make  

start:
		./sbin/start-cpm.sh

stop:
		./sbin/stop-cpm.sh

clean:
		rm -rf $(GOBIN)/*server* $(GOBIN)/*command*
		rm -rf $(GOPATH)/pkg/linux_amd64/github.com/crunchydata/skybridge2/*.a


