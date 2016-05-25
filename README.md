# practice
A practice with nats.io, mgo package and JSON encoding

Step 1. run gnatsd and mongod

Step 2. run subscribe.go and publish.go to send JSON and write into Mongo (make sure data.json sits in the same directory as publish.go)

or

Step 2. run r_subscribe.go and r_publish.go to read from Mongo and output a JSON file

Error Code:
1 - unknown
2 - Empty file
3 - Nats server error
4 - File permission error
5 - File not exist
