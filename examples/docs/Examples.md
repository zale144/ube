+++ categories = ["Development", "Infrastructure", "Lambda"]
date = "2021-11-12"
description = "Examples explained"
slug = "Examples"
title = "Examples"
+++

# Examples

We have a couple of examples, which should showcase the daily usage of UBE.

----------

## Example : Car

This is an example where the client can not send the information that is needed for the UBE Model.
That means there is a Feed-model and an UBE model.

[Details](../examples/internal/car/car-example.md)

----------

## Example : Product

This is an example where the client can  send the information that is needed for the UBE Model.
There are 3 versions of a message.

- the message has information it is a create-message	
- the message has information it is an update-message	
- the message is an S3 filepointer and will create all the records inside the s3-data.


[Details](../example/internal/product/product-example.md)

----------

## Example : Simple Product

This is an example where the client can exactly send the information that is needed for the UBE Model.
In that information is a store-id, which is used to get the store data
and to enrich the simpleProduct with.

[Details](../example/internal/simpleProduct/simpleProduct-example.md)

----------

## Example : Warehouse Stock

This is an example where the client can not send the information that is needed for the UBE Model.
That means there is a Feed-model and an UBE model.

This is only used to create records.

[Details](../example/internal/warehouseStock/warehouseStock-example.md)

----------

