# Chik client

Chik client is an application meant to be used on an headless embedded device connected to some sort of 
actuator and/or input device.

Client is configured using the `client.conf` json config file that can be stored either in the same folder
of the executable or inside `/etc/chik/`. 
Configuration file contains two main parameters:
 - server: domain name or ip address and port of the PC where the server instance is running
 - identity: uuid for the client (this field is randomly generated if empty)
 
A number of parameters may be configurable via the application once the connection with the server is available
When running the application for the first time the configuration gets automatically created and populated with default values.
