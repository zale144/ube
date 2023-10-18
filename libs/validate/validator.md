# Validator

## Uses

- The validator is used to validate a struct:
```go
func main(){
    type InputStruct struct {
        ID int `json:"id" validate:"required,gte=5"`
    }
    
    inputStr := Person{ID:  5}
    err = validate.Struct(&inputStr)
    ...
}
```
in the above example we are evaluating the struct to the required fields set in the 
struct tag `validate:"required,gte=5"`

the validator validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified
and also allows passing of context.Context for contextual validation information.

It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.

for a list of validation tags you may use on structs please read the validator package (https://pkg.go.dev/github.com/go-playground/validator)

There is a maximum of bad fields / bad records and the validator stops if one of those thresholds is reached.
You can influence this if you want other values than the default ones.
The default values are :
50  = maximum number of field errors
10 = maximum number of record errors
