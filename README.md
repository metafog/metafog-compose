<div align="center">
  <a href="https://planetr.io">
    <img src="planetr-compose-white.png"/>
  </a>

  <h1>Planetr Compose</h1>

  <p>
  Task manager with parallel processing functionality based on a YAML task definition file designed for <a href="https://planetr.io">Planetr's decentralised network<a>.
  </p>

  <p>
    See <a href="examples/">examples</a> to learn more.
  </p>
</div>

Credits: [https://github.com/go-task/task](https://github.com/go-task/task)



This is a trimmed and adapted version of the above library in order to support parallel task processing on iteratables like numeric range, files and records in a file.

```loop``` command is added to support the above feature.

## Binaries

Mac, Windows and Linux binaries are in [bin](bin/) folder. They are statically built without any dependancies.

> Note: Binaries are bundled with [Planetr Gateway](https://planetr.io/) installers.

## Build from source

Clone the repository.

```
$ git clone https://github.com/planetrio/planetr-compose
$ cd planetr-compose
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

### Range

Iterate through numbers 1-100 and run ```task1```. 

```ARG``` will be the file number within the range.

```
- loop:
  range: [1, 100] 
  run: task1
  parallel: 2
```

### Folder

Iterate through all files in the folder (ignore sub folders) and run ```task1```.

```ARG``` will be the file name.

```
- loop:
  folder: /tmp/
  run: task1
  parallel: 2
```

### File

Iterate through each line in the file and run ```task1```. 

```ARG``` will be the contents of the line.

```
- loop:
  file: /tmp/urls.txt
  run: task1
  parallel: 2
```

## Examples

There are many [examples](./examples/) for you to get started with planetr-compose.