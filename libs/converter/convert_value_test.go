package converter

/* ------------------------------- Imports --------------------------- */

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/* -------------------------- Methods/Functions ---------------------- */

func Test_ConvertValue_Wrong_Bool(t *testing.T) {
	result, err := ConvertBool("myBool", "ok")
	assert.Equal(t, `myBool: value 'ok' is not bool`, err.Error())
	assert.False(t, result)
}

func Test_ConvertValue_Good_Bool_True(t *testing.T) {
	result, err := ConvertBool("myBool", "1")
	require.NoError(t, err)
	assert.True(t, result)

	result, err = ConvertBool("myBool", "t")
	require.NoError(t, err)
	assert.True(t, result)

	result, err = ConvertBool("myBool", "T")
	require.NoError(t, err)
	assert.True(t, result)

	result, err = ConvertBool("myBool", "T")
	require.NoError(t, err)
	assert.True(t, result)

	result, err = ConvertBool("myBool", "true")
	require.NoError(t, err)
	assert.True(t, result)

	result, err = ConvertBool("myBool", "TRUE")
	require.NoError(t, err)
	assert.True(t, result)

	result, err = ConvertBool("myBool", "True")
	require.NoError(t, err)
	assert.True(t, result)
}

func Test_ConvertValue_Good_Bool_False(t *testing.T) {
	result, err := ConvertBool("myBool", "0")
	require.NoError(t, err)
	assert.False(t, result)

	result, err = ConvertBool("myBool", "f")
	require.NoError(t, err)
	assert.False(t, result)

	result, err = ConvertBool("myBool", "F")
	require.NoError(t, err)
	assert.False(t, result)

	result, err = ConvertBool("myBool", "false")
	require.NoError(t, err)
	assert.False(t, result)

	result, err = ConvertBool("myBool", "FALSE")
	require.NoError(t, err)
	assert.False(t, result)

	result, err = ConvertBool("myBool", "False")
	require.NoError(t, err)
	assert.False(t, result)
}

func Test_ConvertInt_Wrong(t *testing.T) {
	result, err := ConvertInt("myInt", "x")
	assert.Equal(t, `myInt: value 'x' is not int`, err.Error())
	assert.Equal(t, 0, result)
}

func Test_ConvertInt_Good(t *testing.T) {
	result, err := ConvertInt("myInt", "12345")
	require.NoError(t, err)
	assert.Equal(t, 12345, result)
}

func Test_ConvertFloat64_Wrong(t *testing.T) {
	result, err := ConvertFloat64("myFloat64", "x")
	assert.Equal(t, `myFloat64: value 'x' is not float64`, err.Error())
	assert.Equal(t, float64(0), result)
}

func Test_ConvertFloat64_Empty(t *testing.T) {
	result, err := ConvertFloat64("myFloat64", "")
	require.NoError(t, err)
	assert.Equal(t, float64(0), result)
}

func Test_ConvertFloat64_Good(t *testing.T) {
	result, err := ConvertFloat64("myFloat64", "123.45")
	require.NoError(t, err)
	assert.Equal(t, float64(123.45), result)
}

func Test_ConvertDateTime_Wrong(t *testing.T) {
	result, err := ConvertDateTime("myDateTime", "x2021-12-10T05:30:01")
	assert.Equal(t, `myDateTime: value 'x2021-12-10T05:30:01' is not date-time: Could not find format for "x2021-12-10T05:30:01"`, err.Error())
	assert.Nil(t, result)
}

func Test_ConvertDateTime_GoodSimpleUS(t *testing.T) {
	result, err := ConvertDateTime("myDateTime", "07/09/1964")
	require.NoError(t, err)
	assert.Equal(t, "1964-07-09 00:00:00", result.Format("2006-01-02 15:04:05"))
}

func Test_ConvertDateTime_GoodSimpleEU(t *testing.T) {
	result, err := ConvertDateTime("myDateTime", "1964-07-09")
	require.NoError(t, err)
	assert.Equal(t, "1964-07-09 00:00:00", result.Format("2006-01-02 15:04:05"))
}

func Test_ConvertDateTime_Good(t *testing.T) {
	result, err := ConvertDateTime("myDateTime", "2021-12-10T05:30:01Z")
	require.NoError(t, err)
	assert.Equal(t, "2021-12-10 05:30:01", result.Format("2006-01-02 15:04:05"))
}
