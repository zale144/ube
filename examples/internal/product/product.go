package product

import (
	"fmt"

	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"

	"github.com/zale144/ube/examples/internal/store"
)

// TODO:
// alert action

// Product holds the product data
// key: {ProductID_SKU}
type Product struct {
	model.ApparelBase
	EntityKey
	HarmCode        string       `json:"harm_code,omitempty"`
	HarmDescription string       `json:"harm_description,omitempty"`
	Folder          string       `json:"folder,omitempty"`
	Comment         string       `json:"comment,omitempty"`
	Store           *store.Store `json:"store,omitempty"`
	Categories      []string     `json:"categories,omitempty"`
	Images          []string     `json:"images,omitempty"`
	Prices          []*Price     `json:"prices,omitempty"`
}

type EntityKey struct {
	ProductID string `json:"product_id"`
	SKU       string `json:"sku,omitempty"`
	EAN       string `json:"ean,omitempty"`
}

func (e EntityKey) PK() string {
	return fmt.Sprintf("%s_%s", e.ProductID, e.SKU)
}

func (e EntityKey) SK() string {
	return e.EAN
}

// GetKey returns the custom composed key to the entity
func (p Product) GetKey() model.Key {
	return p.EntityKey
}

// Price is the product price structure
type Price struct {
	model.PriceBase
	Campaigns []*model.Campaign `json:"campaigns,omitempty"`
}

const (
	create = "CreateProduct"
	update = "UpdateProduct"

	persistBatchSize = 40
	publishBatchSize = 10
)

// Pipeline constructs a new product event pipeline with the provided dependencies
func (p *Product) Pipeline(
	uploader actions.IUploader,
	downloader actions.IDownloader,
	prodRepo, storeRepo actions.IRepository,
	publisher actions.IPublisher,
	republisher actions.IRepublisher,
) *pl.Pipeline {
	return pl.NewPipeline(
		p,
		pl.InputTransformer(
			// TODO: pl.Or(RecordsFromFilePointer(), RecordsFromKey())
			actions.RecordsFromFilePointer("filePointer", "product", "Zale144", create, downloader),
			actions.RecordsFromKey("product", "Zale144"),
		),
		pl.Uploader(uploader, actions.Skip(update) /*, TODO: action.FailureMandate(model.StopAndRetry)*/),
		pl.Enricher(actions.EnricherMapping{
			create: actions.EnrichEvent(
				actions.WithDedupe(prodRepo),
				actions.WithSubEntity("store", storeRepo),
			),
			update: actions.EnrichEvent(
				actions.WithPatchOriginal(prodRepo),
			),
		}),
		pl.Persister(
			prodRepo,
			actions.BatchSize(persistBatchSize),
			actions.FailureMandate(model.StopAndRetry),
		), // get batch size and such things from config?
		pl.Publisher(publisher, actions.BatchSize(publishBatchSize) /*pl.Async(),*/, actions.FailureMandate(model.StopAndRetry)),
		pl.AfterEach(actions.Republisher(republisher, 3)),
	)
}
