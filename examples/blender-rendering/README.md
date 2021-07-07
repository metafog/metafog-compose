# Parallel Rendering of Blender Model using Planetr

In this example, a blender model ```.blend``` file is rendered in parallel using 2 Planetr DCU instances.

It uses the public docker image: ```docker.io/nytimes/blender:2.92-cpu-ubuntu18.04``` with CYCLES rendering engine on CPU.

Arguements to the composer are:

* ```INSTANCE_TYPE``` - Planetr DCU instance type.
* ```BLEND_FILE``` - Blender model file.

Clone this repo. Copy your blend file (say myblend.blend) to this folder.

```shell
$ cd <repo-folder>/examples/blender-rendering/
$ <repo-folder>/bin/planetr-compose BLEND_FILE=mymodel.blend INSTANCE_TYPE=g.4xlarge
```

Composer YAML file is using ```range``` option of the ```loop```. Change the start and end frame as neeeded.

[Yaml file](Taskfile.yml) snippet:

```
vars:
  BLENDER_DOCKER_IMAGE: docker.io/nytimes/blender:2.92-cpu-ubuntu18.04

tasks:
  default:
    cmds: 
      - loop:
        range: [1, 10] 
        run: render-frame
        parallel: 2
 ...  
```
