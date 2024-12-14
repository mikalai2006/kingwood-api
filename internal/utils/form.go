package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BindAndValid binds and validates data.
func BindAndValid[V any](c *gin.Context, form V) (interface{}, error) {
	var body map[string]interface{}
	if c.Request.Body == nil {
		return nil, fmt.Errorf("not found data for patch")
	}
	if er := json.NewDecoder(c.Request.Body).Decode(&body); er != nil {
		return nil, er
	}
	result := make(map[string]interface{}, len(body))
	var tagValue, primitiveValue, tagJSONValue string
	myDataReflect := reflect.Indirect(reflect.ValueOf(form))

	for i := 0; i < myDataReflect.NumField(); i++ {
		typeField := myDataReflect.Type().Field(i)
		tag := typeField.Tag
		tagValue = tag.Get("bson")
		tagJSONValue = tag.Get("json")
		primitiveValue = tag.Get("primitive")
		if val, ok := body[tagJSONValue]; ok {
			// fmt.Println(tagValue, tagJSONValue, reflect.TypeOf(val))
			switch myDataReflect.Field(i).Kind() {
			case reflect.String:
				result[tagValue] = val.(string)

			case reflect.Bool:
				result[tagValue] = val.(bool)

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// s := val.(string)
				// i, err := strconv.ParseInt(s, 10, 64)
				// if err == nil {
				// 	result[tagValue] = i
				// 	continue
				// }
				// f, err := strconv.ParseFloat(s, 64)
				// if err == nil {
				// 	result[tagValue] = f
				// 	continue
				// }
				result[tagValue] = val
			default:
				if primitiveValue == "true" {
					if reflect.ValueOf(val).Kind() == reflect.Slice {
						l := len(val.([]interface{}))
						idsPrimititiveSlice := make([]primitive.ObjectID, l)
						allValue := val.([]interface{})
						for i := range val.([]interface{}) {
							id, err := primitive.ObjectIDFromHex(allValue[i].(string))
							if err != nil {
								// todo error
								return result, err
							}
							idsPrimititiveSlice[i] = id
						}
						result[tagValue] = idsPrimititiveSlice
						// fmt.Println("default: ", tagValue, reflect.ValueOf(val).Kind())
					} else {

						id, err := primitive.ObjectIDFromHex(val.(string))
						if err != nil {
							// todo error
							return result, err
						}
						// fmt.Println("default: ", tagValue, reflect.ValueOf(val).Kind())
						result[tagValue] = id
					}
				} else {
					result[tagValue] = val
				}
				// value := myDataReflect.Field(i)
				// fmt.Println("   === default: tag=", tagValue, value)
				// fmt.Println("   === default: value=", value)
				// fmt.Println("   === default: tag primitiveValue=", primitiveValue)
				// fmt.Println("   === default: kind= ", myDataReflect.Field(i).Kind())
			}
		}
	}

	// fmt.Println("============result======================")
	// fmt.Println(result)
	// fmt.Println("==========================================")

	result["updated_at"] = time.Now()
	return result, nil
}

// BindAndValid binds and validates data.
func BindAndValidFromMarshal[V any](data interface{}, form V) (interface{}, error) {
	var body map[string]interface{}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Error Marshal patch data")
	}
	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, fmt.Errorf("Error Unmarshal patch data")
	}

	result := make(map[string]interface{}, len(body))
	var tagValue, primitiveValue, tagJSONValue string
	myDataReflect := reflect.Indirect(reflect.ValueOf(form))

	for i := 0; i < myDataReflect.NumField(); i++ {
		typeField := myDataReflect.Type().Field(i)
		tag := typeField.Tag
		tagValue = tag.Get("bson")
		tagJSONValue = tag.Get("json")
		primitiveValue = tag.Get("primitive")
		if val, ok := body[tagJSONValue]; ok {
			// fmt.Println(tagValue, tagJSONValue, reflect.TypeOf(val))
			switch myDataReflect.Field(i).Kind() {
			case reflect.String:
				result[tagValue] = val.(string)

			case reflect.Bool:
				result[tagValue] = val.(bool)

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// s := val.(string)
				// i, err := strconv.ParseInt(s, 10, 64)
				// if err == nil {
				// 	result[tagValue] = i
				// 	continue
				// }
				// f, err := strconv.ParseFloat(s, 64)
				// if err == nil {
				// 	result[tagValue] = f
				// 	continue
				// }
				result[tagValue] = val
			default:
				if primitiveValue == "true" {
					if reflect.ValueOf(val).Kind() == reflect.Slice {
						l := len(val.([]interface{}))
						idsPrimititiveSlice := make([]primitive.ObjectID, l)
						allValue := val.([]interface{})
						for i := range val.([]interface{}) {
							id, err := primitive.ObjectIDFromHex(allValue[i].(string))
							if err != nil {
								// todo error
								return result, err
							}
							idsPrimititiveSlice[i] = id
						}
						result[tagValue] = idsPrimititiveSlice
						// fmt.Println("default: ", tagValue, reflect.ValueOf(val).Kind())
					} else {

						id, err := primitive.ObjectIDFromHex(val.(string))
						if err != nil {
							// todo error
							return result, err
						}
						// fmt.Println("default: ", tagValue, reflect.ValueOf(val).Kind())
						result[tagValue] = id
					}
				} else {
					result[tagValue] = val
				}
				// value := myDataReflect.Field(i)
				// fmt.Println("   === default: tag=", tagValue, value)
				// fmt.Println("   === default: value=", value)
				// fmt.Println("   === default: tag primitiveValue=", primitiveValue)
				// fmt.Println("   === default: kind= ", myDataReflect.Field(i).Kind())
			}
		}
	}

	// fmt.Println("============result======================")
	// fmt.Println(result)
	// fmt.Println("==========================================")

	result["updated_at"] = time.Now()
	return result, nil
}

type StackNode struct {
	Node   interface{}
	Parent string
	Global bool
}

func BindJSON2[V any](raw map[string]json.RawMessage) (V, error) {
	var result V
	parsedData := make(map[string]map[string]string)
	// var raw map[string]json.RawMessage
	// err := json.Unmarshal([]byte(data), &raw)
	// if err != nil {
	// 	panic(err)
	// }

	parsed := make(map[string]interface{}, len(raw))

	for i := range raw {
		key := strings.Split(i, "__i18n__")
		if len(key) == 2 {
			if parsedData[key[1]] == nil {
				parsedData[key[1]] = map[string]string{}
			}
			// if reflect.ValueOf(raw[i]).Kind() != reflect.String {
			// 	return result, fmt.Errorf("field %s[%s] must be string [%v]", i, reflect.ValueOf(raw[i]).Kind(), raw[i])
			// }

			st, err := strconv.Unquote(string(raw[i]))
			if err != nil {
				return result, err
			}
			if st != "" {
				parsedData[key[1]][key[0]] = st
			}
		}
	}
	parsed["locale"] = parsedData

	for key, val := range raw {
		s := string(val)
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			parsed[key] = i
			continue
		}
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			parsed[key] = f
			continue
		}
		var v interface{}
		err = json.Unmarshal(val, &v)
		if err == nil {
			parsed[key] = v
			continue
		}
		parsed[key] = val
	}

	// Convert the map to JSON
	jsonData, err := json.Marshal(parsed)
	if err != nil {
		return result, err
	}
	// Convert the JSON to a struct
	json.Unmarshal(jsonData, &result)

	// for key, val := range parsed {
	// 	fmt.Printf("%T: %v %v \n", val, key, val)
	// }
	return result, nil
}

func BindJSON[V any](data map[string]interface{}) (V, error) {
	var result V
	parsedData := make(map[string]map[string]string)

	// fmt.Println(reflect.TypeOf(result))
	// fmt.Println(reflect.ValueOf(result))
	// fmt.Println(reflect.ValueOf(result).Kind())

	// fmt.Println(reflect.TypeOf(data))
	// fmt.Println(reflect.ValueOf(data))
	// fmt.Println(reflect.ValueOf(data).Kind())

	for i := range data {
		key := strings.Split(i, "__i18n__")
		if len(key) == 2 {
			if parsedData[key[1]] == nil {
				parsedData[key[1]] = map[string]string{}
			}
			if reflect.ValueOf(data[i]).Kind() != reflect.String {
				return result, fmt.Errorf("field %s must be string", i)
			}
			parsedData[key[1]][key[0]] = data[i].(string)
		}
	}
	// for k, v := range parsedData {
	data["locale"] = parsedData
	// fmt.Println("locale=>", data["locale"])
	// }

	// fmt.Println("==========data======================")
	// fmt.Println(data)
	// fmt.Println("==========/data======================")

	// elementsStructure := reflect.ValueOf(result)
	var tagValue, primitiveValue string
	structValue := reflect.ValueOf(&result).Elem()
	for i := 0; i < structValue.NumField(); i++ {
		typeField := structValue.Type().Field(i)
		tag := typeField.Tag
		tagValue = tag.Get("json")
		primitiveValue = tag.Get("primitive")
		// structValue2 := structValue.FieldByName(tagValue)
		// fmt.Println(tagValue, "= structValue2= ", structValue2)
		if !structValue.Field(i).CanSet() {
			return result, fmt.Errorf("no canset field %s", tagValue)
			// fmt.Println("===== nocanset", tagValue)
			// fmt.Println(i, tagValue, structValue.Field(i), structValue.Field(i).Type())
			// continue
		}
		val, ok := data[tagValue]
		if ok {

			if reflect.TypeOf(val) != structValue.Field(i).Type() {
				return result, fmt.Errorf("field %s(%s) must be type %s", tagValue, reflect.TypeOf(val), structValue.Field(i).Type().String())
				// fmt.Println("===== no compare types", tagValue)
				// fmt.Println(i, structValue.Field(i), structValue.Field(i).Type())
				// fmt.Println(i, val, reflect.TypeOf(val))
				// continue
			}

			// fmt.Println("====ok============================")
			// fmt.Println(i, tagValue, structValue.Field(i), structValue.Field(i).Type())
			valStructure := reflect.ValueOf(val)
			// if tagValue == "locale" {
			if primitiveValue == "true" {
				_, err := primitive.ObjectIDFromHex(val.(string))
				if err != nil {
					return result, err
				}
				// structValue.Field(i).Set(id)
				structValue.Field(i).Set(valStructure)
			} else {
				structValue.Field(i).Set(valStructure)
			}
			// }
		}
	}

	// fmt.Println("===result====")
	// fmt.Printf("%v", result)
	// fmt.Println("=============")

	return result, nil
}
