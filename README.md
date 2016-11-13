# go-codelab
Simple codelab tutorial for learning Go language.

*Notes: This is "unofficial" tutorial. It was made by personal. If there are something wrong, please inform or fix that via Issues or Pull Requests.*

We'll make very tiny IoT application service with Go.
Our example is that make the simple server and simple sensor device structures as clients. So the server and devices can communicate each other in realtime.

## Environments

* OS Independent.
* If possible, use the latest version because you are beginner of Go.

## Simple Overview

Assume that there are 3 sensors and we have a realtime server for handling the data from pipeline.

```
|---- Sensor1 ----\            /----- Server1 -----\
                   \          /                     \
|---- Sensor2 ----------------------- Server2 ---------------- Log handler
                   /          \                     /               |
|---- Sensor3 ----/            \----- Server3 -----/                |
         |                               |                          |
[Produce the data]             [Handle the request]         [Logging the data]
``` 

The sensors send produced data to pipeline concurrently using goroutine, so pipeline should queueing the streaming data then the server will fetch the data from pipeline by channel and processing with that.

## What you can learn

This is for beginner of Go language.

If you finish this tutorial, you could get knowledges about 

* Go application structure
* How to work Go application with packages
* Go interface/struct model 
* Concurrency in Go
* Communication way between goroutines using Go channel
