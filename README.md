# LogStack

To run you need a MongoDB running on localhost:

1. Start your mongodb
2. Create a database with name: LogStack
3. Execute: git clone git@github.com:prsolucoes/logstack.git  
4. Execute: make deps  
5. Execute: make build  
6. Execute: make start  
7. Open in your browser: http://localhost:8080  

** The application will try connect on your localhost mongdb

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

1. start  = it will kill current LogStack process and start again
2. stop   = it will kill current LogStack process
3. update = it will update code from git

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
