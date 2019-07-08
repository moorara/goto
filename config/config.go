package config

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// this is used for printing debugging logs
var debug bool

const (
	flagTag   = "flag"
	envTag    = "env"
	fileTag   = "file"
	sepTag    = "sep"
	skipValue = "-"
)

type flagValue struct{}

func (v *flagValue) String() string {
	return ""
}

func (v *flagValue) Set(string) error {
	return nil
}

func print(msg string, args ...interface{}) {
	if debug {
		log.Printf(msg+"\n", args...)
	}
}

/*
 * tokenize breaks a field name into its tokens (generally words).
 *   UserID       -->  User, ID
 *   DatabaseURL  -->  Database, URL
 */
func tokenize(name string) []string {
	tokens := []string{}
	current := string(name[0])
	lastLower := unicode.IsLower(rune(name[0]))

	add := func(slice []string, str string) []string {
		if str == "" {
			return slice
		}
		return append(slice, str)
	}

	for i := 1; i < len(name); i++ {
		r := rune(name[i])

		if unicode.IsUpper(r) && lastLower {
			// The case is changing from lower to upper
			tokens = add(tokens, current)
			current = string(name[i])
		} else if unicode.IsLower(r) && !lastLower {
			// The case is changing from upper to lower
			l := len(current) - 1
			tokens = add(tokens, current[:l])
			current = current[l:] + string(name[i])
		} else {
			// Increment current token
			current += string(name[i])
		}

		lastLower = unicode.IsLower(r)
	}

	tokens = append(tokens, string(current))

	return tokens
}

/*
 * getFlagName returns a canonical flag name for a field.
 *   UserID       -->  user.id
 *   DatabaseURL  -->  database.url
 */
func getFlagName(name string) string {
	parts := tokenize(name)
	result := strings.Join(parts, ".")
	result = strings.ToLower(result)

	return result
}

/*
 * getFlagName returns a canonical environment variable name for a field.
 *   UserID       -->  USER_ID
 *   DatabaseURL  -->  DATABASE_URL
 */
func getEnvVarName(name string) string {
	parts := tokenize(name)
	result := strings.Join(parts, "_")
	result = strings.ToUpper(result)

	return result
}

/*
 * getFileVarName returns a canonical environment variable name for value file of a field.
 *   UserID       -->  USER_ID_FILE
 *   DatabaseURL  -->  DATABASE_URL_FILE
 */
func getFileVarName(name string) string {
	parts := tokenize(name)
	result := strings.Join(parts, "_")
	result = strings.ToUpper(result)
	result = result + "_FILE"

	return result
}

// defineFlag registers a flag name, so it will show up in the help description.
func defineFlag(flagName, defaultValue, envName, fileName string) {
	if flagName == skipValue {
		return
	}

	usage := fmt.Sprintf(
		"%s:\t\t\t\t%s\n%s:\t\t\t%s\n%s:\t%s",
		"default value", defaultValue,
		"environment variable", envName,
		"config file environment variable", fileName,
	)

	if flag.Lookup(flagName) == nil {
		flag.Var(&flagValue{}, flagName, usage)
	}
}

/*
 * getFlagValue returns the value set for a flag.
 *   - The flag name can start with - or --
 *   - The flag value can be separated by space or =
 */
func getFlagValue(flagName string) string {
	flagRegex := regexp.MustCompile("-{1,2}" + flagName)
	genericRegex := regexp.MustCompile("^-{1,2}[A-Za-z].*")

	for i, arg := range os.Args {
		if flagRegex.MatchString(arg) {
			if s := strings.Index(arg, "="); s > 0 {
				return arg[s+1:]
			}

			if i+1 < len(os.Args) {
				val := os.Args[i+1]
				if !genericRegex.MatchString(val) {
					return val
				}
			}

			return "true"
		}
	}

	return ""
}

/*
 * getFieldValue reads and returns the string value for a field from either
 *   - command-line flags,
 *   - environment variables,
 *   - or configuration files
 */
func getFieldValue(flag, env, file string) string {
	var value string

	// First, try reading from flag
	if value == "" && flag != skipValue {
		value = getFlagValue(flag)
		print("value read from flag %s: %s", flag, value)
	}

	// Second, try reading from environment variable
	if value == "" && env != skipValue {
		value = os.Getenv(env)
		print("value read from environment variable %s: %s", env, value)
	}

	// Third, try reading from file
	if value == "" && file != skipValue {
		filepath := os.Getenv(file)
		if content, err := ioutil.ReadFile(filepath); err == nil {
			value = string(content)
			print("value read from file %s: %s", file, value)
		}
	}

	return value
}

func float32Slice(strs []string) []float32 {
	floats := []float32{}
	for _, str := range strs {
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			floats = append(floats, float32(f))
		}
	}
	return floats
}

func float64Slice(strs []string) []float64 {
	floats := []float64{}
	for _, str := range strs {
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			floats = append(floats, f)
		}
	}
	return floats
}

func intSlice(strs []string) []int {
	ints := []int{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, int(i))
		}
	}
	return ints
}

func int8Slice(strs []string) []int8 {
	ints := []int8{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, int8(i))
		}
	}
	return ints
}

func int16Slice(strs []string) []int16 {
	ints := []int16{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, int16(i))
		}
	}
	return ints
}

func int32Slice(strs []string) []int32 {
	ints := []int32{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, int32(i))
		}
	}
	return ints
}

func int64Slice(strs []string) []int64 {
	ints := []int64{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, i)
		}
	}
	return ints
}

func uintSlice(strs []string) []uint {
	uints := []uint{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, uint(u))
		}
	}
	return uints
}

func uint8Slice(strs []string) []uint8 {
	uints := []uint8{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, uint8(u))
		}
	}
	return uints
}

func uint16Slice(strs []string) []uint16 {
	uints := []uint16{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, uint16(u))
		}
	}
	return uints
}

func uint32Slice(strs []string) []uint32 {
	uints := []uint32{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, uint32(u))
		}
	}
	return uints
}

func uint64Slice(strs []string) []uint64 {
	uints := []uint64{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, u)
		}
	}
	return uints
}

func durationSlice(strs []string) []time.Duration {
	durations := []time.Duration{}
	for _, str := range strs {
		if d, err := time.ParseDuration(str); err == nil {
			durations = append(durations, d)
		}
	}
	return durations
}

func urlSlice(strs []string) []url.URL {
	urls := []url.URL{}
	for _, str := range strs {
		if u, err := url.Parse(str); err == nil {
			urls = append(urls, *u)
		}
	}
	return urls
}

func pick(config interface{}) error {
	v := reflect.ValueOf(config) // reflect.Value --> v.Type(), v.Kind(), v.NumField()
	t := reflect.TypeOf(config)  // reflect.Type --> t.Name(), t.Kind(), t.NumField()

	// If a pointer is passed, navigate to the value
	if t.Kind() != reflect.Ptr {
		print("a non-pointer type is passed")
		return errors.New("a non-pointer type is passed")
	}

	// Navigate to the pointer value
	v = v.Elem()
	t = t.Elem()

	if t.Kind() != reflect.Struct {
		print("a non-struct type is passed")
		return errors.New("a non-struct type is passed")
	}

	// Iterate over struct fields
	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i) // reflect.Value --> vField.Kind(), vField.Type().Name(), vField.Type().Kind(), vField.Interface()
		tField := t.Field(i) // reflect.StructField --> tField.Name, tField.Type.Name(), tField.Type.Kind(), tField.Tag.Get(tag)

		// Skip unexported fields
		if !vField.CanSet() {
			continue
		}

		name := tField.Name

		// `flag:"..."`
		flagName := tField.Tag.Get(flagTag)
		if flagName == "" {
			flagName = getFlagName(name)
		}

		print("[%s] expecting flag name: %s", name, flagName)

		// `env:"..."`
		envName := tField.Tag.Get(envTag)
		if envName == "" {
			envName = getEnvVarName(name)
		}

		print("[%s] expecting environment variable name: %s", name, envName)

		// `file:"..."`
		fileName := tField.Tag.Get(fileTag)
		if fileName == "" {
			fileName = getFileVarName(name)
		}

		print("[%s] expecting file environment variable name: %s", name, fileName)

		// `sep:"..."`
		sep := tField.Tag.Get(sepTag)
		if sep == "" {
			sep = ","
		}

		print("[%s] expecting separator for list: %s", name, sep)

		// Define a flag for the field so flag.Parse() can be called
		defaultValue := fmt.Sprintf("%v", vField.Interface())
		defineFlag(flagName, defaultValue, envName, fileName)

		str := getFieldValue(flagName, envName, fileName)
		if str == "" {
			continue
		}

		switch vField.Kind() {
		case reflect.String:
			print("[%s] setting string value: %s", name, str)
			vField.SetString(str)

		case reflect.Bool:
			if b, err := strconv.ParseBool(str); err == nil {
				print("[%s] setting boolean value: %t", name, b)
				vField.SetBool(b)
			}

		case reflect.Float32, reflect.Float64:
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				print("[%s] setting float value: %f", name, f)
				vField.SetFloat(f)
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if t := vField.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
				// time.Duration
				if d, err := time.ParseDuration(str); err == nil {
					print("[%s] setting duration value: %s", name, d)
					vField.Set(reflect.ValueOf(d))
				}
			} else if i, err := strconv.ParseInt(str, 10, 64); err == nil {
				print("[%s] setting integer value: %d", name, i)
				vField.SetInt(i)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if u, err := strconv.ParseUint(str, 10, 64); err == nil {
				print("[%s] setting unsigned integer value: %d", name, u)
				vField.SetUint(u)
			}

		case reflect.Struct:
			if t := vField.Type(); t.PkgPath() == "net/url" && t.Name() == "URL" {
				// url.URL
				if u, err := url.Parse(str); err == nil {
					print("[%s] setting url value: %s", name, str)
					// u is a pointer
					vField.Set(reflect.ValueOf(u).Elem())
				}
			}

		case reflect.Slice:
			iSlice := vField.Interface()
			tSlice := reflect.TypeOf(iSlice).Elem()
			strs := strings.Split(str, sep)

			switch tSlice.Kind() {
			case reflect.String:
				print("[%s] setting string slice: %v", name, str)
				vField.Set(reflect.ValueOf(strs))

			case reflect.Float32:
				floats := float32Slice(strs)
				print("[%s] setting float32 slice: %v", name, floats)
				vField.Set(reflect.ValueOf(floats))

			case reflect.Float64:
				floats := float64Slice(strs)
				print("[%s] setting float64 slice: %v", name, floats)
				vField.Set(reflect.ValueOf(floats))

			case reflect.Int:
				ints := intSlice(strs)
				print("[%s] setting int slice: %v", name, ints)
				vField.Set(reflect.ValueOf(ints))

			case reflect.Int8:
				ints := int8Slice(strs)
				print("[%s] setting int8 slice: %v", name, ints)
				vField.Set(reflect.ValueOf(ints))

			case reflect.Int16:
				ints := int16Slice(strs)
				print("[%s] setting int16 slice: %v", name, ints)
				vField.Set(reflect.ValueOf(ints))

			case reflect.Int32:
				ints := int32Slice(strs)
				print("[%s] setting int32 slice: %v", name, ints)
				vField.Set(reflect.ValueOf(ints))

			case reflect.Int64:
				if tSlice.PkgPath() == "time" && tSlice.Name() == "Duration" {
					// []time.Duration
					durations := durationSlice(strs)
					print("[%s] setting duration slice: %v", name, durations)
					vField.Set(reflect.ValueOf(durations))
				} else {
					ints := int64Slice(strs)
					print("[%s] setting int64 slice: %v", name, ints)
					vField.Set(reflect.ValueOf(ints))
				}

			case reflect.Uint:
				uints := uintSlice(strs)
				print("[%s] setting uint slice: %v", name, uints)
				vField.Set(reflect.ValueOf(uints))

			case reflect.Uint8:
				uints := uint8Slice(strs)
				print("[%s] setting uint8 slice: %v", name, uints)
				vField.Set(reflect.ValueOf(uints))

			case reflect.Uint16:
				uints := uint16Slice(strs)
				print("[%s] setting uint16 slice: %v", name, uints)
				vField.Set(reflect.ValueOf(uints))

			case reflect.Uint32:
				uints := uint32Slice(strs)
				print("[%s] setting uint32 slice: %v", name, uints)
				vField.Set(reflect.ValueOf(uints))

			case reflect.Uint64:
				uints := uint64Slice(strs)
				print("[%s] setting uint64 slice: %v", name, uints)
				vField.Set(reflect.ValueOf(uints))

			case reflect.Struct:
				if tSlice.PkgPath() == "net/url" && tSlice.Name() == "URL" {
					// []url.URL
					urls := urlSlice(strs)
					print("[%s] setting url slice: %v", name, urls)
					vField.Set(reflect.ValueOf(urls))
				}
			}
		}
	}

	return nil
}

// Pick reads values for exported fields of a struct from either command-line flags, environment variables, or configuration files.
// You can also specify default values.
// You can see examples at https://github.com/moorara/goto/tree/master/config
func Pick(config interface{}) error {
	debug = false
	return pick(config)
}

// PickAndLog is same as Pick, but it also logs debugging information.
// You can also specify default values.
// You can see examples at https://github.com/moorara/goto/tree/master/config
func PickAndLog(config interface{}) error {
	debug = true
	return pick(config)
}
