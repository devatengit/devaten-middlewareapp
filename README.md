[![Go Report Card](https://goreportcard.com/badge/github.com/team7mysupermon/devaten_middlewareapp)](https://goreportcard.com/report/github.com/team7mysupermon/mySuperMon_Middleware)

# Devaten_Middleware

This middleware was created to have an easy to set up link between Devaten and Prometheus.

This middleware helps you, the user, moniter your database. Through Devaten and Prometheus all the relevant information about different tasks performed on your database will be monitored and saved, and this information is easily accessable through the Prometheus and Devaten dashboard.

Further down this document, you can find a guide on how to install, run and use this middleware.

You must have a Devaten account to use this middleware. You can create an account on their [website](https://Devaten.com/).

## How to Install

To install the middleware locally, you must have docker and docker compose installed and do the following:

Download the docker compose file from the release.

Open the directory in a terminal where the docker compose file is.

Write the following command:
```docker-compose pull```

This will download the docker images locally.

## How to Configure & Run

1. Configure middleware.env file

Add into your configure file your IP address where docker is runnig.

Example

IPv4 Address.: 192.168.1.9

edit middleware.env-->
APP_HOST="http://192.168.1.9:8081"
After change verify devaten is running on this same address. http://192.168.1.9:8081

LOGIN_USER_NAME="Enter Dahshboard username"
After change verify devaten dashboard Username(Email).Enter correct username that you used in singn up to dashboard.

PASSWORD="Enter Dashboard password"
After change verify devaten dashboard Password.Enter correct password that you used in singn up to dashboard.


RECORDING_MAIL="whenFailure"
You get recording mail with status report when you do the start and stop recording.
1. if you want get recording mail when status failure then keep it as "whenFailure".
2. if you want get recording mail every time then keep it as "enable".
3. if you want disabled recording mail  then keep it as "disabled".

EXPLAIN_JSON="enable"
1. if you want get EXPLAIN_JSON of your recording query when status failure then keep it as "whenFailure".
2. if you want get EXPLAIN_JSON of your recording query every time then keep it as "enable".
3. if you want disabled EXPLAIN_JSON of your recording query then keep it as "disabled".


JIRA="disabled"
1. if you configure jira configuration with devaten dashboard and you want to create ticket in your jira board when status failure then keep it as "whenFailure".
2. if you configure jira configuration with devaten dashboard and you  want to create ticket in your jira board then keep it as "enable".
3. if you configure jira configuration with devaten dashboard and you  dont want to create ticket in your jira board then keep it as "disabled".

REPORT="enable"
1. if you want generate report of your recording when status failure then keep it as "whenFailure".
2. if you want generate report of your recording every time then keep it as "enable".
3. if you want disabled generate report of your recording then keep it as "disabled".

SCRAP_INTERVAL_TIME="10"

2. Run
To start program open a terminal and navigate to the folder containing the docker compose file.
Write following command:
```docker-compose up```

Before proceding login. To login, see : [Login](#login)

## How to Use

When the docker image is running, it is running on the local port **8999**, which is the port you can use to start and stop a Devaten recording.

We expose 2 ports for Prometheus. The first is **9090** which is where the "targets" and "graph" for **Prometheus** are located. Also, the image will open the port **9091** that can be used to access information about the recording through **Prometheus.**

The image will export Grafana on port **3000**.

Once the middleware is up and running, you can do the following API calls, API calls can be made through the address-bar in the browser:



### Start Recording

```
localhost:8999/Start/{Usecase name}/{Application Identifier}
```
http://localhost:8999/Start/getCustomer/861632a7-7fde-46a1-b62b-eae111d00115


**Usecase name** can be anything that you choose. eg , Jmeter test this could your test suite name. 

**Application Identifier** can be found in Devaten, under *Applications* and *Application Management.*

eg. 861632a7-7fde-46a1-b62b-eae111d00115


### Stop Recording

```


```
http://localhost:8999/Stop/getCustomer/861632a7-7fde-46a1-b62b-eae111d00115

**Usecase name** has to be the same as the name used to start the recording.


**Application Identifier** has to be the same as the application identifier used to start the recording.

## **Prometheus**

### **Accessing metrics**

*Please remember to login before hand. See subsection [Login](#login)*
Access prometheus dashboard (in browser) on path:http://mymiddelware.localhost:9090/ 
Access Devaten custom metrics in txt format on path: http://localhost:9091/metrics

## **Grafana**

Access Grafana on path: http://localhost:3000/

*OBS! Beware that first time users of Grafana needs to login with credentials: {uname}: admin, {password}: admin*

### **Steps to connect Prometheus to Grafana**

- Press *Add datasource*.
- Select *Prometheus* as the type.
- Fill out the form, with the following info:

    **HTTP**
    - Name: whatever you wanna call it
    - URL: http://prometheus:9090/
    - Access: Server (Default)
    
    The remaining fields should not be altered
    
    **Auth**
    - Basic auth: on
    
    The remaining fields should be left off
    
    **Basic Auth Details**
    - User: *Username for Devaten*
    - Password: *Password for Devaten*
    
    **Alerting**
    - Scrape interval: 5s
    
    All the remaining fields should be left untouched.

- Press: *Save & test*.
- Access metrics in explore.
- See Grafana tutorials for more.

## Swagger

Once the middleware is up and running, swagger documentation will be up on the following page: [http://localhost:8999/swagger/index.html#/](http://localhost:8999/swagger/index.html#/)

### How to use

When the swagger page is opened the API endpoints can be tested by opening a tab and pressing the “try it out” button. Fill out the required information and press execute.
