# UBE TODO

All on this list are:
- can be done,
- suggestions,
- proposals of improvements


## Main goal for 2022-Q1:

UBE - Beta candidate, first version by 31-03-2022
 - Complete pipeline ready
   - with Persistence [Storage]
 - 100% test coverage
 - Telemetry [i.e logging, monitoring]
 - Must include an example end to end integration flow
   - Source -> Process -> etc.
   - along with tests, and telemetry.
   - The example must show & explain the implementation, the use of the test suite, and the telemetry.
   - Use of a Localstack implementation, so that others can use this example as the basis to learn how to use the UBE Toolkit.
   - Implement Test_ContainerFramework as an End to End integration test harness for UBE.
     - Provide feedback and direction via the Test_ContainerFramework GitHub repo (PR’s & Issues).
     - Where required, reach out and coordinate with Stefan Golban, and Miguel Fonseca to drive the T_CF development to for the requirements of the UBE Beta Candidate.
   - Use a dedicated AWS Playground ( Ask DevOps )
  
## Suggestion : InnerLambda vs Pipeline actions

I've been thinking about that SQSLambda which the handler needs. I would like it, if it was just an action in the pipeline.
In this case, it is the first action in a pipeline, triggered by something and creating a business-event for the next action.
the “advantage” is when we can have custom pipeline actions, that don't exists in UBE but are part of the customers software.
So if people need an SFTP pipeline action, they can make it themselves, without “polluting” UBE.
Or the other way around: if people made some customer-specific pipeline action, we could easily take it in,, when it is worthy.
Perhaps we have 2 kinds of pipeline actions: receivers (or triggers), and just processors, the last ones getting an event,
do something and push the event out for the next processor.

## Change : Directory pipeline split-up

Since it holds all kinds of pipeline functions or pipeline actions,   
my suggestion is to move them accordingly.
(also change the documentation when done)
So we would have a pipeline directory, holding everything to do with a pipeline,  
and a pipelineaction or pipeline function for those specific ones.
(perhaps a subfolder of pipelines)

Also document more about:
- action
- pipeline
- service
- supervisor

## Change: pipeline mocks

Can be moved to a subdirectory, because it isn't a pipeline function.]

## Suggestion : Enricher (pipeline function)

Have it simplified for Create and Update, like: CreateUpdateEnricher, which can be called like this:  
```
pl.CreateUpdateEnricher(
	createEventName: createEnrich(store, nowFn),
	updateEventName: updateEnrich(store, nowFn),
,
```

## Change : variable "store"

This variable is a weird name, perhaps ubeTable is more appropriate.  
(store reads quite confusing, for example most bretailers will have stores)

## Change : enricher (product)

Add tests

Also this part reads confusing:
```
if err := store.GetEntity(ctx, id, &struct{}{}); err != nil {
	if !errors.Is(err, serrors.NotFound) {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
} else {
	log.Println("duplicate incoming product with id: ", id)
	return p, errors.New("attempt to create a duplicate for product ID: " + id)
}
```
I would like to see it like this (because it is less confusing when the else is reached):
When I look at where it is implemented I doubt if it will ever reach any code after that.
People should never doubt when reading code ;-)
```
err := store.GetEntity(ctx, id, &struct{}{});
if err != nil {
	if !errors.Is(err, serrors.NotFound) {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
} else {
	log.Println("duplicate incoming product with id: ", id)
	return p, errors.New("attempt to create a duplicate for product ID: " + id)
}
```

## Suggestion : Don't make use of AWS directly

If we implement the ideas from gocloud.dev, we have a bridge pattern in place that can handle multiple clouds,  
without us having to support them and still be quite flexible.
Only the necessary should be in the AWS directory.

## Change : Zap-logger to bridge pattern

In this way each team can set up its favorite logger.
(perhaps Zap still being the default)

See also: https://github.com/go-logr/logr

## Change : method-receivers of pipeline functions

If we name all method-receivers like this EventXXX, that would be nice.

## Suggestion : Uploader (pipeline function)

Also the working of this function is a bit vague.
I would like to store an incoming message, but I see it deals with business events.
I also want to use it to store outgoing data, can this be possible?

## Suggestion : Errors

Make them more friendly.

Now we have "failed to get product: %w" but it could be "get product ID '12345' fail: %w".

## Suggestion : Api Request - Respone (Pipeline action)

We can have such API calls to our lambda.
So perhaps a solution would be, that this type of lambda calls the pipeline from the body and in the end returns the result

```
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Params", request.QueryStringParameters)

	// call the pipeline here

	request.Body = "return message"

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}
```

## Suggestion : Secure logger aka sanitise logging information

Since we use a logger, we must be able to use it in a secure way: no sensitive info should be logged.
For example: name, phone number, email address could be filtered out before writing the actual log info.
(I wrote this before the Log4J vulnerability :-) )

## Suggestion : Errors in lambda with stacktrace

Package "github.com/pkg/errors" provides a simple way to wrap errors with extra information.

Everything below the lambda should make use of this:
```
	err := mypackage.MyFunction()
	return errors.Wrap(err, "MyFunction")
```

and the final error itself can be like this:
```
	err := handler.Execute()
	if err != nil {
		// log error with stacktrace
	}
```
## Change : Config

Since reading the config is pretty standard, the only thing that differs is which variables to read,
we could move this into the UBE framework.

Perhaps call it with a map : cfg[string]string
Example:
cfg["region"] = "AWS_REGION"
and it will return another map: cfg["region"] = the value

or call it with a struct
Example:
type Config struct {
	Region     string `lookup:"AWS_REGION"`
	QueueURL   string `lookup:"SQS_QUEUE_URL"`
	BqQueueURL string `lookup:"SQS_BQ_QUEUE_URL"`,
	Bucket     string `lookup:"S3_BUCKET"`,
	DbTable    string `lookup:"DB_TABLE_NAME"`
}
where we can use the tag to lookup the value.

The disadvantages of both: values can be manipulated after reading the config,
but I can see we only use this in the init part of the lambda

Or... we can stop using the config altogether:
instead of doing this:
	inQueue := sqs.NewQueue(c.queueURL)
we could just do:
	inQueue := sqs.NewQueue("SQS_QUEUE_URL")
and in there, we lookup the given value.

## Suggestion : (simplify) call to NewEventHandler

I got this idea while thinking about config....

In the lambda->main we have an initialisation which is not needed there.
Example:
h := product.NewEventHandler(uploader, inQueue, bqPublisher, store)
We could just call it like this:
	h := product.NewEventHandler()
The product->NewEventHandler should know how to initialise itself.

That means: Less code in the lambda itself and more logic to test.

## Suggestion : Incoming message validation

It's simply adding a tag to a model field, like this:
	BKSKU        string `json:"EAN" validate:"required"`

And you can do more that just see if it is present:
validate:
- a valid email address
- a slice with a size bigger than X
- a numeric value bigger than X
- a numeric value bigger then or equal to X
- a validate date (in YYYY-MM-DD format)
- a numeric value less then or equal to X
- a string length equal to X
- a string length maximum of X
- a string length minimum of X
- a string value that is on of "ValueA ValueB ValueC"

## Change : Work with a message that ia an envelope pointing to an S3 bucket

Sometimes we get a message in, like this:
["com.amazon.javamessaging.MessageS3Pointer",{"s3BucketName":"zale144-uat-sqs-payloads","s3Key":"0687abee-4a25-48bd-a98d-3c4008e29e9b"}]
and it points to an S3 bucket which holds the json file with the actual data.

## Suggestion : Work with a message that has multiple records inside

Sometimes the customer sends a lot of records in 1 message.
Ideally, we would be able to handle this in a streaming kind of way.

Let's say: we get 1 "core product" with holds all product variations, sizes and colors.
Or: the latest stock information changes from a store.

## Think about : separation of standard models vs client models

We did some separation, where a part of the model is an "UBE-Standard" a.k.a. as baseModel and another part is client specific.
Things like a product or customer always have commonalities, which could come from a baseModel.
It is then easy to spin some basic functionality up for the client without the specific part.
(And we could have matching reporting as a standard too, to give the customer an idea how this could work.)
We could even have baseModels for certain industries, without having to try to have 1 "one size fits all" solution.

## Think about : client data vs UBE data

Normally, we would have an incoming model feed, transfer that into an UBE model and then proceed.
What is the advantage in doing that?
(The disadvantage is that there's always 1 transformation in between.)

If we want a specific incoming model, then we need to be able to call a feed to model transformer in the pipeline.
This means the model.go also gets a Feed{} struct.
