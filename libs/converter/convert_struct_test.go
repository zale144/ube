package converter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvert_convertStruct(t *testing.T) {
	type args struct {
		typ, newVal interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success: all",
			args: args{
				typ: &feedObj{
					Integer:     "1234567",
					Integer8:    "-123",
					Integer16:   "12345",
					Integer32:   "1234567898",
					Integer64:   "1234567898765432101",
					Uinteger:    "12343213431234",
					Uinteger8:   "123",
					Uinteger16:  "12343",
					Uinteger32:  "1234321234",
					Uinteger64:  "12344321123443211234",
					Floating32:  "321.123",
					Floating64:  "321.123321",
					Boolean:     "true",
					Timestamp:   "2022-02-04T15:16:22Z",
					TimestampUS: "02/04/2022 15:16:22",
					Date:        "2022-02-04",
					DateUS:      "02/04/2022",
					Object:      &objectFeed{Integer: "321"},
					ObjectList:  []*objectFeed{{Integer: "222"}, {Integer: "5542"}},
					ObjectMap: map[string]*objectFeedT{
						"key1": {Timestamp: "02/04/2022 15:16:22"},
						"key2": {Timestamp: "2022-02-04T15:16:22Z"},
					},
				},
				newVal: &modelObj{},
			},
			want: &modelObj{
				Integer:     1234567,
				Integer8:    -123,
				Integer16:   12345,
				Integer32:   1234567898,
				Integer64:   1234567898765432101,
				Uinteger:    12343213431234,
				Uinteger8:   123,
				Uinteger16:  12343,
				Uinteger32:  1234321234,
				Uinteger64:  12344321123443211234,
				Floating32:  321.123,
				Floating64:  321.123321,
				Boolean:     true,
				Timestamp:   time.Date(2022, time.February, 4, 15, 16, 22, 0, time.UTC),
				TimestampUS: time.Date(2022, time.February, 4, 15, 16, 22, 0, time.UTC),
				Date:        time.Date(2022, time.February, 4, 0, 0, 0, 0, time.UTC),
				DateUS:      time.Date(2022, time.February, 4, 0, 0, 0, 0, time.UTC),
				Object: &object{
					Integer: 321,
				},
				ObjectList: []*object{{Integer: 222}, {Integer: 5542}},
				ObjectMap: map[string]*objectT{
					"key1": {Timestamp: time.Date(2022, time.February, 4, 15, 16, 22, 0, time.UTC)},
					"key2": {Timestamp: time.Date(2022, time.February, 4, 15, 16, 22, 0, time.UTC)},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		}, {
			name: "failure: invalid int",
			args: args{
				typ: &feedObj{
					Integer: "3lk34#fk",
				},
				newVal: &modelObj{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && strings.Contains(err.Error(), `failed to convert field 'Integer': strconv.ParseInt:`)
			},
		}, {
			name: "failure: unsigned int",
			args: args{
				typ: &feedObj{
					Uinteger: "-12343213431234",
				},
				newVal: &modelObj{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == `failed to convert field 'Uinteger': strconv.ParseUint: parsing "-12343213431234": invalid syntax`
			},
		}, {
			name: "failure: bool",
			args: args{
				typ: &feedObj{
					Boolean: "wrong",
				},
				newVal: &modelObj{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == `failed to convert field 'Boolean': strconv.ParseBool: parsing "wrong": invalid syntax`
			},
		}, {
			name: "failure: float32",
			args: args{
				typ: &feedObj{
					Floating32: "wrong",
				},
				newVal: &modelObj{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == `failed to convert field 'Floating32': strconv.ParseFloat: parsing "wrong": invalid syntax`
			},
		}, {
			name: "failure: invalid timestamp",
			args: args{
				typ: &feedObj{
					TimestampUS: "13/04/2022 15:16:22",
				},
				newVal: &modelObj{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == `failed to convert field 'TimestampUS': parsing time "13/04/2022 15:16:22": month out of range`
			},
		}, {
			name: "failure: invalid slice",
			args: args{
				typ: &feedObj{
					ObjectList: []*objectFeed{{Integer: "wrong"}},
				},
				newVal: &modelObj{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == `failed to convert field 'ObjectList': failed to convert array element: failed to convert struct element: failed to convert field 'Integer': strconv.ParseInt: parsing "wrong": invalid syntax`
			},
		}, {
			name: "failure: invalid map",
			args: args{
				typ: &feedObj{
					ObjectMap: map[string]*objectFeedT{"key1": {Timestamp: "wrong"}},
				},
				newVal: &modelObj{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == `failed to convert field 'ObjectMap': failed to convert array element: failed to convert struct element: failed to convert field 'Timestamp': Could not find format for "wrong"`
			},
		}, {
			name: "failure: non-pointer struct",
			args: args{
				typ:    feedObj{},
				newVal: modelObj{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == `non-pointer structs are not assignable`
			},
		}, {
			name: "failure: non-pointer struct field",
			args: args{
				typ: &struct {
					ObjField objectFeed
				}{
					ObjField: objectFeed{
						Integer: "123",
					},
				},
				newVal: &struct {
					ObjField object
				}{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == `failed to convert field 'ObjField': failed to convert struct element: non-pointer structs are not assignable`
			},
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			err := ConvertStruct(test.args.typ, reflect.ValueOf(test.args.newVal))
			if !test.wantErr(t, err, fmt.Sprintf("convertStruct(%v, %v)", test.args.typ, test.args.newVal)) {
				t.Errorf("unexpected error result: %s", err)
				return
			}
			if err == nil {
				assert.Equal(t, test.want, test.args.newVal)
			}
		})
	}
}

type feedObj struct {
	Integer     string
	Integer8    string
	Integer16   string
	Integer32   string
	Integer64   string
	Uinteger    string
	Uinteger8   string
	Uinteger16  string
	Uinteger32  string
	Uinteger64  string
	Floating32  string
	Floating64  string
	Boolean     string
	Timestamp   string
	TimestampUS string
	Date        string
	DateUS      string
	Object      *objectFeed
	ObjectList  []*objectFeed
	ObjectMap   map[string]*objectFeedT
}

type objectFeed struct {
	Integer string
}

type objectFeedT struct {
	Timestamp string
}

type modelObj struct {
	Integer     int
	Integer8    int8
	Integer16   int16
	Integer32   int32
	Integer64   int64
	Uinteger    uint
	Uinteger8   uint8
	Uinteger16  uint16
	Uinteger32  uint32
	Uinteger64  uint64
	Floating32  float32
	Floating64  float64
	Boolean     bool
	Timestamp   time.Time
	TimestampUS time.Time
	Date        time.Time
	DateUS      time.Time
	Object      *object
	ObjectList  []*object
	ObjectMap   map[string]*objectT
}

type object struct {
	Integer int
}

type objectT struct {
	Timestamp time.Time
}

func TestConvertStruct_Int(t *testing.T) {
	var (
		feedModel = struct {
			Number string
		}{Number: "1"}
		ubeModel = struct {
			Number int
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.NoError(t, err)
	assert.Equal(t, 1, ubeModel.Number)
}

func TestConvertStruct_Uint(t *testing.T) {
	var (
		feedModel = struct {
			Number string
		}{Number: "1"}
		ubeModel = struct {
			Number uint
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.NoError(t, err)
	assert.Equal(t, uint(1), ubeModel.Number)
}

func TestConvertStruct_BoolTrue(t *testing.T) {
	var (
		feedModel = struct {
			HasIt string
		}{HasIt: "1"}
		ubeModel = struct {
			HasIt bool
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.NoError(t, err)
	assert.Equal(t, true, ubeModel.HasIt)
}

func TestConvertStruct_BoolFalse(t *testing.T) {
	var (
		feedModel = struct {
			HasIt string
		}{HasIt: "0"}
		ubeModel = struct {
			HasIt bool
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.NoError(t, err)
	assert.Equal(t, false, ubeModel.HasIt)
}

func TestConvertStruct_Float(t *testing.T) {
	var (
		feedModel = struct {
			Number string
		}{Number: "1"}
		ubeModel = struct {
			Number float64
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.NoError(t, err)
	assert.Equal(t, float64(1), ubeModel.Number)
}

func TestConvertStruct_Time(t *testing.T) {
	var (
		feedModel = struct {
			DateTime string
		}{DateTime: "2021-12-10"}
		ubeModel = struct {
			DateTime time.Time
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.NoError(t, err)
	assert.Equal(t, "2021-12-10 00:00:00 +0000 UTC", ubeModel.DateTime.String())
}

func TestConvertStruct_Slice(t *testing.T) {
	var (
		feedModel = struct {
			Numbers []string
		}{Numbers: []string{"1", "2", "3"}}
		ubeModel = struct {
			Numbers []int
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(ubeModel.Numbers))
	assert.Equal(t, int(1), ubeModel.Numbers[0])
	assert.Equal(t, int(2), ubeModel.Numbers[1])
	assert.Equal(t, int(3), ubeModel.Numbers[2])
}

func TestConvertStruct_Map(t *testing.T) {
	var (
		feedModel = struct {
			Numbers map[string]string
		}{
			Numbers: map[string]string{
				"one":   "1",
				"two":   "2",
				"three": "3",
			},
		}
		ubeModel = struct {
			Numbers map[string]int
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(ubeModel.Numbers))
	assert.Equal(t, int(1), ubeModel.Numbers["one"])
	assert.Equal(t, int(2), ubeModel.Numbers["two"])
	assert.Equal(t, int(3), ubeModel.Numbers["three"])
}

func TestConvertStruct_ErrorComplexToFloat(t *testing.T) {
	var (
		feedModel = struct {
			Complex complex128
		}{
			Complex: 1 + 4i,
		}
		ubeModel = struct {
			Complex float64
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.Error(t, err)
}

func TestConvertStruct_ErrorFloatToComplex(t *testing.T) {
	var (
		feedModel = struct {
			Complex float64
		}{
			Complex: 123.45,
		}
		ubeModel = struct {
			Complex complex128
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.Error(t, err)
}

func TestConvertStruct_ErrorStringToComplex(t *testing.T) {
	var (
		feedModel = struct {
			Complex string
		}{
			Complex: "12345",
		}
		ubeModel = struct {
			Complex complex128
		}{}
	)

	err := ConvertStruct(&feedModel, reflect.ValueOf(&ubeModel))
	assert.Error(t, err)
}
