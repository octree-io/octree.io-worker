package testharness

import (
	"fmt"
	"log"

	utils "octree.io-worker/internal/utils/converters"
)

func PythonHarness(code string, args map[string]string, testCases []map[string]interface{}, returnType string) string {
	pyArgs, err := utils.JsonToPython(args)
	if err != nil {
		log.Fatal("Error converting JSON to Python")
	}

	pyTestCases, err := utils.JsonToPython(testCases)
	if err != nil {
		log.Fatal("Error converting JSON to Python")
	}

	pythonCode := fmt.Sprintf(`from collections import *
from array import *
from bisect import *
from typing import *
import collections
import array
import bisect
import heapq

class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next

class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right

class GraphNode:
    def __init__(self, val=0, neighbors=None):
        self.val = val
        self.neighbors = neighbors if neighbors is not None else []

%s

def list_to_tree(lst):
    if not lst:
        return None

    root = TreeNode(lst[0])
    queue = collections.deque([root])
    index = 1

    while queue and index < len(lst):
        node = queue.popleft()

        if index < len(lst) and lst[index] is not None:
            node.left = TreeNode(lst[index])
            queue.append(node.left)
        else:
            node.left = None
        index += 1

        if index < len(lst) and lst[index] is not None:
            node.right = TreeNode(lst[index])
            queue.append(node.right)
        else:
            node.right = None
        index += 1

    return root

def tree_to_list(root):
    if not root:
        return []
    result = []
    queue = collections.deque([root])
    while queue:
        node = queue.popleft()
        if node:
            result.append(node.val)
            queue.append(node.left)
            queue.append(node.right)
        else:
            result.append(None)

    while result and result[-1] is None:
        result.pop()
    return result

def find_node_by_value(root: TreeNode, value: int) -> TreeNode:
    if not root:
        return None
    if root.val == value:
        return root
    left_result = find_node_by_value(root.left, value)
    if left_result:
        return left_result
    return find_node_by_value(root.right, value)

def linked_list_to_list(ll):
    lst = []
    current = ll
    while current:
        lst.append(current.val)
        current = current.next
    return lst

def list_to_linked_list(lst):
    dummy = ListNode(0)
    tail = dummy
    for val in lst:
        tail.next = ListNode(val)
        tail = tail.next
    return dummy.next

def run_test_cases():
    solution = Solution()
    py_args = %s
    test_cases = %s
    return_type = "%s"

    for i, test_case in enumerate(test_cases):
        arg_names = list(py_args.keys())
        method_args = []
        root = None

        if "root" in test_case:
            root = list_to_tree(test_case["root"])
            method_args.append(root)

        for arg in arg_names:
            arg_type = py_args[arg]
            value = test_case.get(arg, None)

            if arg_type == 'TreeNode':
                if isinstance(value, list):
                    value = list_to_tree(value)
                elif isinstance(value, int):
                    value = find_node_by_value(root, value)
            elif arg_type == "ListNode":
                if isinstance(value, list):
                    value = list_to_linked_list(value)

            if arg != "root":
                method_args.append(value)

        result = solution.solve(*method_args)

        if isinstance(result, TreeNode):
            if return_type == "TreeNode-int":
                print(result.val)
            else:
                print(tree_to_list(result))
        elif isinstance(result, ListNode):
            print(linked_list_to_list(result))
        else:
            if not result and return_type == "ListNode" or return_type == "TreeNode":
                print([])
            else:
                print(result)

run_test_cases()
`, code, pyArgs, pyTestCases, returnType)

	return pythonCode
}
