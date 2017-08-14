# Cocktailbot.recipes

Place this folder into `$YOUR_GO_WORKSPACE/src/github.com/shrwdflrst/cocktailbot`

## Local environment

Starting it up; this will setup the VM and install Go 1.8 and Elasticsearch 5.5

    cd resources
    vagrant up

Stopping

    vagrant suspend

Troubleshooting; if there's issues with the VM, try destroying and rebuilding

    vagrant destroy -f
    vagrant up

Accessing local environment

    vagrant ssh

## Commands

These are accessible after compiling with `go build`.

Import recipes into Elasticsearch:

    ./cocktailbot --recipes resources/data/recipes.json


## References

- https://www.digitalocean.com/community/tutorials/how-to-install-java-with-apt-get-on-ubuntu-16-04
- https://askubuntu.com/questions/190582/installing-java-automatically-with-silent-option
- https://www.digitalocean.com/community/tutorials/how-to-install-and-configure-elasticsearch-on-ubuntu-16-04
- https://github.com/olivere/elastic/wiki
- https://github.com/olivere/elastic/issues/525

## Go docs & utilities

- https://gobyexample.com/reading-files
- https://gobyexample.com/json
- Convert JSON to Go structs: https://mholt.github.io/json-to-go/
- https://github.com/golang/go/wiki/CodeReviewComments

## Testing Elasticsearch

     curl -XGET 'localhost:9200/cocktails/recipe/_search?q=id:861'
     curl -XGET 'localhost:9200/cocktails/recipe/_search?q=title:"The Casino Cocktail"'
