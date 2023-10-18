package warehousestock

import (
	"fmt"
	"time"

	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"
)

// Feed holds the incoming data model
type Feed struct {
	Company                  string `json:"Company" validate:"required"`
	DClocation               string `json:"DClocation" validate:"required"`
	EAN                      string `json:"EAN" validate:"required"`
	StockCategoryCode        string `json:"StockCategoryCode"`
	StockDate                string `json:"StockDate" validate:"required" fake:"{date}"`
	AvailableQuantity        string `json:"AvailableQuantity" fake:"{number:1,10}"`
	OnPOOrderQuantity        string `json:"OnPOOrderQuantity" fake:"{number:1,10}"`
	InTransitQuantity        string `json:"InTransitQuantity" fake:"{number:1,10}"`
	TransferQuantity         string `json:"TransferQuantity" fake:"{number:1,10}"`
	OnSOQuantity             string `json:"OnSOQuantity" fake:"{number:1,10}"`
	OnDeliveryQuantity       string `json:"OnDeliveryQuantity" fake:"{number:1,10}"`
	PackedQuantity           string `json:"PackedQuantity" fake:"{number:1,10}"`
	BlockedQuantity          string `json:"BlockedQuantity" fake:"{number:1,10}"`
	ReservedQuantity         string `json:"ReservedQuantity" fake:"{number:1,10}"`
	InspectionQuantity       string `json:"InspectionQuantity" fake:"{number:1,10}"`
	StockLevelInd            string `json:"StockLevelInd"`
	PAPQuantity              string `json:"PAPQuantity" fake:"{number:1,10}"`
	PAPInTransitQuantity     string `json:"PAPInTransitQuantity" fake:"{number:1,10}"`
	OpenOnSalesOrderQuantity string `json:"OpenOnSalesOrderQuantity" fake:"{number:1,10}"`
	Material                 string `json:"Material"`
	Season                   string `json:"Season"`
	Brand                    string `json:"Brand"`
	Size                     string `json:"Size"`
	Width                    string `json:"Width"`
	JobDateTimeStamp         string `json:"JobDateTimeStamp" validate:"required" fake:"{date}"`
	BatchJobStepName         string `json:"BatchJobStepName"`
}

// UBEModel holds the UBE data model
type UBEModel struct {
	// ID string
	Key
	StockCategoryCode        string    `json:"StockCategoryCode"`
	AvailableQuantity        int       `json:"AvailableQuantity"`
	OnPOOrderQuantity        int       `json:"OnPOOrderQuantity"`
	InTransitQuantity        int       `json:"InTransitQuantity"`
	TransferQuantity         int       `json:"TransferQuantity"`
	OnSOQuantity             int       `json:"OnSOQuantity"`
	OnDeliveryQuantity       int       `json:"OnDeliveryQuantity"`
	PackedQuantity           int       `json:"PackedQuantity"`
	BlockedQuantity          int       `json:"BlockedQuantity"`
	ReservedQuantity         int       `json:"ReservedQuantity"`
	InspectionQuantity       int       `json:"InspectionQuantity"`
	StockLevelInd            string    `json:"StockLevelInd"`
	PAPQuantity              int       `json:"PAPQuantity"`
	PAPInTransitQuantity     int       `json:"PAPInTransitQuantity"`
	OpenOnSalesOrderQuantity int       `json:"OpenOnSalesOrderQuantity"`
	Material                 string    `json:"Material"`
	Season                   string    `json:"Season"`
	Brand                    string    `json:"Brand"`
	Size                     string    `json:"Size"`
	Width                    string    `json:"Width"`
	JobDateTimeStamp         time.Time `json:"JobDateTimeStamp"`
	BatchJobStepName         string    `json:"BatchJobStepName"`
	DocType                  string
}

type Key struct {
	StockDate  time.Time `json:"StockDate,omitempty"`
	EAN        string    `json:"EAN"`
	Company    string    `json:"Company"`
	DClocation string    `json:"DClocation"`
}

func (k Key) PK() string {
	return fmt.Sprintf("%s-%s-%s-%s", k.StockDate, k.EAN, k.Company, k.DClocation)
}

// GetKey returns the key to the entity
func (entity UBEModel) GetKey() model.Key {
	return entity.Key
}

// Pipeline constructs a new warehousestock event pipeline with the provided dependencies
func (entity *UBEModel) Pipeline(
	uploader actions.IUploader,
	warehouseStockRepo actions.IRepository,
	publisher actions.IPublisher,
	rep actions.IRepublisher,
) *pl.Pipeline {
	return pl.NewPipeline(
		entity,
		pl.InputTransformer(actions.FeedToModel(&Feed{}, &UBEModel{}), actions.CreateEvent("warehouseStock", "Zale144")),
		pl.Uploader(uploader, actions.FailureMandate(model.LogFailureAndContinue)),
		pl.Enricher(actions.WithDedupe(warehouseStockRepo)),
		pl.Persister(warehouseStockRepo),
		pl.Publisher(publisher, actions.FailureMandate(model.StopAndRetry)),
		pl.AfterEach(actions.Republisher(rep, 3)),
	)
}
