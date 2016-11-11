# go-codelab
Simple codelab tutorial for learning Go language.

*Notes: This is "unofficial" tutorial. It was made by personal. If there are something wrong, please inform or fix that via Issues or Pull Requests.*

We'll make very simple IoT application service with Go.
Our example is that make the simple server and simple sensor device structures as clients. So the server and devices can communicate each other in realtime using pipeline.

## Environments

* OS Independent.
* If possible, use the latest version because you are beginner of Go.

## Simple Overview

Assume that there are 3 sensors and we have a realtime server for handling the data from pipeline.

```
[Produce the data with own data structure in realtime]
[We use 'JSON' format here for communicating]
   |
   |
Sensor1 \
   |     \
Sensor2 ---------- Queue as Pipeline ---------- IoT Server 
   |     /                                           |
Sensor3 /                                            |
                                           [Processing with data]
                                           [  Logging the data  ]
                                           [ Visualize the data ]
```

The sensors send produced data to pipeline concurrently using goroutine, so pipeline should queueing the streaming data then the server will fetch the data from pipeline by channel and processing with that.

## What you can learn

This is for beginner of Go language.

If you finish this tutorial, you could get knowledges about 'Go application structure', 'Go struct model', 'Concurrency of Go', 'Handling the synchronization using Channel' and 'A way processing the bytes/encoding in Go'.
