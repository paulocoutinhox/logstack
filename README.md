# LogStack

LogStack is a simple log repository made with Go (Golang).

It use a simple plugin mechanism that let you use some different datasources to store and retrive logs.

Today i have implemented MongoDB and ElasticSearch datasources, feel free to implement more :)

# Configuration

LogStack configuration is a simple INI file called "config.ini".

Example for MongoDB as the datasource:

```
[server]
port = 8080
dsname = mongodb
host = localhost
databasename = logstack
```

Example for ElasticSearch as the datasource:

```
[server]
port = 8080
dsname = elasticsearch
host = localhost
databasename = logstack
```

# Starting

1. Start your datastore server
2. Execute: go get github.com/prsolucoes/logstack  
3. Execute: cd $GOPATH/src/github.com/prsolucoes/logstack  
4. Create config file (config.ini) based on some above example  
5. Execute: make deps  
6. Execute: make install  
7. Execute: logstack -f=[ABSOLUTE-PATH-OF-YOUR-CONFIG-FILE]
8. Open in your browser: http://localhost:8080  

# API

1. List(GET): http://localhost:8080/api/log/list?token=[put-your-token-here]&created_at=[start-date-log-optional]
2. Add(POST): http://localhost:8080/api/log/add   [token, type, message]
3. DeleteAll(GET): http://localhost:8080/api/log/deleteAll   [token]
4. StatsByType(GET): http://localhost:8080/api/log/statsByType   [token]

# Log Entity

1. token = your session token, because you can see only logs from specific session token.
2. type = can be any knew type of level log (error, fatal, info, warning, trace, debug, verbose, echo, warning, success)
3. message = any log message

# Command line interface

You can use some make commands to control LogStack service, like start, stop and update from git repository.

1. make stop   = it will kill current LogStack process
2. make update = it will update code from git and install on GOPATH/bin directory

So if you want start your server for example, you only need call "make start".

# Alternative method to Build and Start project

1. go build
2. ./logstack

# Updates In Real Time

You dont need refresh your browser, everything is updated in real time. 

You can leave the stats charts opened in one browser window for example and see the chart being refreshed in real time.  

# Screenshots

[![Main interface](https://github.com/prsolucoes/logstack/raw/master/screenshots/logstack1.png)](http://github.com/prsolucoes/logstack)

[![Stats](https://github.com/prsolucoes/logstack/raw/master/screenshots/logstack2.png)](http://github.com/prsolucoes/logstack)

# Author WebSite

> http://www.pcoutinho.com

# License

MIT
