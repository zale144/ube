+++ categories = ["Development", "Infrastructure", "Lambda"]
date = "2021-11-12"
description = "Simple Product example explained"
slug = "Simple Product example"
title = "Simple Product Example"
+++

# Example : Simple Product

This is an example where the client can exactly send the information that is needed for the UBE-Model.
In that information is a store-id, which is used to get the store data
and to enrich the simpleProduct with.

----------

## Flow part 1 : Receiving messages
```mermaid
graph LR;
    Start((Client)) -- new product<br>JSON --> SQS-1(create queue)
    Start((Client)) -- update product<br>JSON --> SQS-2(update queue)
    SQS-1 -- "SQS-event" --> Lambda-1[Lambda]
	SQS-2 -- "SQS-event" --> Lambda-1[Lambda]
	Lambda-1 --> Flow((Flow<BR>part 2))

	style Lambda-1 fill:#a06500
	style SQS-1 fill:#9A225A
	style SQS-2 fill:#9A225A
	style Flow fill:#7532a8
```
The customer sends messages to a sqs queue.

The queue will trigger the lambda, that will receive an sqs-event

----------

## Flow part 2 : Pipeline
```mermaid
graph LR;
 	Flow((Flow<BR>part 2)) --> InputTransformer(get message<br>from sqs-event)
	InputTransformer --> Decide{Message<br>Source}
	Decide -- Create  --> S3(Store message<BR>on S3)
	Decide -- Update --> Enricher-2(Enricher<HR>patch original<br>and get store info)
	S3 --> Enricher-1(Enricher<HR>check dedupe<br>and get store info)
	Enricher-1 --> DDB[(DDB<BR>store data)]
	Enricher-2 --> DDB
	DDB --> Publisher

	style Flow fill:#7532a8
	style S3 fill:#039839
	style DDB fill:#303090
```
Based on this pipeline, the following actions are done:
- the incoming message is derived for the sqs-event.
- if message source is Create:
  - this message is stored on S3.
- the data is enriched (with the store information).
- the data is stored in the database.
- the business event is published.

----------
