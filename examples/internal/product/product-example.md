+++ categories = ["Development", "Infrastructure", "Lambda"]
date = "2021-11-12"
description = "Product example explained"
slug = "Product example"
title = "Product example"
+++

# Example : Product

This is an example where the client can  send the information that is needed for the UBE-Model.
There are 3 versions of a message.

- the message has information it is a create-message	
- the message has information it is an update-message	
- the message is an S3 filepointer and will create all the records inside the s3-data.

----------

## Flow part 1 : Receiving messages
```mermaid
graph LR;
    Start((Client)) -- "JSON" --> SQS-1("queue")
    SQS-1 -- "SQS-event" --> Lambda-1[Lambda]
	Lambda-1 --> Flow(("Flow<BR>part 2"))

	style Lambda-1 fill:#a06500
	style SQS-1 fill:#9A225A
	style Flow fill:#7532a8
```
The customer sends messages to a sqs queue.

The queue will trigger the lambda, that will receive an sqs-event.

----------

## Flow part 2 : Pipeline
```mermaid
graph LR;
 	Flow((Flow<BR>part 2)) --> InputTransformer(get message<br>from sqs-event)
	InputTransformer --> Decide{Message<br>Source}
	Decide -- Create --> S3-1(Store message<BR>on S3)
	Decide -- Update --> Enricher-2(Enricher<HR>patch original)
	Decide -- S3-pointer --> S3-2(Get data<br>from S3)
	S3-1 --> Enricher-1(Enricher<HR>check dedupe)
	S3-2 --> Enricher-1
	Enricher-1 --> DDB[(DDB<BR>store data)]
	Enricher-2 --> DDB
	DDB --> Publisher

	style Flow fill:#7532a8
	style S3-1 fill:#039839
	style S3-2 fill:#039839
	style DDB fill:#303090
```
Based on this pipeline, the following actions are done:
- the incoming message is derived for the sqs-event.
- if message source is Create:
  - this message is stored on S3.
- if message source is an S3-pointer:
  - data is downloaded from S3.
- the data is enriched.
- the data is stored in the database.
- the business event is published.

----------
