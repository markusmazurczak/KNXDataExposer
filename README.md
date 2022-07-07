# KNXDataExposer

A little GOLang GIN Server which exposes configured KNX data using REST

## Installation

Clone the repository
```bash
git clone [URL]
```
modify the configuration files (See [Configuration](#Configuration)) to your needs and run
```bash
go run .
```
from the directory where the files where clone'd into.
**You must have a running Postgresql Server in order to use the application.**

## Configuration

The application ships with 2 configuration files
1. app_config.yaml
2. log_config.yaml

### app_config.yaml

The application consists of *two* parts. One part is a connection to your KNX IP Gateway. This part is responsible for fetching exposed data from the KNX-Bus.
The other part is a server-part which gives you the ability to *query* data using REST (HTTP-GET).

Both parts are configured in *app_config.yaml*:

```yaml
#Server configuration
server:
  bindIP: "127.0.0.1" #<- Tell the REST-Server component on which IP it has to listen
  bindPort: 12345 #<- Tell the REST-Server component the port it has to bind to

knx:
  gatewayIP: "192.168.178.2" #<- Input the IP Address of your KNX-Gateway
  gatewayPort: 3671 #<-Input the Port of your KNX-Gateway
  #Can be "tunnel" or "router"
  gatewayType: tunnel #<- This option depends on what type of KNX-IP Gateway you  have.
                      #Configure "tunnel" if you connect using a tunnel-
                      #Gateway or "router" if you connect using a router-
                      #Gateway

postgres:
  username: postgres #<- Input the username for the database connection
  password:          #<- Input the password
  server: localhost  #<- IP address or hostname of your database server
  dbName: knxData    #<- Name of the database
  port: 5432         #<- Port of the database server
  timeZone: Europe/Berlin #<- Timezone

datapoints:  #<- See Datapoint Configuration
  DPT_9001:
    - 12/0/0
    - 12/0/1
    - 12/0/2
    - 12/0/3
    - 12/0/4
    - 12/0/5
    - 12/0/6
    - 12/0/7
    - 12/0/9
  DPT_14056:
    - 6/1/0
    - 6/1/1
```
#### Datapoint configuration
This app works in the way that it connects to your KNX-IP Gateway. Every possible value that is pushed on your gateway can be seen but instead of saving everything, you have to configure what *Group Addresses* are interesting for you to expose them using the REST-Server.

In order to do that you have to configure the DPT (KNX Datatype) with a list of group addresses you are interested in there values.

Lets take the above configuration. Here we have two DPTs configured. DPT_9001 and DPT_14056 (the format **must** be DPT_[actual dpt-number without any dot(.)]).
DPT 9.001 is defined as a temperature value whereas DPT 14.056 is defined as a power value.
Below every DPT-Section you have to configure the group addresses you are interested in as array elements.

In this example the app will listen for temperature values on the group addresses
- 12/0/0
- 12/0/1
- 12/0/2
- 12/0/3
- ...

and for power values on the group addresses
- 6/1/0
- 6/1/1

#### Environment variables

You can override every single option of the *app_config.yaml* with corresponding ENV variables. The ENV var will always take precedence over the value from the configuration file.

Lets take the configuration parameter *gatewayType* from the KNX block:
```yaml
knx:
  gatewayIP: "192.168.178.118"
  gatewayPort: 3671
  #Can be "tunnel" or "router"
  gatewayType: tunnel            #<-THIS ONE
```
If you set an environment variable with the following form (all in UPPERCASE)

*EXPOSE_[VAR]*

you override the value from the file. The file in this example configures the *gatewayType* as *tunnel*. If you want to override it using an environment variable you simply do:

```bash
export EXPOSE_KNX_GATEWAYTYPE=router
```

### app_config.yaml

Logging configuration of this app is done in the file *log_config.yaml*. The [ZAP](https://github.com/uber-go/zap) framework is used for logging.
```yaml
level: debug #<- Configure the desired logging level here
development: false
disableCaller: false
disableStacktrace: false
encoding: console
outputPaths: #<- Configure output
  - stdout    
  #- app.log
encoderConfig:
  timeKey: time
  timeEncoder: iso8601
  messageKey: message
  levelKey: level
  levelEncoder: lowercase
```

## Fetching data

If you want to get data, simply do a HTTP-GET to */dataset* and put the group address you want the value from as parameter *ga* in the URL:

```bash
curl http://127.0.0.1:12345/dataset?ga=6/1/0 -H "Accept: application/json"
{"ID":6,"CreatedAt":"2022-07-06T12:23:53.764244+02:00","UpdatedAt":"2022-07-06T12:23:53.764244+02:00","DeletedAt":null,"Group_Address":"6/1/0","Value":"1486.00","Unit":"W"}
```