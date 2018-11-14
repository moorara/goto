package config

import (
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	unexported         string
	FieldString        string          // `flag:"field.string" env:"FIELD_STRING" file:"FIELD_STRING_FILE"`
	FieldBool          bool            // `flag:"field.bool" env:"FIELD_BOOL" file:"FIELD_BOOL_FILE"`
	FieldFloat32       float32         // `flag:"field.float32" env:"FIELD_FLOAT32" file:"FIELD_FLOAT32_FILE"`
	FieldFloat64       float64         // `flag:"field.float64" env:"FIELD_FLOAT64" file:"FIELD_FLOAT64_FILE"`
	FieldInt           int             // `flag:"field.int" env:"FIELD_INT" file:"FIELD_INT_FILE"`
	FieldInt8          int8            // `flag:"field.int8" env:"FIELD_INT8" file:"FIELD_INT8_FILE"`
	FieldInt16         int16           // `flag:"field.int16" env:"FIELD_INT16" file:"FIELD_INT16_FILE"`
	FieldInt32         int32           // `flag:"field.int32" env:"FIELD_INT32" file:"FIELD_INT32_FILE"`
	FieldInt64         int64           // `flag:"field.int64" env:"FIELD_INT64" file:"FIELD_INT64_FILE"`
	FieldUint          uint            // `flag:"field.uint" env:"FIELD_UINT" file:"FIELD_UINT_FILE"`
	FieldUint8         uint8           // `flag:"field.uint8" env:"FIELD_UINT8" file:"FIELD_UINT8_FILE"`
	FieldUint16        uint16          // `flag:"field.uint16" env:"FIELD_UINT16" file:"FIELD_UINT16_FILE"`
	FieldUint32        uint32          // `flag:"field.uint32" env:"FIELD_UINT32" file:"FIELD_UINT32_FILE"`
	FieldUint64        uint64          // `flag:"field.uint64" env:"FIELD_UINT64" file:"FIELD_UINT64_FILE"`
	FieldDuration      time.Duration   // `flag:"field.duration" env:"FIELD_DURATION" file:"FIELD_DURATION_FILE"`
	FieldURL           url.URL         // `flag:"field.url" env:"FIELD_URL" file:"FIELD_URL_FILE"`
	FieldStringArray   []string        // `flag:"field.string.array" env:"FIELD_STRING_ARRAY" file:"FIELD_STRING_ARRAY_FILE" sep:","`
	FieldFloat32Array  []float32       // `flag:"field.float32.array" env:"FIELD_FLOAT32_ARRAY" file:"FIELD_FLOAT32_ARRAY_FILE" sep:","`
	FieldFloat64Array  []float64       // `flag:"field.float64.array" env:"FIELD_FLOAT64_ARRAY" file:"FIELD_FLOAT64_ARRAY_FILE" sep:","`
	FieldIntArray      []int           // `flag:"field.int.array" env:"FIELD_INT_ARRAY" file:"FIELD_INT_ARRAY_FILE" sep:","`
	FieldInt8Array     []int8          // `flag:"field.int8.array" env:"FIELD_INT8_ARRAY" file:"FIELD_INT8_ARRAY_FILE" sep:","`
	FieldInt16Array    []int16         // `flag:"field.int16.array" env:"FIELD_INT16_ARRAY" file:"FIELD_INT16_ARRAY_FILE" sep:","`
	FieldInt32Array    []int32         // `flag:"field.int32.array" env:"FIELD_INT32_ARRAY" file:"FIELD_INT32_ARRAY_FILE" sep:","`
	FieldInt64Array    []int64         // `flag:"field.int64.array" env:"FIELD_INT64_ARRAY" file:"FIELD_INT64_ARRAY_FILE" sep:","`
	FieldUintArray     []uint          // `flag:"field.uint.array" env:"FIELD_UINT_ARRAY" file:"FIELD_UINT_ARRAY_FILE" sep:","`
	FieldUint8Array    []uint8         // `flag:"field.uint8.array" env:"FIELD_UINT8_ARRAY" file:"FIELD_UINT8_ARRAY_FILE" sep:","`
	FieldUint16Array   []uint16        // `flag:"field.uint16.array" env:"FIELD_UINT16_ARRAY" file:"FIELD_UINT16_ARRAY_FILE" sep:","`
	FieldUint32Array   []uint32        // `flag:"field.uint32.array" env:"FIELD_UINT32_ARRAY" file:"FIELD_UINT32_ARRAY_FILE" sep:","`
	FieldUint64Array   []uint64        // `flag:"field.uint64.array" env:"FIELD_UINT64_ARRAY" file:"FIELD_UINT64_ARRAY_FILE" sep:","`
	FieldDurationArray []time.Duration // `flag:"field.duration.array" env:"FIELD_DURATION_ARRAY" file:"FIELD_DURATION_ARRAY_FILE" sep:","`
	FieldURLArray      []url.URL       // `flag:"field.url.array" env:"FIELD_URL_ARRAY" file:"FIELD_URL_ARRAY_FILE" sep:","`
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name           string
		expectedTokens []string
	}{
		{"c", []string{"c"}},
		{"C", []string{"C"}},
		{"camel", []string{"camel"}},
		{"Camel", []string{"Camel"}},
		{"camelCase", []string{"camel", "Case"}},
		{"CamelCase", []string{"Camel", "Case"}},
		{"OneTwoThree", []string{"One", "Two", "Three"}},
		{"DatabaseURL", []string{"Database", "URL"}},
		{"DBEndpoints", []string{"DB", "Endpoints"}},
	}

	for _, tc := range tests {
		tokens := tokenize(tc.name)
		assert.Equal(t, tc.expectedTokens, tokens)
	}
}

func TestGetFlagName(t *testing.T) {
	tests := []struct {
		fieldName        string
		expectedFlagName string
	}{
		{"c", "c"},
		{"C", "c"},
		{"camel", "camel"},
		{"Camel", "camel"},
		{"camelCase", "camel.case"},
		{"CamelCase", "camel.case"},
		{"OneTwoThree", "one.two.three"},
		{"DatabaseURL", "database.url"},
		{"DBEndpoints", "db.endpoints"},
	}

	for _, tc := range tests {
		flagName := getFlagName(tc.fieldName)
		assert.Equal(t, tc.expectedFlagName, flagName)
	}
}

func TestGetEnvVarName(t *testing.T) {
	tests := []struct {
		fieldName          string
		expectedEnvVarName string
	}{
		{"c", "C"},
		{"C", "C"},
		{"camel", "CAMEL"},
		{"Camel", "CAMEL"},
		{"camelCase", "CAMEL_CASE"},
		{"CamelCase", "CAMEL_CASE"},
		{"OneTwoThree", "ONE_TWO_THREE"},
		{"DatabaseURL", "DATABASE_URL"},
		{"DBEndpoints", "DB_ENDPOINTS"},
	}

	for _, tc := range tests {
		envVarName := getEnvVarName(tc.fieldName)
		assert.Equal(t, tc.expectedEnvVarName, envVarName)
	}
}

func TestGetFlagValue(t *testing.T) {
	tests := []struct {
		args              []string
		flagName          string
		expectedFlagValue string
	}{
		{[]string{"exe", "-enabled"}, "enabled", "true"},
		{[]string{"exe", "--enabled"}, "enabled", "true"},
		{[]string{"exe", "-enabled=false"}, "enabled", "false"},
		{[]string{"exe", "--enabled=false"}, "enabled", "false"},
		{[]string{"exe", "-enabled", "false"}, "enabled", "false"},
		{[]string{"exe", "--enabled", "false"}, "enabled", "false"},

		{[]string{"exe", "-port=-10"}, "port", "-10"},
		{[]string{"exe", "--port=-10"}, "port", "-10"},
		{[]string{"exe", "-port", "-10"}, "port", "-10"},
		{[]string{"exe", "--port", "-10"}, "port", "-10"},

		{[]string{"exe", "-text=content"}, "text", "content"},
		{[]string{"exe", "--text=content"}, "text", "content"},
		{[]string{"exe", "-text", "content"}, "text", "content"},
		{[]string{"exe", "--text", "content"}, "text", "content"},

		{[]string{"exe", "-enabled", "-text", "content"}, "enabled", "true"},
		{[]string{"exe", "--enabled", "--text", "content"}, "enabled", "true"},

		{[]string{"exec", "-service.name=go-service"}, "service.name", "go-service"},
		{[]string{"exec", "--service.name=go-service"}, "service.name", "go-service"},
		{[]string{"exec", "-service.name", "go-service"}, "service.name", "go-service"},
		{[]string{"exec", "--service.name", "go-service"}, "service.name", "go-service"},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		os.Args = tc.args
		flagValue := getFlagValue(tc.flagName)

		assert.Equal(t, tc.expectedFlagValue, flagValue)
	}
}

func TestGetFieldValue(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		env, envValue     string
		file, fileContent string
		flag              string
		expectedValue     string
	}{
		{
			"FromFlag#01",
			[]string{"/path/to/executable", "-log.level=debug"},
			"LOG_LEVEL", "info",
			"LOG_LEVEL_FILE", "error",
			"log.level", "debug",
		},
		{
			"FromFlag#02",
			[]string{"/path/to/executable", "--log.level=debug"},
			"LOG_LEVEL", "info",
			"LOG_LEVEL_FILE", "error",
			"log.level", "debug",
		},
		{
			"FromFlag#03",
			[]string{"/path/to/executable", "-log.level", "debug"},
			"LOG_LEVEL", "info",
			"LOG_LEVEL_FILE", "error",
			"log.level", "debug",
		},
		{
			"FromFlag#04",
			[]string{"/path/to/executable", "--log.level", "debug"},
			"LOG_LEVEL", "info",
			"LOG_LEVEL_FILE", "error",
			"log.level", "debug",
		},
		{
			"FromEnvironmentVariable",
			[]string{"/path/to/executable"},
			"LOG_LEVEL", "info",
			"LOG_LEVEL_FILE", "error",
			"log.level", "info",
		},
		{
			"FromFileContent",
			[]string{"/path/to/executable"},
			"LOG_LEVEL", "",
			"LOG_LEVEL_FILE", "error",
			"log.level", "error",
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set value using a flag
			os.Args = tc.args

			// Set value in an environment variable
			err := os.Setenv(tc.env, tc.envValue)
			assert.NoError(t, err)

			// Write value in a temporary file
			tmpfile, err := ioutil.TempFile("", "gotest_")
			assert.NoError(t, err)
			defer os.Remove(tmpfile.Name())
			_, err = tmpfile.WriteString(tc.fileContent)
			assert.NoError(t, err)
			err = tmpfile.Close()
			assert.NoError(t, err)
			err = os.Setenv(tc.file, tmpfile.Name())
			assert.NoError(t, err)

			value := getFieldValue(tc.flag, tc.env, tc.file)
			assert.Equal(t, tc.expectedValue, value)
		})
	}
}

func TestFloat32Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []float32
	}{
		{
			[]string{},
			[]float32{},
		},
		{
			[]string{"3.1415"},
			[]float32{3.1415},
		},
		{
			[]string{"3.1415", "2.7182"},
			[]float32{3.1415, 2.7182},
		},
		{
			[]string{"3.1415", "2.7182", "1.6180"},
			[]float32{3.1415, 2.7182, 1.6180},
		},
	}

	for _, tc := range tests {
		result := float32Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestFloat64Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []float64
	}{
		{
			[]string{},
			[]float64{},
		},
		{
			[]string{"3.14159265"},
			[]float64{3.14159265},
		},
		{
			[]string{"3.14159265", "2.71828182"},
			[]float64{3.14159265, 2.71828182},
		},
		{
			[]string{"3.14159265", "2.71828182", "1.61803398"},
			[]float64{3.14159265, 2.71828182, 1.61803398},
		},
	}

	for _, tc := range tests {
		result := float64Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestIntSlice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []int
	}{
		{
			[]string{},
			[]int{},
		},
		{
			[]string{"-2147483648"},
			[]int{-2147483648},
		},
		{
			[]string{"-2147483648", "0"},
			[]int{-2147483648, 0},
		},
		{
			[]string{"-2147483648", "0", "2147483647"},
			[]int{-2147483648, 0, 2147483647},
		},
	}

	for _, tc := range tests {
		result := intSlice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestInt8Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []int8
	}{
		{
			[]string{},
			[]int8{},
		},
		{
			[]string{"-128"},
			[]int8{-128},
		},
		{
			[]string{"-128", "0"},
			[]int8{-128, 0},
		},
		{
			[]string{"-128", "0", "127"},
			[]int8{-128, 0, 127},
		},
	}

	for _, tc := range tests {
		result := int8Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestInt16Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []int16
	}{
		{
			[]string{},
			[]int16{},
		},
		{
			[]string{"-32768"},
			[]int16{-32768},
		},
		{
			[]string{"-32768", "0"},
			[]int16{-32768, 0},
		},
		{
			[]string{"-32768", "0", "32767"},
			[]int16{-32768, 0, 32767},
		},
	}

	for _, tc := range tests {
		result := int16Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestInt32Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []int32
	}{
		{
			[]string{},
			[]int32{},
		},
		{
			[]string{"-2147483648"},
			[]int32{-2147483648},
		},
		{
			[]string{"-2147483648", "0"},
			[]int32{-2147483648, 0},
		},
		{
			[]string{"-2147483648", "0", "2147483647"},
			[]int32{-2147483648, 0, 2147483647},
		},
	}

	for _, tc := range tests {
		result := int32Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestInt64Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []int64
	}{
		{
			[]string{},
			[]int64{},
		},
		{
			[]string{"-9223372036854775808"},
			[]int64{-9223372036854775808},
		},
		{
			[]string{"-9223372036854775808", "0"},
			[]int64{-9223372036854775808, 0},
		},
		{
			[]string{"-9223372036854775808", "0", "9223372036854775807"},
			[]int64{-9223372036854775808, 0, 9223372036854775807},
		},
	}

	for _, tc := range tests {
		result := int64Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestUintSlice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []uint
	}{
		{
			[]string{},
			[]uint{},
		},
		{
			[]string{"4294967295"},
			[]uint{4294967295},
		},
		{
			[]string{"0", "4294967295"},
			[]uint{0, 4294967295},
		},
		{
			[]string{"0", "2147483648", "4294967295"},
			[]uint{0, 2147483648, 4294967295},
		},
	}

	for _, tc := range tests {
		result := uintSlice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestUint8Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []uint8
	}{
		{
			[]string{},
			[]uint8{},
		},
		{
			[]string{"255"},
			[]uint8{255},
		},
		{
			[]string{"0", "255"},
			[]uint8{0, 255},
		},
		{
			[]string{"0", "128", "255"},
			[]uint8{0, 128, 255},
		},
	}

	for _, tc := range tests {
		result := uint8Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestUint16Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []uint16
	}{
		{
			[]string{},
			[]uint16{},
		},
		{
			[]string{"65535"},
			[]uint16{65535},
		},
		{
			[]string{"0", "65535"},
			[]uint16{0, 65535},
		},
		{
			[]string{"0", "32768", "65535"},
			[]uint16{0, 32768, 65535},
		},
	}

	for _, tc := range tests {
		result := uint16Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestUint32Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []uint32
	}{
		{
			[]string{},
			[]uint32{},
		},
		{
			[]string{"4294967295"},
			[]uint32{4294967295},
		},
		{
			[]string{"0", "4294967295"},
			[]uint32{0, 4294967295},
		},
		{
			[]string{"0", "2147483648", "4294967295"},
			[]uint32{0, 2147483648, 4294967295},
		},
	}

	for _, tc := range tests {
		result := uint32Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestUint64Slice(t *testing.T) {
	tests := []struct {
		strs     []string
		expected []uint64
	}{
		{
			[]string{},
			[]uint64{},
		},
		{
			[]string{"18446744073709551615"},
			[]uint64{18446744073709551615},
		},
		{
			[]string{"0", "18446744073709551615"},
			[]uint64{0, 18446744073709551615},
		},
		{
			[]string{"0", "9223372036854775808", "18446744073709551615"},
			[]uint64{0, 9223372036854775808, 18446744073709551615},
		},
	}

	for _, tc := range tests {
		result := uint64Slice(tc.strs)
		assert.Equal(t, tc.expected, result)
	}
}

func TestPickError(t *testing.T) {
	tests := []struct {
		name          string
		config        interface{}
		expectedError string
	}{
		{
			"NonPointer",
			Config{},
			"non-pointer type is passed",
		},
		{
			"NonStruct",
			new(string),
			"non-struct type is passed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Pick(tc.config)
			assert.Equal(t, tc.expectedError, err.Error())
		})
	}
}

func TestPick(t *testing.T) {
	d90m := 90 * time.Minute
	d120m := 120 * time.Minute
	exampleURL, _ := url.Parse("https://example.com")
	localhostURL, _ := url.Parse("http://localhost:8080")

	tests := []struct {
		name           string
		args           []string
		envs           [][2]string
		files          [][2]string
		config         Config
		expectedConfig Config
	}{
		{
			"Empty",
			[]string{},
			[][2]string{},
			[][2]string{},
			Config{},
			Config{},
		},
		{
			"AllFromDefaults",
			[]string{},
			[][2]string{},
			[][2]string{},
			Config{
				FieldString:        "default",
				FieldBool:          false,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *localhostURL,
				FieldStringArray:   []string{"url1", "url2"},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*exampleURL, *localhostURL},
			},
			Config{
				FieldString:        "default",
				FieldBool:          false,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *localhostURL,
				FieldStringArray:   []string{"url1", "url2"},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*exampleURL, *localhostURL},
			},
		},
		{
			"AllFromFlagsWithSingleDash",
			[]string{
				"-field.string", "content",
				"-field.bool",
				"-field.float32", "3.1415",
				"-field.float64", "3.14159265359",
				"-field.int", "-2147483648",
				"-field.int8", "-128",
				"-field.int16", "-32768",
				"-field.int32", "-2147483648",
				"-field.int64", "-9223372036854775808",
				"-field.uint", "4294967295",
				"-field.uint8", "255",
				"-field.uint16", "65535",
				"-field.uint32", "4294967295",
				"-field.uint64", "18446744073709551615",
				"-field.duration", "90m",
				"-field.url", "http://localhost:8080",
				"-field.string.array", "url1,url2",
				"-field.float32.array", "3.1415,2.7182",
				"-field.float64.array", "3.14159265359,2.71828182845",
				"-field.int.array", "-2147483648,2147483647",
				"-field.int8.array", "-128,127",
				"-field.int16.array", "-32768,32767",
				"-field.int32.array", "-2147483648,2147483647",
				"-field.int64.array", "-9223372036854775808,9223372036854775807",
				"-field.uint.array", "0,4294967295",
				"-field.uint8.array", "0,255",
				"-field.uint16.array", "0,65535",
				"-field.uint32.array", "0,4294967295",
				"-field.uint64.array", "0,18446744073709551615",
				"-field.duration.array", "90m,120m",
				"-field.url.array", "https://example.com,http://localhost:8080",
			},
			[][2]string{},
			[][2]string{},
			Config{},
			Config{
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *localhostURL,
				FieldStringArray:   []string{"url1", "url2"},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*exampleURL, *localhostURL},
			},
		},
		{
			"AllFromFlagsWithDoubleDash",
			[]string{
				"--field.string", "content",
				"--field.bool",
				"--field.float32", "3.1415",
				"--field.float64", "3.14159265359",
				"--field.int", "-2147483648",
				"--field.int8", "-128",
				"--field.int16", "-32768",
				"--field.int32", "-2147483648",
				"--field.int64", "-9223372036854775808",
				"--field.uint", "4294967295",
				"--field.uint8", "255",
				"--field.uint16", "65535",
				"--field.uint32", "4294967295",
				"--field.uint64", "18446744073709551615",
				"--field.duration", "90m",
				"--field.url", "http://localhost:8080",
				"--field.string.array", "url1,url2",
				"--field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
				"--field.int.array", "-2147483648,2147483647",
				"--field.int8.array", "-128,127",
				"--field.int16.array", "-32768,32767",
				"--field.int32.array", "-2147483648,2147483647",
				"--field.int64.array", "-9223372036854775808,9223372036854775807",
				"--field.uint.array", "0,4294967295",
				"--field.uint8.array", "0,255",
				"--field.uint16.array", "0,65535",
				"--field.uint32.array", "0,4294967295",
				"--field.uint64.array", "0,18446744073709551615",
				"--field.duration.array", "90m,120m",
				"--field.url.array", "https://example.com,http://localhost:8080",
			},
			[][2]string{},
			[][2]string{},
			Config{},
			Config{
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *localhostURL,
				FieldStringArray:   []string{"url1", "url2"},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*exampleURL, *localhostURL},
			},
		},
		{
			"AllFromFlagsWithSingleDashAndAssignment",
			[]string{
				"-field.string=content",
				"-field.bool",
				"-field.float32=3.1415",
				"-field.float64=3.14159265359",
				"-field.int=-2147483648",
				"-field.int8=-128",
				"-field.int16=-32768",
				"-field.int32=-2147483648",
				"-field.int64=-9223372036854775808",
				"-field.uint=4294967295",
				"-field.uint8=255",
				"-field.uint16=65535",
				"-field.uint32=4294967295",
				"-field.uint64=18446744073709551615",
				"-field.duration=90m",
				"-field.url=http://localhost:8080",
				"-field.string.array=url1,url2",
				"-field.float32.array=3.1415,2.7182",
				"-field.float64.array=3.14159265359,2.71828182845",
				"-field.int.array=-2147483648,2147483647",
				"-field.int8.array=-128,127",
				"-field.int16.array=-32768,32767",
				"-field.int32.array=-2147483648,2147483647",
				"-field.int64.array=-9223372036854775808,9223372036854775807",
				"-field.uint.array=0,4294967295",
				"-field.uint8.array=0,255",
				"-field.uint16.array=0,65535",
				"-field.uint32.array=0,4294967295",
				"-field.uint64.array=0,18446744073709551615",
				"-field.duration.array=90m,120m",
				"-field.url.array=https://example.com,http://localhost:8080",
			},
			[][2]string{},
			[][2]string{},
			Config{},
			Config{
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *localhostURL,
				FieldStringArray:   []string{"url1", "url2"},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*exampleURL, *localhostURL},
			},
		},
		{
			"AllFromFlagsWithDoubleDashAndAssignment",
			[]string{
				"--field.string=content",
				"--field.bool",
				"--field.float32=3.1415",
				"--field.float64=3.14159265359",
				"--field.int=-2147483648",
				"--field.int8=-128",
				"--field.int16=-32768",
				"--field.int32=-2147483648",
				"--field.int64=-9223372036854775808",
				"--field.uint=4294967295",
				"--field.uint8=255",
				"--field.uint16=65535",
				"--field.uint32=4294967295",
				"--field.uint64=18446744073709551615",
				"--field.duration=90m",
				"--field.url=http://localhost:8080",
				"--field.string.array=url1,url2",
				"--field.float32.array=3.1415,2.7182",
				"--field.float64.array=3.14159265359,2.71828182845",
				"--field.int.array=-2147483648,2147483647",
				"--field.int8.array=-128,127",
				"--field.int16.array=-32768,32767",
				"--field.int32.array=-2147483648,2147483647",
				"--field.int64.array=-9223372036854775808,9223372036854775807",
				"--field.uint.array=0,4294967295",
				"--field.uint8.array=0,255",
				"--field.uint16.array=0,65535",
				"--field.uint32.array=0,4294967295",
				"--field.uint64.array=0,18446744073709551615",
				"--field.duration.array=90m,120m",
				"--field.url.array=https://example.com,http://localhost:8080",
			},
			[][2]string{},
			[][2]string{},
			Config{},
			Config{
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *localhostURL,
				FieldStringArray:   []string{"url1", "url2"},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*exampleURL, *localhostURL},
			},
		},
		{
			"AllFromEnvironmentVariables",
			[]string{},
			[][2]string{
				[2]string{"FIELD_STRING", "content"},
				[2]string{"FIELD_BOOL", "true"},
				[2]string{"FIELD_FLOAT32", "3.1415"},
				[2]string{"FIELD_FLOAT64", "3.14159265359"},
				[2]string{"FIELD_INT", "-2147483648"},
				[2]string{"FIELD_INT8", "-128"},
				[2]string{"FIELD_INT16", "-32768"},
				[2]string{"FIELD_INT32", "-2147483648"},
				[2]string{"FIELD_INT64", "-9223372036854775808"},
				[2]string{"FIELD_UINT", "4294967295"},
				[2]string{"FIELD_UINT8", "255"},
				[2]string{"FIELD_UINT16", "65535"},
				[2]string{"FIELD_UINT32", "4294967295"},
				[2]string{"FIELD_UINT64", "18446744073709551615"},
				[2]string{"FIELD_DURATION", "90m"},
				[2]string{"FIELD_URL", "http://localhost:8080"},
				[2]string{"FIELD_STRING_ARRAY", "url1,url2"},
				[2]string{"FIELD_FLOAT32_ARRAY", "3.1415,2.7182"},
				[2]string{"FIELD_FLOAT64_ARRAY", "3.14159265359,2.71828182845"},
				[2]string{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				[2]string{"FIELD_INT8_ARRAY", "-128,127"},
				[2]string{"FIELD_INT16_ARRAY", "-32768,32767"},
				[2]string{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				[2]string{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
				[2]string{"FIELD_UINT_ARRAY", "0,4294967295"},
				[2]string{"FIELD_UINT8_ARRAY", "0,255"},
				[2]string{"FIELD_UINT16_ARRAY", "0,65535"},
				[2]string{"FIELD_UINT32_ARRAY", "0,4294967295"},
				[2]string{"FIELD_UINT64_ARRAY", "0,18446744073709551615"},
				[2]string{"FIELD_DURATION_ARRAY", "90m,120m"},
				[2]string{"FIELD_URL_ARRAY", "https://example.com,http://localhost:8080"},
			},
			[][2]string{},
			Config{},
			Config{
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *localhostURL,
				FieldStringArray:   []string{"url1", "url2"},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*exampleURL, *localhostURL},
			},
		},
		{
			"AllFromFromFileContent",
			[]string{},
			[][2]string{},
			[][2]string{
				[2]string{"FIELD_STRING_FILE", "content"},
				[2]string{"FIELD_BOOL_FILE", "true"},
				[2]string{"FIELD_FLOAT32_FILE", "3.1415"},
				[2]string{"FIELD_FLOAT64_FILE", "3.14159265359"},
				[2]string{"FIELD_INT_FILE", "-2147483648"},
				[2]string{"FIELD_INT8_FILE", "-128"},
				[2]string{"FIELD_INT16_FILE", "-32768"},
				[2]string{"FIELD_INT32_FILE", "-2147483648"},
				[2]string{"FIELD_INT64_FILE", "-9223372036854775808"},
				[2]string{"FIELD_UINT_FILE", "4294967295"},
				[2]string{"FIELD_UINT8_FILE", "255"},
				[2]string{"FIELD_UINT16_FILE", "65535"},
				[2]string{"FIELD_UINT32_FILE", "4294967295"},
				[2]string{"FIELD_UINT64_FILE", "18446744073709551615"},
				[2]string{"FIELD_DURATION_FILE", "90m"},
				[2]string{"FIELD_URL_FILE", "http://localhost:8080"},
				[2]string{"FIELD_STRING_ARRAY_FILE", "url1,url2"},
				[2]string{"FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182"},
				[2]string{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845"},
				[2]string{"FIELD_INT_ARRAY_FILE", "-2147483648,2147483647"},
				[2]string{"FIELD_INT8_ARRAY_FILE", "-128,127"},
				[2]string{"FIELD_INT16_ARRAY_FILE", "-32768,32767"},
				[2]string{"FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647"},
				[2]string{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807"},
				[2]string{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				[2]string{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				[2]string{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				[2]string{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				[2]string{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				[2]string{"FIELD_DURATION_ARRAY_FILE", "90m,120m"},
				[2]string{"FIELD_URL_ARRAY_FILE", "https://example.com,http://localhost:8080"},
			},
			Config{},
			Config{
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *localhostURL,
				FieldStringArray:   []string{"url1", "url2"},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*exampleURL, *localhostURL},
			},
		},
		{
			"Mixed",
			[]string{
				"-field.bool",
				"-field.float32=3.1415",
				"--field.float64=3.14159265359",
				"-field.duration=90m",
				"--field.url", "http://localhost:8080",
				"-field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
			},
			[][2]string{
				[2]string{"FIELD_INT", "-2147483648"},
				[2]string{"FIELD_INT8", "-128"},
				[2]string{"FIELD_INT16", "-32768"},
				[2]string{"FIELD_INT32", "-2147483648"},
				[2]string{"FIELD_INT64", "-9223372036854775808"},
				[2]string{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				[2]string{"FIELD_INT8_ARRAY", "-128,127"},
				[2]string{"FIELD_INT16_ARRAY", "-32768,32767"},
				[2]string{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				[2]string{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
				[2]string{"FIELD_DURATION_ARRAY", "90m,120m"},
			},
			[][2]string{
				[2]string{"FIELD_UINT_FILE", "4294967295"},
				[2]string{"FIELD_UINT8_FILE", "255"},
				[2]string{"FIELD_UINT16_FILE", "65535"},
				[2]string{"FIELD_UINT32_FILE", "4294967295"},
				[2]string{"FIELD_UINT64_FILE", "18446744073709551615"},
				[2]string{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				[2]string{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				[2]string{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				[2]string{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				[2]string{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				[2]string{"FIELD_URL_ARRAY_FILE", "https://example.com,http://localhost:8080"},
			},
			Config{
				FieldString:      "default",
				FieldStringArray: []string{"url1", "url2"},
			},
			Config{
				FieldString:        "default",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *localhostURL,
				FieldStringArray:   []string{"url1", "url2"},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*exampleURL, *localhostURL},
			},
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set arguments
			os.Args = tc.args

			// Set environment variables
			for _, env := range tc.envs {
				err := os.Setenv(env[0], env[1])
				assert.NoError(t, err)
				defer os.Unsetenv(env[0])
			}

			// Write files
			for _, file := range tc.files {
				tmpfile, err := ioutil.TempFile("", "gotest_")
				assert.NoError(t, err)
				defer os.Remove(tmpfile.Name())
				_, err = tmpfile.WriteString(file[1])
				assert.NoError(t, err)
				err = tmpfile.Close()
				assert.NoError(t, err)
				err = os.Setenv(file[0], tmpfile.Name())
				assert.NoError(t, err)
				defer os.Unsetenv(file[0])
			}

			err := Pick(&tc.config)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedConfig, tc.config)
		})
	}

	// flag.Parse() can be called only once
	flag.Parse()
}
