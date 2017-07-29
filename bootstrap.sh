#!/usr/bin/env bash
# The output of all these installation steps is noisy. With this utility
# the progress report is nice and concise.
function install {
    echo installing $1
    shift
    apt-get -y install "$@" >/dev/null 2>&1
}

echo adding swap file
fallocate -l 2G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
echo '/swapfile none swap defaults 0 0' >> /etc/fstab

sudo apt-get update >/dev/null 2>&1
sudo apt-get -y upgrade

install 'development tools' build-essential
install Git git
install Curl curl

echo "Download and install Go 1.8"
# https://medium.com/@patdhlk/how-to-install-go-1-8-on-ubuntu-16-04-710967aa53c9
sudo curl -sO https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
sudo tar -xf go1.8.linux-amd64.tar.gz
sudo mv go /usr/local
echo "export PATH=\$PATH:/usr/local/go/bin" > /home/ubuntu/.bash_profile
echo "export GOPATH=\$HOME" >> /home/ubuntu/.bash_profile

echo "Finished!"
