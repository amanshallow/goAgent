# goAgent
A simple Go agent to facilitate HTTP requests to an API, parse the JSON response body and display required contents in a formatted manner.
Contains the ability to be built as a Docker image and run inside of a Docker container with error and success messages
being submitted to Loggly via a token included in a env.list file.

**Note: Must run Docker image with "sudo docker run --env-file env.list agent" command for the environment variable to set properly
inside of the docker container. Additional flags such as "-d" may also be passed in as needed.**

Changelog:
[9/25]: Agent now includes the ability to run autonomously at the polling rate of 60 seconds. Meaning the agent will fetch information 
every 1 minute from RatesAPI. If the user specifies the polling rate by setting the environment variable "DELAY", it will be given preference and used instead.
