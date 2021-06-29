# ActiveMQ Planetr-Compose Example

Many organizations use ActiveMQ as their message broker for projects like IoT. If you want to run decentralized network tasks for each of the message in the ActiveMQ topic queue, planetr-compose makes it easy for you.

## Install and run ActiveMQ (locally)

If you already have ActiveMQ running, you can skip this section. We will use a Docker image to run ActiveMQ for this example.

```
$ docker run -d --rm -p 61616:61616 -p 61613:61613 -p 8161:8161 rmohr/activemq
13ca273b4ff956048b19a8231d6888418ba34c61ccb31c7162441596a326f1c9
```

Once ActiveMQ is started, you can post test messages using HTTP. We will use ```curl``` to post messages to a topic named ```planetr```.

```
$ curl -XPOST -u admin:admin -d "body=FooBar123" http://localhost:8161/api/message/planetr?type=queue
```

ActiveMQ is running now and you know how to post messages to the queue ```planetr```. You can use any other queue name here as well.

## Run planetr-compose and consume the ActiveMQ messages

Refer [Taskfile.yml](Taskfile.yml).

```
    ...
    cmds: 
      - loop:
        activemq: ["localhost:61613", "planetr"]
        run: capture
        parallel: 0
    ...
```

Where "localhost:61613" is the ActiveMQ connection URL and "planetr" is the topic name.

> Change parallel to non zero for parallel processing.

```shell
$ planetr-compose 
Loop > ActiveMQ localhost:61613 ...
```

## Take screenshots of website URLs consumed from queue

Open another terminal window and start posting messages using curl command to ActiveMQ. We will add ```https://apache.org/``` as the URL to take the screenshot.

```
$ curl -XPOST -u admin:admin \
    -d "body=https://apache.org/" \
    http://localhost:8161/api/message/planetr?type=queue
```

```planetr-compose``` will consume the message ```https://apache.org/``` and run the task ```render```.

Thats it!

Press Ctr+C to stop the planetr-compose. Stop the ActiveMQ container by ```docker stop```.
