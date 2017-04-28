1. Generate keyfiles and certificates:
```bash
        CUR_DIR=`pwd`
    	cd $GOROOT/src/crypto/tls && go build generate_cert.go
    	./generate_cert --host <TARGET_DNS_HOSTNAME> --ca
    	mv *.pem $CUR_DIR
    	cd $CUR_DIR

    	openssl genrsa -out lara.rsa 2048
    	openssl rsa -in lara.rsa -pubout > lara.rsa.pub
```
1. Run `build.sh`
1. Copy all files to target machine and run:
    * `install.sh` as root to install application first time
    * `upgrade.sh` as root to upgrade existing installation (restarts application!)
1. On target machine edit /etc/rc.conf.local:
    * set lara_flags if necessary
    * add lara to pkg_scripts
1. On target machine run dist/migration/*.sql if necessary
1. On target machine start service: `/etc/rc.d/lara start`
