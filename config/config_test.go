package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSpec struct {
	FieldString  string  `flag:"field.string" env:"FIELD_STRING" file:"FIELD_STRING_FILE"`
	FieldBool    bool    `flag:"field.bool" env:"FIELD_BOOL" file:"FIELD_BOOL_FILE"`
	FieldFloat32 float32 `flag:"field.float32" env:"FIELD_FLOAT32" file:"FIELD_FLOAT32_FILE"`
	FieldFloat64 float64 `flag:"field.float64" env:"FIELD_FLOAT64" file:"FIELD_FLOAT64_FILE"`
	FieldInt     int     `flag:"field.int" env:"FIELD_INT" file:"FIELD_INT_FILE"`
	FieldInt8    int8    `flag:"field.int8" env:"FIELD_INT8" file:"FIELD_INT8_FILE"`
	FieldInt16   int16   `flag:"field.int16" env:"FIELD_INT16" file:"FIELD_INT16_FILE"`
	FieldInt32   int32   `flag:"field.int32" env:"FIELD_INT32" file:"FIELD_INT32_FILE"`
	FieldInt64   int64   `flag:"field.int64" env:"FIELD_INT64" file:"FIELD_INT64_FILE"`
	FieldUint    uint    `flag:"field.uint" env:"FIELD_UINT" file:"FIELD_UINT_FILE"`
	FieldUint8   uint8   `flag:"field.uint8" env:"FIELD_UINT8" file:"FIELD_UINT8_FILE"`
	FieldUint16  uint16  `flag:"field.uint16" env:"FIELD_UINT16" file:"FIELD_UINT16_FILE"`
	FieldUint32  uint32  `flag:"field.uint32" env:"FIELD_UINT32" file:"FIELD_UINT32_FILE"`
	FieldUint64  uint64  `flag:"field.uint64" env:"FIELD_UINT64" file:"FIELD_UINT64_FILE"`
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
		{"MoreCamelCase", "more.camel.case"},
		{"DatabaseURL", "database.url"},
	}

	for _, tc := range tests {
		result := getFlagName(tc.fieldName)
		assert.Equal(t, tc.expectedFlagName, result)
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
		{"MoreCamelCase", "MORE_CAMEL_CASE"},
		{"DatabaseURL", "DATABASE_URL"},
	}

	for _, tc := range tests {
		result := getEnvVarName(tc.fieldName)
		assert.Equal(t, tc.expectedEnvVarName, result)
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

func TestPick(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		envs         [][2]string
		files        [][2]string
		spec         testSpec
		expectedSpec testSpec
	}{
		{
			"Empty",
			[]string{},
			[][2]string{},
			[][2]string{},
			testSpec{},
			testSpec{},
		},
		{
			"AllFromDefaults",
			[]string{},
			[][2]string{},
			[][2]string{},
			testSpec{
				FieldString:  "default",
				FieldBool:    false,
				FieldFloat32: 3.1415,
				FieldFloat64: 3.14159265359,
				FieldInt:     -2147483648,
				FieldInt8:    -128,
				FieldInt16:   -32768,
				FieldInt32:   -2147483648,
				FieldInt64:   -9223372036854775808,
				FieldUint:    4294967295,
				FieldUint8:   255,
				FieldUint16:  65535,
				FieldUint32:  4294967295,
				FieldUint64:  18446744073709551615,
			},
			testSpec{
				FieldString:  "default",
				FieldBool:    false,
				FieldFloat32: 3.1415,
				FieldFloat64: 3.14159265359,
				FieldInt:     -2147483648,
				FieldInt8:    -128,
				FieldInt16:   -32768,
				FieldInt32:   -2147483648,
				FieldInt64:   -9223372036854775808,
				FieldUint:    4294967295,
				FieldUint8:   255,
				FieldUint16:  65535,
				FieldUint32:  4294967295,
				FieldUint64:  18446744073709551615,
			},
		},
		{
			"AllFromFlags#01",
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
			},
			[][2]string{},
			[][2]string{},
			testSpec{},
			testSpec{
				FieldString:  "content",
				FieldBool:    true,
				FieldFloat32: 3.1415,
				FieldFloat64: 3.14159265359,
				FieldInt:     -2147483648,
				FieldInt8:    -128,
				FieldInt16:   -32768,
				FieldInt32:   -2147483648,
				FieldInt64:   -9223372036854775808,
				FieldUint:    4294967295,
				FieldUint8:   255,
				FieldUint16:  65535,
				FieldUint32:  4294967295,
				FieldUint64:  18446744073709551615,
			},
		},
		{
			"AllFromFlags#02",
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
			},
			[][2]string{},
			[][2]string{},
			testSpec{},
			testSpec{
				FieldString:  "content",
				FieldBool:    true,
				FieldFloat32: 3.1415,
				FieldFloat64: 3.14159265359,
				FieldInt:     -2147483648,
				FieldInt8:    -128,
				FieldInt16:   -32768,
				FieldInt32:   -2147483648,
				FieldInt64:   -9223372036854775808,
				FieldUint:    4294967295,
				FieldUint8:   255,
				FieldUint16:  65535,
				FieldUint32:  4294967295,
				FieldUint64:  18446744073709551615,
			},
		},
		{
			"AllFromFlags#03",
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
			},
			[][2]string{},
			[][2]string{},
			testSpec{},
			testSpec{
				FieldString:  "content",
				FieldBool:    true,
				FieldFloat32: 3.1415,
				FieldFloat64: 3.14159265359,
				FieldInt:     -2147483648,
				FieldInt8:    -128,
				FieldInt16:   -32768,
				FieldInt32:   -2147483648,
				FieldInt64:   -9223372036854775808,
				FieldUint:    4294967295,
				FieldUint8:   255,
				FieldUint16:  65535,
				FieldUint32:  4294967295,
				FieldUint64:  18446744073709551615,
			},
		},
		{
			"AllFromFlags#04",
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
			},
			[][2]string{},
			[][2]string{},
			testSpec{},
			testSpec{
				FieldString:  "content",
				FieldBool:    true,
				FieldFloat32: 3.1415,
				FieldFloat64: 3.14159265359,
				FieldInt:     -2147483648,
				FieldInt8:    -128,
				FieldInt16:   -32768,
				FieldInt32:   -2147483648,
				FieldInt64:   -9223372036854775808,
				FieldUint:    4294967295,
				FieldUint8:   255,
				FieldUint16:  65535,
				FieldUint32:  4294967295,
				FieldUint64:  18446744073709551615,
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
			},
			[][2]string{},
			testSpec{},
			testSpec{
				FieldString:  "content",
				FieldBool:    true,
				FieldFloat32: 3.1415,
				FieldFloat64: 3.14159265359,
				FieldInt:     -2147483648,
				FieldInt8:    -128,
				FieldInt16:   -32768,
				FieldInt32:   -2147483648,
				FieldInt64:   -9223372036854775808,
				FieldUint:    4294967295,
				FieldUint8:   255,
				FieldUint16:  65535,
				FieldUint32:  4294967295,
				FieldUint64:  18446744073709551615,
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
				[2]string{"FIELD_UINT64_FILE", "18446744073709551615"}},
			testSpec{},
			testSpec{
				FieldString:  "content",
				FieldBool:    true,
				FieldFloat32: 3.1415,
				FieldFloat64: 3.14159265359,
				FieldInt:     -2147483648,
				FieldInt8:    -128,
				FieldInt16:   -32768,
				FieldInt32:   -2147483648,
				FieldInt64:   -9223372036854775808,
				FieldUint:    4294967295,
				FieldUint8:   255,
				FieldUint16:  65535,
				FieldUint32:  4294967295,
				FieldUint64:  18446744073709551615,
			},
		},
		{
			"Combinatorial",
			[]string{
				"-field.bool",
				"-field.float32=3.1415",
				"--field.float64", "3.14159265359",
			},
			[][2]string{
				[2]string{"FIELD_INT", "-2147483648"},
				[2]string{"FIELD_INT8", "-128"},
				[2]string{"FIELD_INT16", "-32768"},
				[2]string{"FIELD_INT32", "-2147483648"},
				[2]string{"FIELD_INT64", "-9223372036854775808"},
			},
			[][2]string{
				[2]string{"FIELD_UINT_FILE", "4294967295"},
				[2]string{"FIELD_UINT8_FILE", "255"},
				[2]string{"FIELD_UINT16_FILE", "65535"},
				[2]string{"FIELD_UINT32_FILE", "4294967295"},
				[2]string{"FIELD_UINT64_FILE", "18446744073709551615"}},
			testSpec{
				FieldString: "default",
			},
			testSpec{
				FieldString:  "default",
				FieldBool:    true,
				FieldFloat32: 3.1415,
				FieldFloat64: 3.14159265359,
				FieldInt:     -2147483648,
				FieldInt8:    -128,
				FieldInt16:   -32768,
				FieldInt32:   -2147483648,
				FieldInt64:   -9223372036854775808,
				FieldUint:    4294967295,
				FieldUint8:   255,
				FieldUint16:  65535,
				FieldUint32:  4294967295,
				FieldUint64:  18446744073709551615,
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

			err := Pick(&tc.spec)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedSpec, tc.spec)
		})
	}
}
