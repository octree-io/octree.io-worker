package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertBsonMToStringMap(bsonMap bson.M) map[string]string {
	stringMap := make(map[string]string)

	for key, value := range bsonMap {
		strValue, ok := value.(string)
		if ok {
			stringMap[key] = strValue
		} else {
			fmt.Printf("Key %s has a non-string value: %v\n", key, value)
		}
	}

	return stringMap
}

func ConvertBsonToNative(bsonVal interface{}) interface{} {
	switch v := bsonVal.(type) {
	case primitive.A:
		var result []interface{}
		for _, elem := range v {
			result = append(result, ConvertBsonToNative(elem))
		}
		return result
	case primitive.M:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = ConvertBsonToNative(val)
		}
		return result
	default:
		return v
	}
}

func ConvertBSONValue(value interface{}) interface{} {
	switch v := value.(type) {
	case primitive.A:
		return []interface{}(v)
	case primitive.D:
		return bsonDToMap(v)
	case primitive.M:
		return map[string]interface{}(v)
	case primitive.ObjectID:
		return v.Hex()
	default:
		return value
	}
}

func bsonDToMap(d primitive.D) map[string]interface{} {
	result := make(map[string]interface{})
	for _, elem := range d {
		result[elem.Key] = elem.Value
	}
	return result
}
