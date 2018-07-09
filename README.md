System Management - Pharos Anchor
=======================================

This provides functionalities to deploy, update, terminate a container or containers to a certain edge device or a group of edge devices. Also, this provides APIs to create, update, and delete a group of edge devices which container(s) can be deployed at the same time.

![](https://github.sec.samsung.net/RS7-EdgeComputing/system-pharos-anchor-go/tree/master/doc/images/system_arch.png?raw=true)

- Pharos Node
    - This container is running on every edge device to handle service deployment requests
    - Please visit [Pharos Node project](https://github.sec.samsung.net/RS7-EdgeComputing/system-pharos-node-go) to know how to build and run Pharos Node service

- Pharos Web Client
    - A web based GUI-Tool for Pharos
    - Please visit [Pharos Node project](https://github.sec.samsung.net/RS7-EdgeComputing/system-pharos-web-client) to know how to build and run Web Client

## Prerequisites ##
- docker-ce
    - Version: 17.09
    - [How to install](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/)
- go compiler
    - Version: 1.8 or above
    - [How to install](https://golang.org/dl/)

## How to build ##
This provides how to build sources codes to an excutable binary and dockerize it to create a Docker image.

#### 1. Executable binary ####
```shell
$ ./build.sh
```
If source codes are successfully built, you can find an output binary file, **pharos-anchor**, on a root of project folder.

#### 2. Docker Image  ####
Next, you can create it to a Docker image.
```shell
$ docker build -t system-pharos-anchor-go-ubuntu -f Dockerfile .
```
If it succeeds, you can see the built image as follows:
```shell
$ sudo docker images
REPOSITORY                         TAG        IMAGE ID        CREATED           SIZE
system-pharos-anchor-go-ubuntu     latest     fcbbd4c401c2    31 seconds ago    157MB
```
## How to download docker image without building project ##
This provides how to download pre-built Docker image.

#### 1. Download Docker image ####
Please visit [Downloads-ubuntu](https://github.sec.samsung.net/RS7-EdgeComputing/system-pharos-anchor-go/releases/download/alpha-1.1_rel/pharos_anchor_ubuntu_x86_64.tar)

#### 2. Load Docker image from tar file ####
```shell
$ docker load -i pharos_anchor_ubuntu_x86_64.tar
```
If it succeeds, you can see the Docker image as follows:
```shell
$ sudo docker images
REPOSITORY                                                                TAG      IMAGE ID        CREATED        SIZE
docker.sec.samsung.net:5000/edge/system-pharos-anchor-go/ubuntu_x86_64    alpha    899dd9fc0f3b    7 weeks ago    156MB
```

## How to run with Docker image ##
Required options to run Docker image
- port
    - 48099:48099
- volume
    - "host folder"/data/db:/data/db (Note that you should replace "host folder" to a desired folder on your host machine)

You can execute it with a Docker image as follows:
```shell
$ docker run -it -p 48099:48099 -v /pharos-anchor/data/db:/data/db system-pharos-anchor-go-ubuntu
```
If it succeeds, you can see log messages on your screen as follows:
```shell
$ docker run -it -p 48099:48099 -v /pharos-anchor/data/db:/data/db system-pharos-anchor-go-ubuntu
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten] MongoDB starting : pid=6 port=27017 dbpath=/data/db 64-bit host=d0a6b9ae16a5
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten] db version v3.4.4
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten] git version: 888390515874a9debd1b6c5d36559ca86b44babd
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten] OpenSSL version: LibreSSL 2.5.5
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten] allocator: system
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten] modules: none
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten] build environment:
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten]     distarch: x86_64
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten]     target_arch: x86_64
2018-01-17T10:29:52.410+0000 I CONTROL  [initandlisten] options: { storage: { mmapv1: { smallFiles: true } } }
2018-01-17T10:29:52.410+0000 W -        [initandlisten] Detected unclean shutdown - /data/db/mongod.lock is not empty.
2018-01-17T10:29:52.413+0000 I -        [initandlisten] Detected data files in /data/db created by the 'wiredTiger' storage engine, so setting the active storage engine to 'wiredTiger'.
2018-01-17T10:29:52.413+0000 W STORAGE  [initandlisten] Recovering data from the last clean checkpoint.
2018-01-17T10:29:52.413+0000 I STORAGE  [initandlisten] 
2018-01-17T10:29:52.413+0000 I STORAGE  [initandlisten] ** WARNING: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine
2018-01-17T10:29:52.413+0000 I STORAGE  [initandlisten] **          See http://dochub.mongodb.org/core/prodnotes-filesystem
2018-01-17T10:29:52.413+0000 I STORAGE  [initandlisten] wiredtiger_open config: create,cache_size=11515M,session_max=20000,eviction=(threads_min=4,threads_max=4),config_base=false,statistics=(fast),log=(enabled=true,archive=true,path=journal,compressor=snappy),file_manager=(close_idle_time=100000),checkpoint=(wait=60,log_size=2GB),statistics_log=(wait=0),
2018-01-17T10:29:52.875+0000 W STORAGE  [initandlisten] Detected configuration for non-active storage engine mmapv1 when current storage engine is wiredTiger
2018-01-17T10:29:52.875+0000 I CONTROL  [initandlisten] 
2018-01-17T10:29:52.875+0000 I CONTROL  [initandlisten] ** WARNING: Access control is not enabled for the database.
2018-01-17T10:29:52.875+0000 I CONTROL  [initandlisten] **          Read and write access to data and configuration is unrestricted.
2018-01-17T10:29:52.875+0000 I CONTROL  [initandlisten] ** WARNING: You are running this process as the root user, which is not recommended.
2018-01-17T10:29:52.875+0000 I CONTROL  [initandlisten] 
2018-01-17T10:29:52.877+0000 I FTDC     [initandlisten] Initializing full-time diagnostic data capture with directory '/data/db/diagnostic.data'
2018-01-17T10:29:52.878+0000 I NETWORK  [thread1] waiting for connections on port 27017
2018-01-17T10:29:53.023+0000 I FTDC     [ftdc] Unclean full-time diagnostic data capture shutdown detected, found interim file, some metrics may have been lost. OK
```
## API Document ##
Pharos Anchor provides a set of REST APIs for its operations. Descriptions for the APIs are stored in <root>/doc folder.
- **[pharos_anchor_api_for_single_device.yaml](https://github.sec.samsung.net/RS7-EdgeComputing/system-pharos-anchor-go/blob/master/doc/pharos_anchor_api_for_single_device.yaml)**
  - Describes APIs of service deployment for single device
- **[pharos_anchor_api_for_multiple_devices.yaml](https://github.sec.samsung.net/RS7-EdgeComputing/system-pharos-anchor-go/blob/master/doc/pharos_anchor_api_for_multiple_devices.yaml)**
  - Describes APIs of group management with a number of devices and service deployment in a group manner

Note that you can visit [Swagger Editor](https://editor.swagger.io/) to graphically investigate the REST APIs in YAML.

## How to work ##
#### 0. Prerequisites ####
  - 1 PC with Ubuntu 14.04(or above) and Docker
  - Pharos Anchor Docker image
    - Please see the above explaination to know how to build Pharos Anchor Docker image
  - Pharos Node Docker image
    - Please visit [Pharos Node project](https://github.sec.samsung.net/RS7-EdgeComputing/system-pharos-node-go) to know how to build Pharos Node Docker image

#### 1. Run Pharos Anchor and Pharos Node containers ####
Run Pharos Anchor container:
```shell
$ docker run -it -p 48099:48099 -v /pharos-anchor/data/db:/data/db system-pharos-anchor-go-ubuntu
```
Run Pharos Node container:
```shell
$ docker run -it -p 48098:48098 -e ANCHOR_ADDRESS=<Pharos Anchor IP> -e NODE_ADDRESS=<Pharos Node IP> -v /pharos-node/data/db:/data/db -v /var/run/docker.sock:/var/run/docker.sock system-pharos-node-go-ubuntu
```

#### 2. Register Pharos Node to Pharos Anchor ####
***TO BE UPDATED***

If you want to verify if the Node is successfully registered to the Anchor, you can send a request as below:

```shell
$ curl -X GET "http://<Pharos Anchor IP>:48099/api/v1/management/nodes" -H "accept: application/json"

{"nodes":[{"apps":[],"config":{"deviceid":"54919CA5-4101-4AE4-595B-353C51AA983C","devicename":"Edge #1","location":"Human readable location","manufacturer":"Manufacturer Name","modelnumber":"Model number as designated by the manufacturer","os":"Operationg system name and version","pinginterval":"10","platform":"Platform name and version","serialnumber":"Serial number"},"id":" "","ip":"<Pharos Anchor IP>","status":"connected"}]}
$
```
Note that values of the above configuration will be changed, later. Only you care is "id" field which is used for an unique identifier of a device that the Node is running on. With the ID, you can send a request of service deployment to the device. Also, you can inspect the device with the ID, not all registered devices as below:
```shell
$ curl -X GET "http://<Pharos Anchor IP>:48099/api/v1/management/nodes/5a69599a1fee050007865f3a" -H "accept: application/json"

{"apps":[],"config":{"deviceid":"54919CA5-4101-4AE4-595B-353C51AA983C","devicename":"Edge #1","location":"Human readable location","manufacturer":"Manufacturer Name","modelnumber":"Model number as designated by the manufacturer","os":"Operationg system name and version","pinginterval":"10","platform":"Platform name and version","serialnumber":"Serial number"},"id":"5a695f2ad5fd9300089dbd91","ip":"<Pharos Anchor IP>","status":"connected"}
$
```
Note that **5a69599a1fee050007865f3a** will be used as a **node ID** of Pharos Node device.

#### 3. Deploy a new service to Pharos Node device ####
Before deploying a new service to a device, you should write a docker-compose YAML file to describe how to run the service. For example, let's deploy MongoDB service. For the deployment, the below description in YAML will be used:
```shell
version: '2'
services:
  mongodb:
    image: mongo:latest
    container_name: "mongodb"
    volumes:
      - ./data/db:/data/db
    ports:
      - 27017:27017
```
Note that, the above YAML description is absolutly identical to one for docker-compose tool and can define a number of containers. We called a unit of the containers an **application**.

Then, let's send a request to install a new application consisting of a MongoDB service to the Node's device:
```shell
$ curl -X POST "http://<Pharos Anchor IP>:48099/api/v1/management/nodes/5a695f2ad5fd9300089dbd91/apps/deploy" -H "accept: application/json" -H "Content-Type: application/json" -d "version: '2'
services:
  mongodb:
    image: mongo:latest
    container_name: "mongodb"
    volumes:
      - ./data/db:/data/db
    ports:
      - 27017:27017"

{"description":"services:\n  mongodb:\n    container_name: mongodb\n    image: mongo:latest\n    ports:\n    - 27017:27017\n    volumes:\n    - ./data/db:/data/db\nversion: \"2\"\n","id":"1f04ccc14635062ad8a478d08dd94ebdd934efa5","images":[{"name":"mongo"}],"services":[{"name":"mongodb","state":{"ExitCode":"0","Status":"running"}}],"state":"DEPLOY"}
$
```
If the MongoDB service is successfully installed, you can check it as below:
```shell
<Pharos Node Machine> $ docker ps
CONTAINER ID        IMAGE                            COMMAND                  CREATED             STATUS              PORTS                      NAMES
3bd2b0ec8126        mongo:latest                     "docker-entrypoint.s…"   8 seconds ago       Up 7 seconds        0.0.0.0:27017->27017/tcp   mongodb
```
Note that **1f04ccc14635062ad8a478d08dd94ebdd934efa5** will be used as a **App ID** of a service of Pharos Node device.
#### 4. Update/Stop/Delete an existing service in Pharos Node device ####

Note that starting, stopping, updating, and deleting an application will be applied to all corresponding containers defined in the application. For example, deleting an application will stop all containers and erase all Docker images defined in an application.

To update a MongoDB service in Pharos Node device if an update of the service is availble, send a corresponding request as below:
```shell
$ curl -X POST "http://<Pharos Anchor IP>:48099/api/v1/management/nodes/5a695f2ad5fd9300089dbd91/apps/1f04ccc14635062ad8a478d08dd94ebdd934efa5/update" -H "accept: application/json"
```

To stop a running service or start a paused service, send corresponding requests as below, respectively:
```shell
(For stop a running service)
$ curl -X POST "http://<Pharos Anchor IP>:48099/api/v1/management/nodes/5a695f2ad5fd9300089dbd91/apps/1f04ccc14635062ad8a478d08dd94ebdd934efa5/stop" -H "accept: application/json"

(For start a paused service)
$ curl -X POST "http://<Pharos Anchor IP>:48099/api/v1/management/nodes/5a695f2ad5fd9300089dbd91/apps/1f04ccc14635062ad8a478d08dd94ebdd934efa5/start" -H "accept: application/json"
```

To stop a service and delete a corresponding Docker image, send a corresponding request as below:
```shell
$ curl -X DELETE "http://<Pharos Anchor IP>:48099/api/v1/management/nodes/5a695f2ad5fd9300089dbd91/apps/1f04ccc14635062ad8a478d08dd94ebdd934efa5" -H "accept: application/json"
```
