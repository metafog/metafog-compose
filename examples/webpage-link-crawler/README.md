# Parallel Crawling of Website Page Links

This example uses Planetr Serverless Functions to crawl a list of web pages and extract the hyper link from each page. This is done parallely using planetr-compose.

Please refer to the [WebCrawler Example](https://planetr.io/docs/tutorial-run-decentralized-functions-to-crawl-website-links.html) to understand how Planetr functions (web assembly) is built. For your convinience, the compiled web assembly file [crawler.wasm](crawler.wasm) is included in this folder.


Clone this repo. 

```shell
$ cd <repo-folder>/examples/webpage-link-crawler/
$ planetr-compose 
```

This will crawl each page by executing the crawler.wasm function on the decentralized Planetr network and the result is saved as JSON files (1.json, 2.json etc...) in the current folder. 

Composer YAML file is using ```file``` option of the ```loop```. Planetr function is created once using the global variable declaration to avail the function identifier. See  ```PLANETRFUNCID``` in the YAML file.

```
vars:
  PLANETRFUNCID: 
    sh: planetr func-create crawler.wasm
...
```

Once the crawl is completed, ```planetr func-rm``` is called to delete the function.

```
tasks:
  default:
    cmds: 
      - loop:
        file: urls.txt 
        run: crawl-link
        parallel: 3
      - planetr func-rm {{.PLANETRFUNCID}}
```

```planetr func-run``` will be called against each URL in the file [urls.txt](urls.txt).

```
  crawl-link:
    cmds:
      - echo "Webpage Link Crawler - {{.ARG}}"
      - planetr func-run {{.PLANETRFUNCID}} -a "{\"url\":\"{{.ARG}}\"}" -p > {{.INDX}}.json
```

