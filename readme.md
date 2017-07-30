# Cocktailbot.recipes


## Local environment

Starting it up; this will setup the VM and install Go 1.8 and Elasticsearch 5.5

    vagrant up

Stopping

    vagrant suspend

Troubleshooting; if there's issues with the VM, try destroying and rebuilding

    vagrant destroy -f
    vagrant up

Accessing local environment

    vagrant ssh


## References

- https://www.digitalocean.com/community/tutorials/how-to-install-java-with-apt-get-on-ubuntu-16-04
- https://askubuntu.com/questions/190582/installing-java-automatically-with-silent-option
- https://www.digitalocean.com/community/tutorials/how-to-install-and-configure-elasticsearch-on-ubuntu-16-04
