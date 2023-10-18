+++ categories = ["Development", "Infrastructure", "Lambda"]
date = "2021-11-12"
description = "Inner workings"
slug = "inner workings"
title = "Inner workings"
+++

# Inner workings

This is information about the inner workings of UBE.

![Overview-Diagram](Overview-Diagram.png)

## Lambda

A lambda is the starting point and is nothing more than calling a configured corresponding handler.

## Handler

The handler is based on the concept of a pipeline. A pipeline is a chain of actions,  
where each action can work with the result of the previous one.

## Pipeline

UBE has several standard pipeline actions to be used which should cover most of the needs.

## Pipeline action

A pipeline action gets a business event in, processes it and pushes a business event out.
Below you can read about the various pipeline actions that exist.

### Acknowledger (pipeline actions)

This will be configured at the end of a pipeline and gives a queue/stream a signal that the incoming message was handled and can be deleted.

### Enricher (pipeline actions)

This enriches the business event entity with a provided EnricherFn function

In the example, the handler calls the Enricher:
```
pl.Enricher(pl.EnricherMapping{
	createEventName: createEnrich(store, nowFn),
	updateEventName: updateEnrich(store, nowFn),
}, nil),
```

Which translates into:  
- if an object is new, call the createEnrich function
- if an object exists, call the updateEnrich function

### InputConverter (pipeline actions)

This converts the input (incoming feed-model) into another model (UBE model)

### Logger (pipeline actions)

This logs the business event.

### Persister (pipeline actions)

This saves the object into the UBE database.

### Publisher (pipeline actions)

This publishes the business event to a next queue.

### Service (pipeline actions)

This steps out to an external solution to manipulate the business event with.
### Uploader (pipeline actions)

This uploads data to a (cloud) filesystem.
