# Cocktailbot.recipes

Place this folder into `$YOUR_GO_WORKSPACE/src/github.com/cocktailbot/app`

    http://127.0.0.1:8080/recipes?ingredients=vodka%20gin

## Local environment

Starting it up; this will setup the VM and install Go 1.8 and Elasticsearch 5.5

    cd resources
    vagrant up
    vagrant ssh

Inside Vagrant

    cd /home/ubuntu/go_workspace/src/github.com/cocktailbot/app
    go run site/server.go

Then you can access the site at: `http://127.0.0.1:8080/`


Stopping

    vagrant suspend

Troubleshooting; if there's issues with the VM, try destroying and rebuilding

    vagrant destroy -f
    vagrant up

Import runs automatically, but if needed to run manually:

    go run importer/cocktail.go --import resources/data/recipes.json resources/data/categories.json

Running firehose/stream script

    set -o allexport; source .env; set +o allexport
    go run twitter/stream/stream.go


<!-- Debug: Search for recipes with `lemon` and `apple` as ingredients

    go run importer/cocktail.go --search lemon apple -->


## References

- https://www.digitalocean.com/community/tutorials/how-to-install-java-with-apt-get-on-ubuntu-16-04
- https://askubuntu.com/questions/190582/installing-java-automatically-with-silent-option
- https://www.digitalocean.com/community/tutorials/how-to-install-and-configure-elasticsearch-on-ubuntu-16-04
- https://github.com/olivere/elastic/wiki
- https://github.com/olivere/elastic/issues/525
- https://serverfault.com/questions/112795/how-to-run-a-server-on-port-80-as-a-normal-user-on-linux
- https://www.elastic.co/guide/en/elasticsearch/guide/current/_index_time_search_as_you_type.html

## Go docs & utilities

- https://gobyexample.com/reading-files
- https://gobyexample.com/json
- Convert JSON to Go structs: https://mholt.github.io/json-to-go/
- https://github.com/golang/go/wiki/CodeReviewComments

## Testing Elasticsearch

    # check logs
    sudo journalctl --unit elasticsearch
    curl -XGET 'localhost:9200/cocktails/\_mapping/'
    curl -XGET 'localhost:9200/cocktails/recipe/\_search?q=id:861'
    curl -XGET 'localhost:9200/cocktails/recipe/\_search?q=title:"The Casino Cocktail"'

## Ansible

    # ping all hosts
    ansible all -m ping -u root -e 'ansible_python_interpreter=/usr/bin/python3'
    ansible all -a "/bin/echo hello" -u root -e 'ansible_python_interpreter=/usr/bin/python3'

Setting up server:

    cd resources/ansible
    ansible-playbook main.yml -i hosts -u root --limit dev
