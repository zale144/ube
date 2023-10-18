+++ categories = ["Development", "Infrastructure", "Lambda"]
date = "2021-11-12"
description = "Inner workings"
slug = "Questions and Answers"
title = "Questions and Answers"
+++

# Questions and Answers

> The customer sends a Feed-Model, but we want to use an UBE Model

In the warehouseStock example, you can see how that's done.
It has a Feed-Model (which is the incoming model) and an UBE Model.

In this case, the incoming Feed-Model is transformed first into that UBE Model.
The pipeline just works with the UBE Model.

> We need to do a custom transformation

> The customer sends multiple records

See example
