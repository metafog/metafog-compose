# Kafka Planetr-Compose Example

If you want to run decentralized network tasks against a Kafka stream, planetr-compose makes it easy for you.

## Install and run Kafka (locally)

If you already have Kafka running, you can skip this section. We will use a Docker compose to run Kafka (minimal version) for this example.


### Run Kafka

Refer [docker-compose.yml](docker-compose.yml) for more information.

```
$ docker-compose up -d
```

Once Kafka/Zookeeper is started, you can publish test messages to Kafka. 

### Create topic
Let us create a topic first.

```
docker exec -it kafka_kafka_1 kafka-topics.sh --create --bootstrap-server kafka:9092 --topic planetr-topic
```

### Publish messages

```
$ curl -XPOST -u admin:admin -d "body=FooBar123" http://localhost:8161/api/message/planetr?type=queue
```


## Run planetr-compose and consume the Kafka messages

Refer [Taskfile.yml](Taskfile.yml).

```
    ...
    cmds: 
      - loop:
        kafka: ["localhost:9092", "planetr-topic"]
        run: capture
        parallel: 0
    ...
```

Where "localhost:9092" is the Kafka connection URL and "planetr-topic" is the topic name.


```shell
$ planetr-compose 
Loop > Kafka localhost:9092 planetr-topic ...
```

## Take screenshots of website URLs consumed from queue

Open another terminal window and start posting messages using curl command to Kafka. We will post ```https://apache.org/``` as the URL to take the screenshot.

```
$ docker exec -it kafka_kafka_1 kafka-console-producer.sh --bootstrap-server kafka:9092 --topic planetr-topic
> 
```

Type ```https://apache.org/``` in the prompt to post to Kafka.


```planetr-compose``` will consume the message ```https://apache.org/``` and run the task ```render```.

Thats it!

Press Ctr+C to stop the planetr-compose. Stop the Kafka containers by ```docker compose down```.
