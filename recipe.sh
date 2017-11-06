sudo apt-get update
sudo apt-get -y upgrade
sudo curl -O https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
sudo tar -xvf go1.8.linux-amd64.tar.gz
sudo mv go /usr/local

mkdir go
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

go get github.com/bcollazo/pokenalysis
pokenalysis -command=serve
