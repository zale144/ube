package product

import (
	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"

	"github.com/zale144/ube/examples/internal/store"
)

// TODO:
// alert action

// SimpleProduct holds the model data
// key: {ProductID_SKU}
type SimpleProduct struct {
	model.ApparelBase
	EAN             string       `json:"ean,omitempty"`
	HarmCode        string       `json:"harm_code,omitempty"`
	HarmDescription string       `json:"harm_description,omitempty"`
	Folder          string       `json:"folder,omitempty"`
	Comment         string       `json:"comment,omitempty"`
	Store           *store.Store `json:"store,omitempty"`
	Categories      []string     `json:"categories,omitempty"`
	Images          []string     `json:"images,omitempty"`
	Prices          []*Price     `json:"prices,omitempty"`
}

// Price is the model price structure
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

var eventMap = map[string]string{ // TODO: move mapping above because of ARN names?
	"createQueue": create,
	"updateQueue": update,
}

// Pipeline constructs a new model event pipeline with the provided dependencies
func (p *SimpleProduct) Pipeline(
	uploader actions.IUploader,
	prodRepo,
	storeRepo actions.IRepository,
	publisher actions.IPublisher,
	republisher actions.IRepublisher,
) *pl.Pipeline {
	return pl.NewPipeline(
		p,
		pl.InputTransformer(actions.EventFromQueueSource("model", "Zale144", eventMap)),
		pl.Uploader(uploader, actions.Skip(update), actions.FailureMandate(model.LogFailureAndContinue)),
		pl.Enricher(actions.EnricherMapping{
			create: actions.EnrichEvent(
				actions.WithDedupe(prodRepo),
				actions.WithSubEntity("store", storeRepo),
			),
			update: actions.EnrichEvent(
				actions.WithPatchOriginal(prodRepo),
			),
		}),
		pl.Persister(prodRepo, actions.BatchSize(persistBatchSize)), // get batch size and such things from config?
		pl.Publisher(publisher, actions.BatchSize(publishBatchSize) /*pl.Async(), */, actions.FailureMandate(model.StopAndRetry)),
		pl.AfterEach(actions.Republisher(republisher, 3)),
	)
}
