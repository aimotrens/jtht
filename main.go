package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

var writeHeaderOnce = false
var headerWritten = map[string]bool{}

func main() {
	var sf, df string

	flag.StringVar(&sf, "src", "", "Source json file to convert")
	flag.StringVar(&df, "dest", "<src>.html", "Destination html file")
	flag.BoolVar(&writeHeaderOnce, "who", false, "If true, header are only written once")

	flag.Parse()

	if sf == "" {
		fmt.Println("Argument \"src\" must be set.")
		os.Exit(2)
	}

	if df == "<src>.html" {
		df = sf + ".html"
	}

	fmt.Printf("Converting %s to %s...\n\n", sf, df)

	convert(sf, df)
}

func convert(src, dst string) {
	fileContent, err := os.ReadFile(src)
	if err != nil {
		panic("Error reading file: " + err.Error())
	}

	var jsonData any
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		panic("Could not unmarshal json:" + err.Error())
	}

	jsonStruct := jsonData.(map[string]any)
	fmt.Printf("Struct: %v\n\n", jsonStruct)

	var sb strings.Builder

	sb.WriteString("<!doctype html>\n")
	sb.WriteString("<html><head><meta charset=\"utf-8\"><link rel=\"stylesheet\" href=\"./style.css\" /></head><body>")

	processObject("", &jsonStruct, &sb)

	sb.WriteString("</body></html>")

	//fmt.Printf("\nsb: %v\n", sb.String())

	os.WriteFile(dst, []byte(sb.String()), os.ModePerm)
}

func processObject(key string, object *map[string]any, output *strings.Builder) {
	output.WriteString("<table>")

	// Sort headers
	headers := make([]string, 0)
	for key := range *object {
		headers = append(headers, key)
	}

	sort.Strings(headers)

	if _, ok := headerWritten[key]; !ok || !writeHeaderOnce {
		headerWritten[key] = true

		// Column Header
		output.WriteString("<tr>")
		for _, key := range headers {
			output.WriteString("<th>")
			output.WriteString(key)
			output.WriteString("</th>")
		}
		output.WriteString("</tr>")
	}

	// Content
	output.WriteString("<tr>")
	for _, key := range headers {
		value := (*object)[key]

		output.WriteString("<td>")

		switch val := value.(type) {
		case map[string]any:
			processObject(key, &val, output)
		case []any:
			processArray(key, &val, output)
		default:
			processValue(&val, output)
		}

		output.WriteString("</td>")
	}
	output.WriteString("</tr>")

	output.WriteString("</table>")
}

func processArray(key string, array *[]any, output *strings.Builder) {
	for _, value := range *array {

		switch val := value.(type) {
		case map[string]any:
			processObject(key, &val, output)
		default:
			panic("Can only process array of json objects.")
		}
	}
}

func processValue(value *any, output *strings.Builder) {
	switch (*value).(type) {
	case nil:
		output.WriteString("NULL")
	case string:
		output.WriteString((*value).(string))
	case float64:
		output.WriteString(fmt.Sprint(*value))
	case bool:
		output.WriteString(fmt.Sprint(*value))
	default:
		fmt.Printf("Unknown type: %s / Value: %s\n", reflect.TypeOf(value).String(), *value)
	}
}

// func getMaxPropertyCount(jsonStruct *any) int {
// 	switch (*jsonStruct).(type) {
// 	case map[string]any:
// 		return len((*jsonStruct).(map[string]any))

// 	case []any:
// 		max := 0
// 		for _, item := range (*jsonStruct).([]any) {
// 			childMax := getMaxPropertyCount(&item)

// 			if childMax > max {
// 				max = childMax
// 			}
// 		}
// 		return max
// 	default:
// 		return 1
// 	}
// }
