#!/usr/bin/env bash
# The output of all these installation steps is noisy. With this utility
# the progress report is nice and concise.

PATH_SELF=/home/ubuntu/go_workspace/src/github.com/shrwdflrst/cocktailbot

function install {
    echo installing $1
    shift
    apt-get -y install "$@" >/dev/null 2>&1
}

function install_go {
    if [ -f /usr/local/go/bin/go ]
    then
        echo "Go already installed."
        exit 0
    fi

    echo "Download and install Go 1.8"
    # https://medium.com/@patdhlk/how-to-install-go-1-8-on-ubuntu-16-04-710967aa53c9
    sudo curl -sO https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
    sudo tar -xf go1.8.linux-amd64.tar.gz
    sudo mv go /usr/local
}

function install_elasticsearch {
    if [ -f ~/.installed_elasticsearch ]
    then
        echo "Elasticsearch already installed."
        exit 0
    fi

    touch ~/.installed_elasticsearch
    # https://www.digitalocean.com/community/tutorials/how-to-install-java-with-apt-get-on-ubuntu-16-04
    sudo add-apt-repository -y ppa:webupd8team/java
    sudo apt-get update >/dev/null 2>&1
    # https://askubuntu.com/questions/190582/installing-java-automatically-with-silent-option
    echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 select true" | sudo debconf-set-selections
    sudo apt-get -y install oracle-java8-installer

    # https://www.digitalocean.com/community/tutorials/how-to-install-and-configure-elasticsearch-on-ubuntu-16-04
    sudo curl -sO https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.5.1.deb
    sudo dpkg -i elasticsearch-5.5.1.deb

    sudo systemctl enable elasticsearch.service
    sudo cp -f ${PATH_SELF}/resources/elasticsearch.yml /etc/elasticsearch/elasticsearch.yml
    sudo systemctl restart elasticsearch
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
install Vim vim
install_go
install_elasticsearch

# Setup bash_profile
cp ${PATH_SELF}/resources/.bash_profile /home/ubuntu
chown ubuntu:ubuntu /home/ubuntu/.bash_profile
mkdir -p /home/ubuntu/go_workspace/{pkg,bin}
mkdir -p /home/ubuntu/go_workspace/src/github.com/shrwdflrst/cocktailbot
chown -R ubuntu:ubuntu /home/ubuntu/go_workspace

# install packages
sudo -u ubuntu /usr/local/go/bin/go get gopkg.in/olivere/elastic.v5

echo "Finished!"
