package testharness

import (
	"encoding/json"
	"fmt"
	"log"
)

func JavaScriptHarness(code string, args map[string]string, testCases []map[string]interface{}, returnType string) string {
	jsArgs, err := convertJsArgToJson(args)
	if err != nil {
		log.Fatal("Error converting JSON to JavaScript")
	}

	jsTestCases, err := convertJsArgToJson(testCases)
	if err != nil {
		log.Fatal("Error converting JSON to JavaScript")
	}

	javaScriptCode := fmt.Sprintf(`class ListNode {
    constructor(val = 0, next = null) {
        this.val = val;
        this.next = next;
    }
}

class TreeNode {
    constructor(val = 0, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

class GraphNode {
    constructor(val = 0, neighbors = []) {
        this.val = val;
        this.neighbors = neighbors;
    }
}

// Code
%s

function listToTree(lst) {
    if (!lst.length) return null;

    const root = new TreeNode(lst[0]);
    const queue = [root];
    let index = 1;

    while (queue.length && index < lst.length) {
        const node = queue.shift();

        if (index < lst.length && lst[index] !== null) {
            node.left = new TreeNode(lst[index]);
            queue.push(node.left);
        }
        index++;

        if (index < lst.length && lst[index] !== null) {
            node.right = new TreeNode(lst[index]);
            queue.push(node.right);
        }
        index++;
    }

    return root;
}

function treeToList(root) {
    if (!root) return [];
    const result = [];
    const queue = [root];

    while (queue.length) {
        const node = queue.shift();
        if (node) {
            result.push(node.val);
            queue.push(node.left);
            queue.push(node.right);
        } else {
            result.push(null);
        }
    }

    while (result.length && result[result.length - 1] === null) {
        result.pop();
    }

    return result;
}

function findNodeByValue(root, value) {
    if (!root) return null;
    if (root.val === value) return root;

    const leftResult = findNodeByValue(root.left, value);
    if (leftResult) return leftResult;

    return findNodeByValue(root.right, value);
}

function linkedListToList(ll) {
    const lst = [];
    let current = ll;
    while (current) {
        lst.push(current.val);
        current = current.next;
    }
    return lst;
}

function listToLinkedList(lst) {
    const dummy = new ListNode(0);
    let tail = dummy;
    for (const val of lst) {
        tail.next = new ListNode(val);
        tail = tail.next;
    }
    return dummy.next;
}

function runTestCases() {
    const jsArgs = %s;
    const testCases = %s;
    const returnType = "%s";

    for (let i = 0; i < testCases.length; i++) {
        const test_case = testCases[i];
        const argNames = Object.keys(jsArgs);
        const methodArgs = [];
        let root = null;

        if (test_case.hasOwnProperty("root")) {
            root = listToTree(test_case["root"]);
            methodArgs.push(root);
        }

        for (const arg of argNames) {
            const argType = jsArgs[arg];
            let value = test_case[arg];

            if (argType === 'TreeNode') {
                if (Array.isArray(value)) {
                    value = listToTree(value);
                } else if (typeof value === 'number') {
                    value = findNodeByValue(root, value);
                }
            } else if (argType === "ListNode") {
                if (Array.isArray(value)) {
                    value = listToLinkedList(value);
                }
            }

            if (arg !== "root") {
                methodArgs.push(value);
            }
        }

        const result = solve(...methodArgs);

        if (result instanceof TreeNode) {
            if (returnType === "TreeNode-int") {
                console.log(result.val);
            } else {
                console.log(JSON.stringify(treeToList(result)));
            }
        } else if (result instanceof ListNode) {
            console.log(JSON.stringify(linkedListToList(result)));
        } else {
            if (!result && (returnType === "ListNode" || returnType === "TreeNode")) {
                console.log([]);
            } else {
                console.log(JSON.stringify(result));
            }
        }
    }
}

runTestCases();
`, code, jsArgs, jsTestCases, returnType)

	return javaScriptCode
}

func convertJsArgToJson(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error converting object to JSON: %w", err)
	}

	return string(jsonData), nil
}
