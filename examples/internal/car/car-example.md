+++ categories = ["Development", "Infrastructure", "Lambda"]
date = "2021-11-12"
description = "Car example explained"
slug = "Car example"
title = "Car example"
+++

# Example : Car

This is an example where the client can not send the information that is needed for the UBE-Model.
That means there is a Feed-model and an UBE-model.

----------

## Flow part 1 : Receiving messages
```mermaid
graph LR;
    Start((Client)) -- new car<br>JSON --> SQS-1(create queue)
    Start((Client)) -- update car JSON --> SQS-2(update queue)
    SQS-1 -- SQS<br>event --> Lambda-1[Lambda]
	SQS-2 -- SQS<br>event --> Lambda-1[Lambda]
	Lambda-1 --> Flow((Flow<BR>part 2))

	style Lambda-1 fill:#a06500
	style SQS-1 fill:#9A225A
	style SQS-2 fill:#9A225A
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
	Decide -- Create --> S3(Store message<BR>on S3)
	Decide -- Update --> Enricher-2(Enricher<HR>patch original)
	S3 --> Enricher-1(Enricher<HR>check dedupe)
	Enricher-1 --> DDB[(DDB<BR>store data)]
	Enricher-2 --> DDB
	DDB --> Publisher

	style Flow fill:#7532a8
	style S3 fill:#039839
	style DDB fill:#303090
```
## Flow part 2 : Pipeline
```mermaid
graph LR;
    Start((Client)) -- CAR<br>JSON --> SQS-1(queue)
    SQS-1 -- SQS<br>event --> Lambda-1[Lambda]
	Lambda-1 --> InputTransformer(get<br>message<br>from<br>sqs-event)
	subgraph inside Lambda
		InputTransformer --> Decide{Message<br>Source}
		Decide -- Create --> S3(Store message<BR>on S3)
		S3 --> Enricher(Enrich & Save<HR>check dedupe)
		Decide -- Update --> Enricher
		Enricher --> Publisher
	end

	Enricher --> DDB[(DDB<BR>store data)]
	DDB --> Enricher
	Publisher --> Stop((Stop))

	style Lambda-1 fill:#a06500
	style SQS-1 fill:#9A225A
	style S3 fill:#039839
	style DDB fill:#303090
	style Stop fill:#600000
```

Based on this pipeline, the following actions are done:
- the incoming message is derived for the sqs-event.
- if message source is Create:
  - this message is stored on S3.
- the data is enriched.
- the data is stored in the database.
- the business event is published.

----------
