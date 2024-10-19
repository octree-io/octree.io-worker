package facade

import (
	"log"

	"octree.io-worker/internal/helpers"
	"octree.io-worker/internal/utils"
)

func JudgeTestCases(
	outputs []map[string]interface{},
	stdout string,
	answerAnyOrder bool,
	deepSort bool,
	returnType string,
) bool {
	parts := utils.SplitStringIntoParts(stdout, "\n")

	if len(outputs) != len(parts) {
		log.Println("Outputs and parts are different lengths")
		return false
	}

	for i := 0; i < len(parts); i++ {
		outputJsonString, err := utils.ConvertToJSONString(outputs[i]["output"])
		if err != nil {
			log.Println("Failed to convert output to JSON string")
			return false
		}

		result, err := helpers.CompareTestCaseOutputs(outputJsonString, parts[i], answerAnyOrder, deepSort, returnType)
		if err != nil {
			log.Printf("Test case %d failed: %v\n", i, err)
			return false
		}

		if !result {
			log.Printf("Test case %d failed. Expected %s but got %s\n", i, outputJsonString, parts[i])
			return false
		}
	}

	return true
}
