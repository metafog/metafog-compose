# Parallel Rendering of Blender Model

In this example, a blender model ```.blend``` file is rendered in parallel using 2 Planetr DCU instances.

It uses the public docker image: ```docker.io/nytimes/blender:2.92-cpu-ubuntu18.04``` with CYCLES rendering engine on CPU.

Arguements to the composer are:

* ```INSTANCE_TYPE``` - Planetr DCU instance type.
* ```BLEND_FILE``` - Blender model file.

```shell
$ planetr-compose BLEND_FILE=/temp/animation.blend INSTANCE_TYPE=g.4xlarge 
```

Composer YAML file is using ```range``` option of the ```loop```. Change the start and end frame as neeeded.

Yaml file snippet:

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