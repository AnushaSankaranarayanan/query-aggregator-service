# Query Aggregator Service

A golang based microservice that  :
 - accepts GET requests with “sortKey” and “limit” parameters.
 - queries three URLs mentioned below:
   - https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json 
   - https://raw.githubusercontent.com/assignment132/assignment/main/google.json 
   - https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json
 - combines the results from all three URLs
 - sorts them by the sortKey
 - returns the response limited by the "limit" parameter. 

## Structure
The structure of the project is following the architecture proposed by Robert C. Martin - [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

The layers proposed by this architecture are inside of internal folder with the following names and its respectives in The Clean Architecture:

* adapter (Interface Adapter)
* entity (Entities - Enterprise Business Roles)
* framework (Framwework and Drivers)
* service (Use Cases - Application Business Rules)

```
|-- api
|   |-- openapi.json
|-- build
|-- cmd
|   |-- microservice
|-- internal
|   |-- adapter
        |-- webserver
            |-- probes
|   |-- consts
|   |-- entity
|   |-- framework
        |-- httpclient
|   |-- service
|-- kube
|-- tests
|   |-- results
|-- .gitignore
|-- CHANGELOG.md
|-- docker-compose.yml
|-- Dockerfile
|-- go.mod
|-- go.sum
|-- Makefile
|-- README.md

```

## Prerequisites
```
Go version >= 1.18
```

## Install
    $ git clone https://github.com/anushasankaranarayanan/query-aggregator-service.git
    $ cd query-aggregator-service

## Build, Tests and Coverage

#### Clean, test and build

- \$ make  all

#### Build

- \$ make build

#### Run tests

- \$ make  test

#### Run tests with coverage

- \$ make cover

---

## Documentation
The module contains the documentation from code, it means that all comments in packages, funcs can generate documentation in HTML.

**Get the godoc**
```
go get golang.org/x/tools/cmd/godoc
```

**Run the command below**
```
godoc -http=:6060
```

**View docs**

Open godocs [here](http://localhost:6060/pkg/github.com/anushasankaranarayanan/query-aggregator-service/)

## Running the service locally
### Non containerized
Create .env file and set the values for below properties. (Sample values given below)
```
LOG_LEVEL=INFO
SERVER_PORT=9000
NAME=query-aggregator-service
VERSION=1.0.0
HTTP_RETRY_MIN_WAIT=5s
HTTP_RETRY_MAX_WAIT=30s
HTTP_MAX_RETRIES=3
```
Navigate to directory:
```
cd cmd/microservice

```
Run `go run -tags real main.go` . Issue requests to : `http://localhost:9000`

### Docker compose
#### Install [Docker](https://www.docker.com/) on your system.

* [Install instructions](https://docs.docker.com/installation/mac/) for Mac OS X
* [Install instructions](https://docs.docker.com/installation/ubuntulinux/) for Ubuntu Linux
* [Install instructions](https://docs.docker.com/installation/) for other platforms

#### Install [Docker Compose](http://docs.docker.com/compose/) on your system.

* Python/pip: `sudo pip install -U docker-compose`
* Other: ``curl -L https://github.com/docker/compose/releases/download/1.1.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose; chmod +x /usr/local/bin/docker-compose``

#### Run `make build` command on the project root directory.

#### Edit `docker-compose.yml` with appropriate environment variables.

#### Run the below command to bring up the service.
```
 docker-compose up --build
```
#### Once the application is up - issue requests to : `http://localhost:9000`.

#### To stop the containers
 ```
docker-compose stop
```
### Kubernetes cluster

#### Ensure that a kubernetes cluster is up and running.For minikube, please refer [here](https://minikube.sigs.k8s.io/docs/start/)

####  Ensure that `kubectl` has access to the running cluster.

####  Run the below command
```
cd kube
./start.sh
```

The service should be up and running.Get the URL of the service by running:
```
minikube service query-aggregator-service --url
```
and issue requests
#### To stop and remove the deployments from kubernetes , run 
```
/kube/clean-up.sh
```


## Sample responses
```
curl --location --request GET 'http://localhost:9000/api/v1/probes/liveness'

{ "name": "query-aggregator-service" }

curl --location --request GET 'http://localhost:9000/api/v1/query?sortKey=relevancescore&limit=3'

{
    "code": 200,
    "status": "OK",
    "message": "Success",
    "data": [
        {
            "url": "www.wikipedia.com/abc1",
            "views": 11000,
            "relevanceScore": 0.1
        },
        {
            "url": "www.example.com/abc1",
            "views": 1000,
            "relevanceScore": 0.1
        },
        {
            "url": "www.wikipedia.com/abc2",
            "views": 12000,
            "relevanceScore": 0.2
        }
    ],
    "count": 3
}

curl --location --request GET 'http://localhost:9000/api/v1/query?sortKey=views&limit=2'

{
    "code": 200,
    "status": "OK",
    "message": "Success",
    "data": [
        {
            "url": "www.example.com/abc1",
            "views": 1000,
            "relevanceScore": 0.1
        },
        {
            "url": "www.example.com/abc2",
            "views": 2000,
            "relevanceScore": 0.2
        }
    ],
    "count": 2
}
curl --location --request GET 'http://localhost:9000/api/v1/query?sortKey=views&limit=200'

{
    "code": 400,
    "status": "ERROR",
    "message": "QueryHandler non recognizable Parameter: limit. Received: 200 Expected: 2 to 199",
    "data": null,
    "count": 0
}

curl --location --request GET 'http://localhost:9000/api/v1/query?sortKey=data'

{
    "code": 400,
    "status": "ERROR",
    "message": "QueryHandler non recognizable Parameter: sortKey. Received: data Expected: relevancescore or views",
    "data": null,
    "count": 0
}


```
## Known caveats
* Swagger assets are included in the service. Moving that to a common module would be a sensible choice
* Currently , this service does not return an error if either of the URLs doesn't respond. It just logs the errors in such cases and proceeds with the rest. This could be improved with a robust error handling mechanism
* Currently, the sort is on ascending order. This could be driven by a query parameter 
* Build tags(fake and real) are used for unit testing. Alternatively , we could use mocks
* The service is not guarded currently
* Logrus is being used for logging. We could use this as a middleware to prevent initializing in multiple places
* sonar.properties file could be included