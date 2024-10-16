package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
	"octree.io-worker/internal/facade"
	testharness "octree.io-worker/internal/test_harness"
	"octree.io-worker/internal/utils/converters"
)

type StdoutItem struct {
	Text string `json:"text"`
}

type BuildResult struct {
	Code               int      `json:"code"`
	TimedOut           bool     `json:"timedOut"`
	Stdout             []string `json:"stdout"`
	Stderr             []string `json:"stderr"`
	Downloads          []string `json:"downloads"`
	ExecutableFilename string   `json:"executableFilename"`
	CompilationOptions []string `json:"compilationOptions"`
}

type CompilerExplorerResponse struct {
	Code                       int          `json:"code"`
	OkToCache                  bool         `json:"okToCache"`
	TimedOut                   bool         `json:"timedOut"`
	Stdout                     []StdoutItem `json:"stdout"`
	Stderr                     []string     `json:"stderr"`
	Truncated                  bool         `json:"truncated"`
	ExecTime                   string       `json:"execTime"`
	ProcessExecutionResultTime float64      `json:"processExecutionResultTime"`
	DidExecute                 bool         `json:"didExecute"`
	BuildResult                BuildResult  `json:"buildResult"`
}

func processCompilationRequest() {
	start := time.Now()

	language := "typescript"

	code := `function solve(words: any) {
	  let anagramGroups: any = {};

	  for (let word of words) {
	    // Sort the word to use as the key for grouping anagrams
	    let sortedWord = word.split("").sort().join("");

	    // Insert the word into the corresponding anagram group
	    if (!anagramGroups[sortedWord]) {
	      anagramGroups[sortedWord] = [];
	    }
	    anagramGroups[sortedWord].push(word);
	  }

	  // Collect all the grouped anagrams into a result array
	  let result = [];
	  for (let key in anagramGroups) {
	    result.push(anagramGroups[key]);
	  }

	  return result;
	}`

	args := map[string]string{
		"words": "string[]",
	}

	testCases := []map[string]interface{}{
		{
			"words":  []interface{}{"eat", "tea", "tan", "ate", "nat", "bat"},
			"output": []interface{}{[]interface{}{"eat", "tea", "ate"}, []interface{}{"tan", "nat"}, []interface{}{"bat"}},
		},
		{
			"words":  []interface{}{"abc", "bca", "cab", "dog", "god", "xyz"},
			"output": []interface{}{[]interface{}{"abc", "bca", "cab"}, []interface{}{"dog", "god"}, []interface{}{"xyz"}},
		},
		{
			"words":  []interface{}{"apple", "pale", "peal", "leap"},
			"output": []interface{}{[]interface{}{"apple"}, []interface{}{"pale", "peal", "leap"}},
		},
	}

	returnType := "string[][]"

	// 	code := `function solve(root, targetSum) {
	//     if (!root) return false;

	//     if (!root.left && !root.right) {
	//         return targetSum === root.val;
	//     }

	//     const remainingSum = targetSum - root.val;
	//     return solve(root.left, remainingSum) || solve(root.right, remainingSum);
	// }`

	// 	args := map[string]string{
	// 		"root":      "TreeNode",
	// 		"targetSum": "int",
	// 	}

	// 	testCases := []map[string]interface{}{
	// 		{
	// 			"root":      []interface{}{5, 4, 8, 11, nil, 13, 4, 7, 2, nil, nil, nil, 1},
	// 			"targetSum": 22,
	// 			"output":    true,
	// 		},
	// 		{
	// 			"root":      []interface{}{1, 2, 3},
	// 			"targetSum": 5,
	// 			"output":    false,
	// 		},
	// 		{
	// 			"root":      []interface{}{1, 2},
	// 			"targetSum": 1,
	// 			"output":    false,
	// 		},
	// 	}

	// 	returnType := "bool"

	// 	code := `function solve(root: TreeNode | null, p: TreeNode, q: TreeNode): TreeNode | null {
	//     if (!root || root === p || root === q) return root;

	//     const left = solve(root.left, p, q);
	//     const right = solve(root.right, p, q);

	//     if (left && right) return root;

	//     return left ? left : right;
	// }`

	// 	args := map[string]string{
	// 		"root": "TreeNode",
	// 		"p":    "TreeNode",
	// 		"q":    "TreeNode",
	// 	}

	// 	testCases := []map[string]interface{}{
	// 		{
	// 			"root":   []interface{}{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, // Tree structure
	// 			"p":      5,                                                  // Node p
	// 			"q":      1,                                                  // Node q
	// 			"output": 3,                                                  // Expected LCA
	// 		},
	// 		{
	// 			"root":   []interface{}{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, // Tree structure
	// 			"p":      5,                                                  // Node p
	// 			"q":      4,                                                  // Node q
	// 			"output": 5,                                                  // Expected LCA
	// 		},
	// 		{
	// 			"root":   []interface{}{1, 2}, // Tree structure
	// 			"p":      1,                   // Node p
	// 			"q":      2,                   // Node q
	// 			"output": 1,                   // Expected LCA
	// 		},
	// 	}

	// 	returnType := "TreeNode-int"

	// 	code := `function solve(head: ListNode | null) {
	//     let prev = null;
	//     let current = head;

	//     while (current !== null) {
	//         let nextNode = current.next;
	//         current.next = prev;
	//         prev = current;
	//         current = nextNode;
	//     }

	//     return prev;
	// }`

	// 	// Input type: A linked list (represented as an array for test cases)
	// 	args := map[string]string{
	// 		"head": "ListNode", // The linked list to be reversed
	// 	}

	// 	// Test cases for the reverse linked list problem
	// 	testCases := []map[string]interface{}{
	// 		{
	// 			"head":   []interface{}{1, 2, 3, 4, 5}, // Initial linked list: 1 -> 2 -> 3 -> 4 -> 5
	// 			"output": []interface{}{5, 4, 3, 2, 1}, // Reversed linked list: 5 -> 4 -> 3 -> 2 -> 1
	// 		},
	// 		{
	// 			"head":   []interface{}{1, 2}, // Initial linked list: 1 -> 2
	// 			"output": []interface{}{2, 1}, // Reversed linked list: 2 -> 1
	// 		},
	// 		{
	// 			"head":   []interface{}{}, // Empty list
	// 			"output": []interface{}{}, // Still empty after reversal
	// 		},
	// 	}

	// 	returnType := "ListNode"

	var wrappedCode string

	switch language {
	case "python":
		wrappedCode = testharness.PythonHarness(code, args, testCases, returnType)

	case "cpp":
		wrappedCode = testharness.CppHarness(code, args, testCases, returnType)

	case "csharp":
		wrappedCode = testharness.CsharpHarness(code, args, testCases, returnType)

	case "java":
		wrappedCode = testharness.JavaHarness(code, args, testCases, returnType)

	case "ruby":
		wrappedCode = testharness.RubyHarness(code, args, testCases, returnType)

	case "javascript":
		wrappedCode = testharness.JavaScriptHarness(code, args, testCases, returnType)

	case "typescript":
		wrappedCode = testharness.TypeScriptHarness(code, args, testCases, returnType)

	default:
		fmt.Println("Unsupported language")
		return
	}

	output, err := facade.CompilerExplorer(language, wrappedCode)
	if err != nil {
		log.Printf("Error while executing compile: %v", err)
	}

	var jsonOutput CompilerExplorerResponse
	json.Unmarshal(([]byte)(output), &jsonOutput)

	elapsed := time.Since(start)

	log.Printf("Request took %s to execute", elapsed)

	if len(testCases) != len(jsonOutput.Stdout) {
		log.Println("Test cases and stdout lengths differ")
	} else {
		for i, item := range jsonOutput.Stdout {
			convertedJson, _ := converters.JsonToPython(testCases[i]["output"])
			log.Printf("Expected: %s, Actual: %s, Equivalence: %t",
				convertedJson, item.Text, convertedJson == item.Text)
		}
	}
}

func SpawnCompilationWorker(id int, msgs <-chan ampq.Delivery) {
	for msg := range msgs {
		log.Printf("[Compilation Worker %d] Received message: %s", id, msg.Body)

		processCompilationRequest()

		if err := msg.Ack(false); err != nil {
			log.Printf("[Compilation Worker %d] Failed to ack message: %v", id, err)
		} else {
			log.Printf("[Compilation Worker %d] Message ack'd", id)
		}
	}
}
