package car

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zale144/ube/libs/testengine"
	"github.com/zale144/ube/model"
)

func TestCar(t *testing.T) {
	model.Now = func() time.Time { return dateNow }
	testengine.EventHandler(t, &UBEModel{})
}

var dateNow = time.Date(2021, 11, 22, 3, 4, 5, 0, time.UTC)

func TestFeed_ConvertToModel(t *testing.T) {
	tests := []struct {
		name    string
		feed    Feed
		want    model.Entity
		wantErr bool
	}{
		{
			name: "success",
			feed: Feed{
				Base: model.Base{
					Product:          "one",
					CreatedAt:        "2021-11-22T03:04:05Z",
					Name:             "Opel",
					Description:      "Not the worst car",
					ShortDescription: "Could be better",
				},
				ProductID:       "1",
				SerialNumber:    "123",
				Brand:           "Opel",
				Model:           "Astra",
				BodyType:        "Hmm",
				Length:          "2",
				Width:           "1",
				Height:          "1",
				BootCapacity:    "4",
				BootCapacityMax: "5",
				ModifiedOn:      "2021-11-22T03:04:05Z",
			},
			want: &UBEModel{
				Base: model.Base{
					BaseKey: model.BaseKey{
						ProductID: 1,
					},
					Product:          "one",
					CreatedAt:        "2021-11-22T03:04:05Z",
					Name:             "Opel",
					Description:      "Not the worst car",
					ShortDescription: "Could be better",
				},
				SerialNumber:    "123",
				Brand:           "Opel",
				Model:           "Astra",
				BodyType:        "Hmm",
				Length:          2,
				Width:           1,
				Height:          1,
				BootCapacity:    4,
				BootCapacityMax: 5,
				ModifiedOn:      &dateNow,
			},
			wantErr: false,
		}, {
			name: "failure",
			feed: Feed{
				Base: model.Base{
					Product:          "one",
					CreatedAt:        "2021-11-22T03:04:05Z",
					Name:             "Opel",
					Description:      "Not the worst car",
					ShortDescription: "Could be better",
				},
				ProductID:       "1a",
				SerialNumber:    "123a",
				Brand:           "Opel",
				Model:           "Astra",
				BodyType:        "Hmm",
				Length:          "2a",
				Width:           "1a",
				Height:          "1a",
				BootCapacity:    "4a",
				BootCapacityMax: "5a",
				ModifiedOn:      "bla bla",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.feed.ConvertToModel()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
