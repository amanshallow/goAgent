# goAgent
A simple Go agent to facilitate HTTP requests to an API, parse the JSON response body and display required contents in a formatted manner.
Contains the ability to be built as a Docker image and run inside of a Docker container with error and success messages being submitted to Loggly.

Recommended Docker build and run procedure (clutterless):

	- $ docker build -t agent --rm --quiet .	 					// Builds quietly 
	- $ docker run --env-file env.list -d --rm --name myagent agent				// Runs detached, auto remove when stopped.
				**OR**
	- $ docker run --env DELAY=10 -d --rm --name myagent agent				// 10 second polling rate, auto remove when stopped.
				**OR**
	- $ docker run --env-file env.list --env DELAY=10 -d --rm --name myagent agent		// Loggly token in env list file, delay in CMDline.
	- $ docker logs myagent -f								// Live container output
	
Process for removing container and images:

	- $ docker stop myagent				// Stop container myagent
	- $ docker rm myagent				// Remove container
	- $ docker rmi $(docker images -a -q)		// Remove all stopped images
	- $ docker rm $(docker ps -a -q)		// Remove all stopped containers
	
Changelog:
-------------------------------------------------------------
[9/25/20]: 

	- Agent now includes the ability to run autonomously at the polling rate of 60 seconds. Meaning the agent will fetch information every 1 minute from RatesAPI. 
	- If the user specifies the polling rate by setting the environment variable "DELAY", it will be given preference instead.

[9/26/20]: 

	- Implemented: Loggly success and error messages for worker fuction.
	- Fixed: Error level messages for loggly not working properly.
	- Tested: All functions, loggly messages, errors and such.
	- Updated: Set default polling rate to 300 seconds or 5 minutes (env.list).
	- Updated: README.md
	
[10/14/20]: 

	- Implemented: Loggly success and error messages for DynamoDB.
	- Implemented: DynamoDB functionality to push data recieved from RatesAPI.
	- Tested: DynamoDB success and failure inside AWS console and verfied data was recieved properly.
	- Tested: Loggly messages with DynamoDB and existing functionality.
	- Updated: README.md
