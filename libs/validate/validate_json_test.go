package validate_test

/* ------------------------------- Imports --------------------------- */

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zale144/ube/libs/validate"
)

/* -------------------------- Methods/Functions ---------------------- */

func TestStruct_Wrong_Model(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required"`
	}

	person := Person{}
	err = validate.Struct(person)
	assert.EqualError(t, err, `model 'validate_test.Person' should be a pointer`)
}

func TestStruct_PointerToAPointer(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required"`
	}

	person := &Person{}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'**validate_test.Person' is a pointer to a pointer, need a pointer to a struct`)
}

func TestStruct_Wrong_IntMaxValue(t *testing.T) {
	var err error

	type Person struct {
		ID   uint8 `json:"id,string" validate:"required"`
		Test uint8 `json:"test,string" validate:"required,max=10"`
	}

	person := Person{ID: 1, Test: 12}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.test=12' (uint8) should have a maximal value of '10'`)
}

func TestStruct_Good_IntMaxValue(t *testing.T) {
	var err error

	type Person struct {
		ID   uint8 `json:"id,string" validate:"required"`
		Test uint8 `json:"test,string" validate:"required,max=10"`
	}

	person := Person{ID: 1, Test: 10}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Int_Required(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required"`
	}

	person := Person{}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id' is required`)
}

func TestStruct_Good_Int_Required(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required"`
	}

	person := Person{ID: 15}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Int_LTE(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,lte=5"`
	}

	person := Person{ID: 15}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id=15' (int) should have a maximal value of '5'`)
}

func TestStruct_Good_Int_LTE(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,lte=5"`
	}

	person := Person{ID: 5}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Int_LEN(t *testing.T) {
	var err error

	type Person struct {
		ID string `json:"id" validate:"required,len=7"`
	}

	person := Person{ID: "Z144"}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id=Z144' should have a length of '7'`)
}

func TestStruct_Good_String_Len(t *testing.T) {
	var err error

	type Person struct {
		ID string `json:"id" validate:"required,len=7"`
	}

	person := Person{ID: "Zale144"}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_String_Numeric(t *testing.T) {
	var err error

	type Person struct {
		ID string `json:"id" validate:"required,numeric"`
	}

	person := Person{ID: "abc12345.67"}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id=abc12345.67' should have a numeric value`)
}

func TestStruct_Good_String_Numeric(t *testing.T) {
	var err error

	type Person struct {
		ID string `json:"id" validate:"required,numeric"`
	}

	person := Person{ID: "12345.67"}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Int_MAX(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,max=5"`
	}

	person := Person{ID: 16}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id=16' (int) should have a maximal value of '5'`)
}

func TestStruct_Good_Int_MAX(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,max=5"`
	}

	person := Person{ID: 5}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Int_GT(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,gt=5"`
	}

	person := Person{ID: 5}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id=5' (int) should have a value greater than '5'`)
}

func TestStruct_Good_Int_GT(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,gt=5"`
	}

	person := Person{ID: 6}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Int_GTE(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,gte=5"`
	}

	person := Person{ID: 2}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id=2' (int) should have a minimal value of '5'`)
}

func TestStruct_Good_Int_GTE(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,gte=5"`
	}

	person := Person{ID: 5}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Int_LT(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,lt=5"`
	}

	person := Person{ID: 5}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id=5' (int) should have a value less than '5'`)
}

func TestStruct_Good_Int_LT(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,lt=5"`
	}

	person := Person{ID: 4}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Int_MIN(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,min=5"`
	}

	person := Person{ID: 3}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id=3' (int) should have a minimal value of '5'`)
}

func TestStruct_Good_Int_MIN(t *testing.T) {
	var err error

	type Person struct {
		ID int `json:"id" validate:"required,min=5"`
	}

	person := Person{ID: 5}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Bool_Required(t *testing.T) {
	var err error

	type Person struct {
		HasBrain bool `json:"brain" validate:"required"`
	}

	person := Person{}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.brain' is required`)
}

func TestStruct_Good_Bool_True_Required(t *testing.T) {
	var err error

	type Person struct {
		HasBrain bool `json:"brain" validate:"required"`
	}

	person := Person{HasBrain: true}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Good_Bool_False_Required(t *testing.T) {
	var err error

	type Person struct {
		HasBrain bool `json:"brain" validate:""`
	}

	person := Person{HasBrain: false}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_SubStruct_CheckAddress_StreetRequired(t *testing.T) {
	var err error

	type Address struct {
		Street string `json:"street" validate:"required"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	person := Person{Name: "John Doe"}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.address.street' is required`)
}

func TestStruct_Good_SubStruct_CheckAddress_StreetRequired(t *testing.T) {
	var err error

	type Address struct {
		Street string `json:"street" validate:"required"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	person := Person{Name: "John Doe", Address: Address{Street: "road 12"}}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_SubSubStruct_CheckAddress_CountryNameRequired(t *testing.T) {
	var err error

	type Country struct {
		Name string `json:"name" validate:"required"`
	}

	type Address struct {
		Street  string  `json:"street" validate:"required"`
		Country Country `json:"country"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	person := Person{Name: "John Doe", Address: Address{Street: "road 12"}}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.address.country.name' is required`)
}

func TestStruct_Wrong_SubStruct_CheckAddress_Street_WrongMinLength(t *testing.T) {
	var err error

	type Address struct {
		Street string `json:"street" validate:"required,min=5"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	person := Person{Name: "John Doe", Address: Address{Street: "Road"}}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.address.street=Road' should have a minimal length of '5'`)
}

func TestStruct_Good_SubStruct_CheckAddress_Street_MinLength(t *testing.T) {
	var err error

	type Address struct {
		Street string `json:"street" validate:"required,min=5"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	person := Person{Name: "John Doe", Address: Address{Street: "Road 12"}}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Good_SubSubStruct_CheckAddress_CountryNameRequired(t *testing.T) {
	var err error

	type Country struct {
		Name string `json:"name" validate:"required"`
	}

	type Address struct {
		Street  string  `json:"street" validate:"required"`
		Country Country `json:"country"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	person := Person{Name: "John Doe", Address: Address{Street: "road 12", Country: Country{Name: "UK"}}}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_SubSubStruct_CheckAddress_CountryName_MaxLength(t *testing.T) {
	var err error

	type Country struct {
		Name string `json:"name" validate:"required,max=10"`
	}

	type Address struct {
		Street  string  `json:"street" validate:"required"`
		Country Country `json:"country"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address" validate:"required"`
	}

	person := Person{Name: "John Doe", Address: Address{Street: "road 12", Country: Country{Name: "The Netherlands"}}}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.address.country.name=The Netherlands' should have a maximal length of '10'`)
}

func TestStruct_Good_SubSubStruct_CheckAddress_CountryName_MaxLength(t *testing.T) {
	var err error

	type Country struct {
		Name string `json:"name" validate:"required,max=10"`
	}

	type Address struct {
		Street  string  `json:"street" validate:"required"`
		Country Country `json:"country"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address" validate:"required"`
	}

	person := Person{Name: "John Doe", Address: Address{Street: "road 12", Country: Country{Name: "Holland"}}}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Combinations(t *testing.T) {
	var err error

	type Person struct {
		ID  int `json:"id" validate:"required,gte=5"`
		Age int `json:"age" validate:"required,gte=50"`
	}

	person := Person{}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id' is required,'Person.age' is required`)
}

func TestStruct_Wrong_Combinations_OneFieldPresent(t *testing.T) {
	var err error

	type Person struct {
		ID  int `json:"id" validate:"required,gte=5"`
		Age int `json:"age" validate:"required,gte=50"`
	}

	person := Person{Age: 16}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id' is required,'Person.age=16' (int) should have a minimal value of '50'`)
}

func TestStruct_Good_Combinations(t *testing.T) {
	var err error

	type Person struct {
		ID  int `json:"id" validate:"required,gte=5"`
		Age int `json:"age" validate:"required,gte=50"`
	}

	person := Person{ID: 5, Age: 55}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_EmailAddress(t *testing.T) {
	var err error

	type Person struct {
		Email string `json:"email" validate:"email"`
	}

	person := Person{Email: "myaddress"}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.email=myaddress' (string) should be a valid emailaddress`)
}

func TestStruct_Wrong_EmailAddress_LooksOk(t *testing.T) {
	var err error

	type Person struct {
		Email string `json:"email" validate:"email"`
	}

	person := Person{Email: "info@zale144.com"}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.email=info@zale144.com' (string) should be a valid emailaddress`)
}

func TestStruct_Good_EmailAddress(t *testing.T) {
	var err error

	type Person struct {
		Email string `json:"email" validate:"email"`
	}

	person := Person{Email: "info@zale144.com"}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_UnknownValidation(t *testing.T) {
	var err error

	type Person struct {
		Info string `json:"info" validate:"stuff"`
	}

	person := Person{}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `Undefined validation function 'stuff' on field 'Info'`)
}

func TestStruct_Wrong_UnhandledValidation(t *testing.T) {
	var err error

	type Person struct {
		Info string `json:"info" validate:"ipv6"`
	}

	person := Person{}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.info' has an unhandled validation error 'ipv6'`)
}

// TestJSON_Wrong_Dash checks if the fieldname is the real fieldname, not the JSON one
func TestStruct_Wrong_Dash(t *testing.T) {
	var err error

	type Person struct {
		Name string `json:"-" validate:"required"`
	}

	person := Person{}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.Name' is required`)
}

// TestJSON_Wrong_Dash checks if the fieldname is the real fieldname, not the JSON one
func TestStruct_Good_Dash(t *testing.T) {
	var err error

	type Person struct {
		Name string `json:"-" validate:""`
	}

	person := Person{Name: "John Doe"}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_IsDateValid(t *testing.T) {
	var err error

	type Person struct {
		CreatedOn string `json:"createdon" validate:"isDateValid"`
	}

	person := Person{CreatedOn: "details"}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.createdon=details' should have a date (YYYY-MM-DD) value`)
}

func TestStruct_Good_IsDateValid(t *testing.T) {
	var err error

	type Person struct {
		CreatedOn string `json:"createdon" validate:"isDateValid"`
	}

	person := Person{CreatedOn: "2016-01-01"}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Good(t *testing.T) {
	var err error

	type Person struct {
		Info string `json:"info" validate:""`
	}

	person := Person{Info: "details"}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Slice_Required(t *testing.T) {
	var err error

	type Address struct {
		Street string `json:"street" validate:"required"`
	}

	type Person struct {
		Name      string     `json:"name"`
		Addresses []*Address `json:"addresses" validate:"gt=0,dive"`
	}

	person := Person{Name: "John Doe"}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.addresses=[]' (slice) should have a more elements than '0'`)
}

func TestStruct_Wrong_Slice_Required_Empty(t *testing.T) {
	var err error

	type Address struct {
		Street string `json:"street" validate:"required"`
	}

	type Person struct {
		Name      string     `json:"name"`
		Addresses []*Address `json:"addresses" validate:"gt=0,dive"`
	}

	person := Person{Name: "John Doe", Addresses: []*Address{}}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.addresses=[]' (slice) should have a more elements than '0'`)
}

func TestStruct_Good_Slice_Required_Empty(t *testing.T) {
	var err error

	type Address struct {
		Street string `json:"street" validate:"required"`
	}

	type Person struct {
		Name      string     `json:"name"`
		Addresses []*Address `json:"addresses" validate:"gt=0,dive"`
	}

	person := Person{Name: "John Doe", Addresses: []*Address{{Street: "road"}}}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_Wrong_Slice_Required_1_address(t *testing.T) {
	var err error

	type Address struct {
		Street string `json:"street" validate:"required"`
	}

	type Person struct {
		Name      string     `json:"name"`
		Addresses []*Address `json:"addresses" validate:"lt=2,dive"`
	}

	person := Person{Name: "John Doe", Addresses: []*Address{{Street: "road"}, {Street: "blvd"}}}
	err = validate.Struct(&person)

	assert.Contains(t, err.Error(), `(slice) should have a less elements than '2'`)
}

func TestStruct_Good_Slice_Required_1_address(t *testing.T) {
	var err error

	type Address struct {
		Street string `json:"street" validate:"required"`
	}

	type Person struct {
		Name      string     `json:"name"`
		Addresses []*Address `json:"addresses" validate:"lt=2,dive"`
	}

	person := Person{Name: "John Doe", Addresses: []*Address{{Street: "road"}}}
	err = validate.Struct(&person)
	assert.NoError(t, err)
}

func TestStruct_OneOf(t *testing.T) {
	var err error

	type Car struct {
		Make string `json:"make" validate:"oneof=Audi Opel"`
	}

	car := Car{Make: "Audi"}
	err = validate.Struct(&car)
	assert.NoError(t, err)

	car = Car{Make: "Opel"}
	err = validate.Struct(&car)
	assert.NoError(t, err)

	car = Car{Make: "Ford"}
	err = validate.Struct(&car)
	expectedError := fmt.Errorf("'%s' should be one of these allowed values: '%s'", "Car.make", "Audi Opel")
	assert.EqualError(t, err, expectedError.Error())
}

func TestStruct_Wrong_LowValue(t *testing.T) {
	var err error

	type Person struct {
		ID uint8 `json:"id" validate:"required,gt=10"`
	}

	person := Person{ID: 9}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.id=9' (uint8) should have a value greater than '10'`)
}

func TestStruct_Good_Skip(t *testing.T) {
	var err error

	type Person struct {
		ID   uint8 `json:"id" validate:"required"`
		Test uint8 `json:"-" validate:"required,max=10"`
	}

	person := Person{ID: 11}
	err = validate.Struct(&person)
	assert.EqualError(t, err, `'Person.Test' is required`)
}
