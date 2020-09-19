# goAgent
A simple Go agent to facilitate HTTP requests to an API, parse the JSON response body and display required contents in a formatted manner.
Contains the ability to be built as a Docker image and run inside of a Docker container with messages being error and success messages
being submitted to Loggly via env.list file. 
**Note: Must run Docker image with "sudo docker run --env-file env.list agent" command for the environment variable to set properly
inside of the docker container. Additional flags such as "-d" may also be passed in as needed.**
