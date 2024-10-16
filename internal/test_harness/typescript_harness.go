package testharness

import (
	"encoding/json"
	"fmt"
	"log"
)

func TypeScriptHarness(code string, args map[string]string, testCases []map[string]interface{}, returnType string) string {
	tsArgs, err := convertTsArgToJson(args)
	if err != nil {
		log.Fatal("Error converting JSON to TypeScript")
	}

	tsTestCases, err := convertTsArgToJson(testCases)
	if err != nil {
		log.Fatal("Error converting JSON to TypeScript")
	}

	typeScriptCode := fmt.Sprintf(`class ListNode {
    val: number;
    next: ListNode | null;

    constructor(val: number = 0, next: ListNode | null = null) {
        this.val = val;
        this.next = next;
    }
}

class TreeNode {
    val: number | null;
    left: TreeNode | null;
    right: TreeNode | null;

    constructor(val: number | null = 0, left: TreeNode | null = null, right: TreeNode | null = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

class GraphNode {
    val: number;
    neighbors: GraphNode[];

    constructor(val: number = 0, neighbors: GraphNode[] = []) {
        this.val = val;
        this.neighbors = neighbors;
    }
}

// Code
%s

function listToTree(lst: (number | null)[]): TreeNode | null {
    if (!lst.length) return null;

    const root = new TreeNode(lst[0]);
    const queue: TreeNode[] = [root];
    let index = 1;

    while (queue.length && index < lst.length) {
        const node = queue.shift()!;

        if (index < lst.length && lst[index] !== null) {
            node.left = new TreeNode(lst[index]!);
            queue.push(node.left);
        }
        index++;

        if (index < lst.length && lst[index] !== null) {
            node.right = new TreeNode(lst[index]!);
            queue.push(node.right);
        }
        index++;
    }

    return root;
}

function treeToList(root: TreeNode | null): (number | null)[] {
    if (!root) return [];
    const result: (number | null)[] = [];
    const queue: (TreeNode | null)[] = [root];

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

function findNodeByValue(root: TreeNode | null, value: number): TreeNode | null {
    if (!root) return null;
    if (root.val === value) return root;

    const leftResult = findNodeByValue(root.left, value);
    if (leftResult) return leftResult;

    return findNodeByValue(root.right, value);
}

function linkedListToList(ll: ListNode | null): number[] {
    const lst: number[] = [];
    let current = ll;
    while (current) {
        lst.push(current.val);
        current = current.next;
    }
    return lst;
}

function listToLinkedList(lst: number[]): ListNode | null {
    const dummy = new ListNode(0);
    let tail = dummy;
    for (const val of lst) {
        tail.next = new ListNode(val);
        tail = tail.next;
    }
    return dummy.next;
}

function runTestCases() {
    const tsArgs: Record<string, string> = %s;
    const testCases: Record<string, any>[] = %s;
    const returnType: string = "%s";

    for (let i = 0; i < testCases.length; i++) {
        const test_case = testCases[i];
        const argNames = Object.keys(tsArgs);
        const methodArgs: any[] = [];
        let root: TreeNode | null = null;

        if (test_case.hasOwnProperty("root")) {
            root = listToTree(test_case["root"]);
            methodArgs.push(root);
        }

        for (const arg of argNames) {
            const argType = tsArgs[arg];
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

        const result: any = solve(...(methodArgs as [any]));

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
`, code, tsArgs, tsTestCases, returnType)

	return typeScriptCode
}

func convertTsArgToJson(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error converting object to JSON: %w", err)
	}

	return string(jsonData), nil
}
