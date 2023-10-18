package car

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/libs/converter"
	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"
)

// Feed holds the incoming car data
type Feed struct {
	model.Base
	ProductID       string `json:"product_id,omitempty" fake:"{number:100,1000}"`
	SerialNumber    string `json:"serial_number" fake:"{uuid}"`
	Brand           string `json:"brand" fake:"{carmaker}"`
	Model           string `json:"model" fake:"{carmodel}"`
	BodyType        string `json:"body_type,omitempty" fake:"{cartype}"`
	Length          string `json:"length,omitempty" fake:"{number:1,10}"`
	Width           string `json:"width,omitempty" fake:"{number:1,10}"`
	Height          string `json:"height,omitempty" fake:"{number:1,10}"`
	BootCapacity    string `json:"boot_capacity,omitempty" fake:"{number:1,10}"`
	BootCapacityMax string `json:"boot_capacity_max,omitempty" fake:"{number:1,10}"`
	ModifiedOn      string `json:"modifiedOn" fake:"{date}"`
}

/*
{
	"SourceSystem": "GK",
	"EAN": "8718513342389",
	"StoreID": "1002",
	"StockDateTime": "2021-10-23T20:33:28.000Z",
	"StockLevelQuantity": "-1"
}
{
	"CreateProduct": {
		"SourceSystem": "GK",
		"EAN": "8718513342389",
		"StoreID": "1002",
		"StockDateTime": "2021-10-23T20:33:28.000Z",
		"StockLevelQuantity": "-1"
	}
}
[
	{
		"SourceSystem": "GK",
		"EAN": "8718513342389",
		"StoreID": "1002",
		"StockDateTime": "2021-10-23T20:33:28.000Z",
		"StockLevelQuantity": "-1"
	}
]
{
	"com.amazon.javamessaging.MessageS3Pointer": {
		"bucket": "./testdata/stock",
		"key": "bulk.json"
	}
}
[
	"com.amazon.javamessaging.MessageS3Pointer",
	{
		"s3BucketName": "./testdata/stock",
		"s3Key": "bulk.json"
	}
]
*/

// ConvertToModel should be implemented if you like to have custom conversion logic instead of the automatic one
func (carFeed *Feed) ConvertToModel() (model.Entity, error) {
	var errMessages []string

	// The automatic converter is called first
	entity := &UBEModel{}
	err := converter.ConvertStruct(carFeed, reflect.ValueOf(entity))
	if err != nil {
		errMessages = append(errMessages, err.Error())
	}

	// the "handmade" converter can be written below
	modifiedOn, err := converter.ConvertDateTime("carFeed.ModifiedOn", carFeed.ModifiedOn)
	if err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if len(errMessages) > 0 {
		return nil, errors.New(strings.Join(errMessages, "; "))
	}

	entity.ModifiedOn = modifiedOn

	return entity, nil
}

// UBEModel holds the UBE car data
type UBEModel struct {
	model.Base
	SerialNumber    string     `json:"serial_number"`
	Brand           string     `json:"brand"`
	Model           string     `json:"model"`
	BodyType        string     `json:"body_type"`
	Length          int        `json:"length"`            // mm
	Width           int        `json:"width"`             // mm
	Height          int        `json:"height"`            // mm
	BootCapacity    int        `json:"boot_capacity"`     // dm3
	BootCapacityMax int        `json:"boot_capacity_max"` // dm3
	ModifiedOn      *time.Time `json:"modifiedOn"`
}

func (entity UBEModel) GetKey() model.Key {
	return entity.BaseKey
}

const (
	create = "CreateCar"
	update = "UpdateCar"
)

var eventMap = map[string]string{ // TODO: use env vars for queue source URIs
	"createQueue": create,
	"updateQueue": update,
}

// Pipeline constructs a new car event pipeline with the provided dependencies
func (entity *UBEModel) Pipeline(
	uploader actions.IUploader,
	carRepo actions.IRepository,
	publisher actions.IPublisher,
) *pl.Pipeline { // TODO: remove
	return pl.NewPipeline(
		entity,
		// what needs to happen before having a usable model
		pl.InputTransformer(
			// TODO: pl.And(FeedToModel(), EventFromQueueSource())
			actions.FeedToModel(&Feed{}, &UBEModel{}),
			actions.EventFromQueueSource("car", "Zale144", eventMap),
		),
		// upload only the create message, not the update message
		pl.Uploader(uploader, actions.Skip(update), actions.FailureMandate(model.LogFailureAndContinue)),
		pl.Enricher(actions.EnricherMapping{
			create: actions.EnrichEvent(
				actions.WithDedupe(carRepo),
			),
			update: actions.EnrichEvent(
				actions.WithPatchOriginal(carRepo),
			),
		}),
		pl.Persister(carRepo),
		pl.Publisher(publisher),
	)
}
