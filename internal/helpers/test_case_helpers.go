package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

func CompareTestCaseOutputs(expected string, actual string, answerAnyOrder bool, deepSort bool, returnType string) (bool, error) {
	switch returnType {
	case "int", "bool", "string", "TreeNode-int":
		return expected == actual, nil
	case "string[]":
		return compareStringArray(expected, actual, answerAnyOrder)
	case "int[]":
		return compareIntArray(expected, actual, answerAnyOrder)
	case "string[][]":
		return compareNestedStringArray(expected, actual, answerAnyOrder, deepSort)
	case "int[][]":
		return compareNestedIntArray(expected, actual, answerAnyOrder, deepSort)
	default:
		return false, fmt.Errorf("unsupported return type: %s", returnType)
	}
}

func compareStringArray(expected string, actual string, answerAnyOrder bool) (bool, error) {
	if expected == "null" {
		expected = "[]"
	}

	var expectedArray, actualArray []string
	err := json.Unmarshal([]byte(expected), &expectedArray)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal expected string[]: %v", err)
	}

	re := regexp.MustCompile(`'([^']*)'`)
	actual = re.ReplaceAllString(actual, `"$1"`)

	err = json.Unmarshal([]byte(actual), &actualArray)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal actual string[]: %v", err)
	}

	if answerAnyOrder {
		sort.Strings(expectedArray)
		sort.Strings(actualArray)
	}

	return reflect.DeepEqual(expectedArray, actualArray), nil
}

func compareIntArray(expected string, actual string, answerAnyOrder bool) (bool, error) {
	var expectedArray, actualArray []int
	err := json.Unmarshal([]byte(expected), &expectedArray)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal expected int[]: %v", err)
	}

	err = json.Unmarshal([]byte(actual), &actualArray)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal actual int[]: %v", err)
	}

	if answerAnyOrder {
		sort.Ints(expectedArray)
		sort.Ints(actualArray)
	}

	return reflect.DeepEqual(expectedArray, actualArray), nil
}

func compareNestedStringArray(expected string, actual string, answerAnyOrder bool, deepSort bool) (bool, error) {
	var expectedArray, actualArray [][]string
	err := json.Unmarshal([]byte(expected), &expectedArray)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal expected string[][]: %v", err)
	}

	re := regexp.MustCompile(`'([^']*)'`)
	actual = re.ReplaceAllString(actual, `"$1"`)

	err = json.Unmarshal([]byte(actual), &actualArray)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal actual string[][]: %v", err)
	}

	// Sort the inner lists
	if deepSort {
		for i := range expectedArray {
			sort.Strings(expectedArray[i])
		}
		for i := range actualArray {
			sort.Strings(actualArray[i])
		}
	}

	// Sort the outer lists
	if answerAnyOrder {
		sort.Slice(expectedArray, func(i, j int) bool {
			return strings.Join(expectedArray[i], ",") < strings.Join(expectedArray[j], ",")
		})
		sort.Slice(actualArray, func(i, j int) bool {
			return strings.Join(actualArray[i], ",") < strings.Join(actualArray[j], ",")
		})
	}

	return reflect.DeepEqual(expectedArray, actualArray), nil
}

func compareNestedIntArray(expected string, actual string, answerAnyOrder bool, deepSort bool) (bool, error) {
	var expectedArray, actualArray [][]int
	err := json.Unmarshal([]byte(expected), &expectedArray)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal expected int[][]: %v", err)
	}

	err = json.Unmarshal([]byte(actual), &actualArray)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal actual int[][]: %v", err)
	}

	if deepSort {
		for i := range expectedArray {
			sort.Ints(expectedArray[i])
		}
		for i := range actualArray {
			sort.Ints(actualArray[i])
		}
	}

	if answerAnyOrder {
		sort.Slice(expectedArray, func(i, j int) bool {
			return fmt.Sprint(expectedArray[i]) < fmt.Sprint(expectedArray[j])
		})
		sort.Slice(actualArray, func(i, j int) bool {
			return fmt.Sprint(actualArray[i]) < fmt.Sprint(actualArray[j])
		})
	}

	return reflect.DeepEqual(expectedArray, actualArray), nil
}
