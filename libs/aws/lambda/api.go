package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	"github.com/zale144/ube/model" // TODO: decouple
)

type (
	// APIGatewayLambda is the standard AWS APIGateway lambda signature
	APIGatewayLambda func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	// requestHandler provides a standard request handler signature
	requestHandler func(context.Context, *model.Request) (*model.Response, error)
)

// ToAPILambda wraps the pipeline.EventHandler signature with an AWS APIGateway Lambda handler signature
func ToAPILambda(handle requestHandler) APIGatewayLambda {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		innerRequest := &model.Request{
			Path:       request.Path,
			HTTPMethod: request.HTTPMethod,
			Reference: model.Reference{
				Headers:               request.Headers,
				QueryStringParameters: request.QueryStringParameters,
				PathParameters:        request.PathParameters,
			},
			RequestID: request.RequestContext.RequestID,
			SourceURI: request.Resource, // TODO?
			Body:      request.Body,
		}

		innerResponse, err := handle(ctx, innerRequest)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		response := events.APIGatewayProxyResponse{
			StatusCode: innerResponse.StatusCode,
			Headers:    innerResponse.Headers,
			Body:       innerResponse.Body,
		}

		return response, err
	}
}
