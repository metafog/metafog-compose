<div align="center">

  <h1>Metafog Compose</h1>

  <p>
  Task manager with parallel processing functionality based on a YAML task definition file designed for <a href="https://metafog.io">Metafog's decentralised network<a>. See <a href="examples/">examples</a> to learn more.
  </p>

</div>

Credits: [https://github.com/go-task/task](https://github.com/go-task/task)



This is a trimmed and adapted version of the above library in order to support parallel task processing on iteratables like numeric range, files and records in a file.

```loop``` command is added to support the above feature.

## Binaries

Mac, Windows and Linux binaries are in [bin](bin/) folder. They are statically built without any dependancies.

> Note: Binaries are bundled with [Metafog Gateway](https://metafog.io/) installers.

## Build from source

Clone the repository.

```
$ git clone https://github.com/metafog/metafog-compose
$ cd metafog-compose
$ go mod tidy
$ go run cmd/task/task.go 
```

## Example (sequential)

```
tasks:
  default:
    cmds: 
      - loop:
        range: [10, 20] 
        run: mytask
        parallel: 0
  
  mytask:
    cmds:
      - echo "My Task - Index:{{.INDX}}, Value:{{.ARG}}"
```

```INDX``` is the iteration number starting 1. 

```ARG``` is the number within the range for the iteration. 

## Example (parallel)

```
tasks:
  default:
    cmds: 
      - loop:
        range: [1, 100] 
        run: mytask
        parallel: 5
  
  mytask:
    cmds:
      - echo "My Task - Index:{{.INDX}}, Value:{{.ARG}}"
```

Here, mytask will be executed in parallel with concurrency of 5. Loop command will exit only after all tasks are executed.

## Loop Options 

  - [Range](#range)
  - [List](#list)
  - [Folder](#folder)
  - [File](#file)
  - [Timer](#timer)
  - [FolderWatch](#folderwatch)
  - [ActiveMQ](#activemq)
  - [Kafka](#kafka)


### Range

Iterate through numbers <start>-<end> and run ```task```. 

```ARG``` will be the number in the range.

```
- loop:
  range: [1, 100] 
  run: task1
  parallel: 0
```

Parameters: ```range: [<start>, <end>]``` (Both are inclusive)

### List

Iterate through a list of values and run ```task```. 

```ARG``` will be the value in the list.

```
- loop:
  list: ["Orange", "Apple"] 
  run: task1
  parallel: 0
```

Parameters: ```list: [<value>, ...]```

### Folder

Iterate through all files in the folder (ignore sub folders) and run ```task1```.

```ARG``` will be the file name.

```
- loop:
  folder: /tmp/
  run: task1
  parallel: 0
```

Parameters: ```folder: <folder-path>```


### File

Iterate through each line in the file and run ```task1```. 

```ARG``` will be the contents of the line.

```
- loop:
  file: /tmp/urls.txt
  run: task1
  parallel: 2
```

Parameters: ```file: <file-path>```


### Timer

Run a timer with ```seconds``` interval and run ```task1```. 

```ARG``` will be the time of that execution.

```
- loop:
  timer: 3
  run: task1
  parallel: 2
```

Parameters: ```timer: <interval-in-seconds>```

> This will never end. You have to exit the process by Ctrl+C.

### FolderWatch

Monitor a folder for new files and run ```task1```.

```ARG``` will be the new file path.

```
- loop:
  folder_watch: /tmp/
  run: task1
  parallel: 2
```

Parameters: ```folder_watch: <folder-path>```

> This will never end. You have to exit the process by Ctrl+C.

### ActiveMQ

Subscribe to ActiveMQ topic and run ```task1```.

```ARG``` will be the message body.

```
- loop:
  activemq: ["localhost:61613", "metafog"]
  run: task1
  parallel: 2
```

Parameters: ```activemq: [<connection-url>, <topic-name>]```

> This will never end. You have to exit the process by Ctrl+C.

### Kafka

Subscribe to Kafka messages and run ```task1```.

```ARG``` will be the message body.

```
- loop:
  kafka: ["localhost:9092", "metafog-topic"]
  run: task1
  parallel: 2
```

Parameters: ```kafka: [<connection-url>, <topic-name>]```

> This will never end. You have to exit the process by Ctrl+C.

## Examples

There are many [examples](./examples/) for you to get started with metafog-compose.
