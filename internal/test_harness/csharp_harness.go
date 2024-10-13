package testharness

import (
	"fmt"
	"strings"
)

func CsharpHarness(code string, args map[string]string, testCases []map[string]interface{}, returnType string) string {
	harnessCode := fmt.Sprintf(`using System;
using System.Collections.Generic;
using System.Linq;

public static class Globals {
    public static string returnType = "%s";
}

public class ListNode {
    public int val;
    public ListNode next;

    public ListNode() {
        this.val = 0;
        this.next = null;
    }

    public ListNode(int val) {
        this.val = val;
        this.next = null;
    }

    public ListNode(int val, ListNode next) {
        this.val = val;
        this.next = next;
    }
}

public class TreeNode {
    public int val;
    public TreeNode left;
    public TreeNode right;

    public TreeNode() {
        this.val = 0;
        this.left = null;
        this.right = null;
    }

    public TreeNode(int val) {
        this.val = val;
        this.left = null;
        this.right = null;
    }

    public TreeNode(int val, TreeNode left, TreeNode right) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

public class GraphNode {
    public int val;
    public List<GraphNode> neighbors;

    public GraphNode() {
        this.val = 0;
        this.neighbors = new List<GraphNode>();
    }

    public GraphNode(int val) {
        this.val = val;
        this.neighbors = new List<GraphNode>();
    }

    public GraphNode(int val, List<GraphNode> neighbors) {
        this.val = val;
        this.neighbors = neighbors;
    }
}

// Code
%s

public static class DSAHelpers
{
    public static TreeNode ListToTree(List<int?> lst)
    {
        if (lst.Count == 0 || lst[0] == null) return null;

        TreeNode root = new TreeNode(lst[0].Value);
        Queue<TreeNode> queue = new Queue<TreeNode>();
        queue.Enqueue(root);
        int index = 1;

        while (queue.Count > 0 && index < lst.Count)
        {
            TreeNode node = queue.Dequeue();

            if (index < lst.Count && lst[index] != null)
            {
                node.left = new TreeNode(lst[index].Value);
                queue.Enqueue(node.left);
            }
            index++;

            if (index < lst.Count && lst[index] != null)
            {
                node.right = new TreeNode(lst[index].Value);
                queue.Enqueue(node.right);
            }
            index++;
        }

        return root;
    }

    public static List<int?> TreeToList(TreeNode root)
    {
        List<int?> result = new List<int?>();
        if (root == null) return result;

        Queue<TreeNode> queue = new Queue<TreeNode>();
        queue.Enqueue(root);

        while (queue.Count > 0)
        {
            TreeNode node = queue.Dequeue();

            if (node != null)
            {
                result.Add(node.val);
                queue.Enqueue(node.left);
                queue.Enqueue(node.right);
            }
            else
            {
                result.Add(null);
            }
        }

        while (result.Count > 0 && result[result.Count - 1] == null)
        {
            result.RemoveAt(result.Count - 1);
        }

        return result;
    }

    public static TreeNode FindNodeByValue(TreeNode root, int value)
    {
        if (root == null) return null;
        if (root.val == value) return root;

        TreeNode leftResult = FindNodeByValue(root.left, value);
        if (leftResult != null) return leftResult;

        return FindNodeByValue(root.right, value);
    }

    public static ListNode ListToLinkedList(List<int> lst)
    {
        if (lst.Count == 0) return null;

        ListNode dummy = new ListNode(0);
        ListNode tail = dummy;

        foreach (int val in lst)
        {
            tail.next = new ListNode(val);
            tail = tail.next;
        }

        return dummy.next;
    }
}

public static class TestHelper
{
    public static void PrintResult(object result)
    {
        if (result is int[])
        {
            Console.WriteLine(string.Join(", ", (int[])result));
        }
        else if (result is List<int>)
        {
            Console.WriteLine(string.Join(", ", (List<int>)result));
        }
        else if (result is int[][])
        {
            foreach (var array in (int[][])result)
            {
                Console.WriteLine("[" + string.Join(", ", array) + "]");
            }
        }
        else if (result is List<List<int>>)
        {
            foreach (var list in (List<List<int>>)result)
            {
                Console.WriteLine("[" + string.Join(", ", list) + "]");
            }
        }
        else if (result is List<string>)
        {
            Console.WriteLine(string.Join(", ", (List<string>)result));
        }
        else if (result is string[][])
        {
            foreach (var array in (string[][])result)
            {
                Console.WriteLine("[" + string.Join(", ", array) + "]");
            }
        }
        else if (result is List<List<string>>)
        {
            foreach (var list in (List<List<string>>)result)
            {
                Console.WriteLine("[" + string.Join(", ", list) + "]");
            }
        }
        else if (result is ListNode)
        {
            PrintListNode((ListNode)result);
        }
        else if (result is TreeNode)
        {
            PrintTreeNode((TreeNode)result);
        }
        else
        {
            Console.WriteLine(result);
        }
    }

    private static void PrintTreeNode(TreeNode root)
    {
        if (Globals.returnType == "TreeNode-int") {
            Console.WriteLine(root.val.ToString());
            return;
        }

        var result = new List<string>();
        var queue = new Queue<TreeNode>();
        queue.Enqueue(root);

        while (queue.Count > 0)
        {
            var node = queue.Dequeue();
            if (node == null)
            {
                result.Add("null");
            }
            else
            {
                result.Add(node.val.ToString());
                queue.Enqueue(node.left);
                queue.Enqueue(node.right);
            }
        }

        while (result.Count > 0 && result[result.Count - 1] == "null")
        {
            result.RemoveAt(result.Count - 1);
        }

        Console.WriteLine("[" + string.Join(", ", result) + "]");
    }

    private static void PrintListNode(ListNode head)
    {
        var result = new List<int>();
        var current = head;

        while (current != null)
        {
            result.Add(current.val);
            current = current.next;
        }

        Console.WriteLine("[" + string.Join(", ", result) + "]");
    }
}

public class TestHarness
{
    public static void Main(string[] args)
    {
        Solution solution = new Solution();

        var testCases = new List<Dictionary<string, object>>();

        // Test cases
        %s

        // CSharp args
        %s

        string[] argNames = new string[] { %s };

        foreach (var testCase in testCases)
        {
            object[] methodArgs = new object[argNames.Length];

            TreeNode root = null;
            if (testCase.ContainsKey("root"))
            {
                root = (TreeNode)testCase["root"];
            }

            for (int j = 0; j < argNames.Length; j++)
            {
                string argName = argNames[j];
                string argType = csharpArgs[argName];
                object value = testCase[argName];

                if (argType == "TreeNode" && value is int)
                {
                    int nodeValue = (int)value;
                    value = DSAHelpers.FindNodeByValue(root, nodeValue);
                }

                methodArgs[j] = value;
            }

            var result = solution.Solve(%s);

            TestHelper.PrintResult(result);
        }
    }
}

`, returnType, code, generateCsharpTestCases(testCases, args), generateCsharpArgs(args), generateCsharpArgNames(args), generateCsharpParameters(args))

	fmt.Println(harnessCode)

	return harnessCode
}

func generateCsharpTestCases(testCases []map[string]interface{}, args map[string]string) string {
	var result strings.Builder

	for index, testCase := range testCases {
		result.WriteString(fmt.Sprintf("\nvar testCase%d = new Dictionary<string, object>();\n", index))

		for arg, argType := range args {
			value := testCase[arg]

			switch argType {
			case "int[]":
				intArray := value.([]int)
				intArrayValues := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(intArray)), ", "), "[]")
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = new int[]{%s};\n", index, arg, intArrayValues))

			case "int[][]":
				int2DArray := value.([][]int)
				int2DArrayValues := ""
				for _, arr := range int2DArray {
					int2DArrayValues += fmt.Sprintf("new int[]{%s}, ", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]"))
				}
				int2DArrayValues = strings.TrimSuffix(int2DArrayValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = new int[][]{%s};\n", index, arg, int2DArrayValues))

			case "List<int[]>":
				listOfIntArray := value.([][]int)
				listValues := ""
				for _, arr := range listOfIntArray {
					listValues += fmt.Sprintf("new int[]{%s}, ", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]"))
				}
				listValues = strings.TrimSuffix(listValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = new List<int[]> { %s };\n", index, arg, listValues))

			case "List<List<int>>":
				listOfLists := value.([][]int)
				listValues := ""
				for _, list := range listOfLists {
					listValues += fmt.Sprintf("new List<int> { %s }, ", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), ", "), "[]"))
				}
				listValues = strings.TrimSuffix(listValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = new List<List<int>> { %s };\n", index, arg, listValues))

			case "List<string[]>":
				listOfStringArray := value.([][]string)
				listValues := ""
				for _, arr := range listOfStringArray {
					listValues += fmt.Sprintf("new string[]{%s}, ", strings.Trim(strings.Join(arr, ", "), "[]"))
				}
				listValues = strings.TrimSuffix(listValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = new List<string[]> { %s };\n", index, arg, listValues))

			case "string[][]":
				string2DArray := value.([][]string)
				string2DArrayValues := ""
				for _, arr := range string2DArray {
					string2DArrayValues += fmt.Sprintf("new string[]{%s}, ", strings.Join(arr, ", "))
				}
				string2DArrayValues = strings.TrimSuffix(string2DArrayValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = new string[][]{%s};\n", index, arg, string2DArrayValues))

			case "List<List<string>>":
				listOfStringLists := value.([][]string)
				listValues := ""
				for _, list := range listOfStringLists {
					listValues += fmt.Sprintf("new List<string> { %s }, ", strings.Join(list, ", "))
				}
				listValues = strings.TrimSuffix(listValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = new List<List<string>> { %s };\n", index, arg, listValues))

			case "int":
				intValue := value.(int)
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = %d;\n", index, arg, intValue))

			case "string":
				stringValue := value.(string)
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = \"%s\";\n", index, arg, stringValue))

			case "TreeNode":
				switch value.(type) {
				case int:
					result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = %d;\n", index, arg, value))

				default:
					treeNodeList := value.([]interface{})
					treeNodeListValues := ""
					for _, val := range treeNodeList {
						if val != nil {
							treeNodeListValues += fmt.Sprintf("%d, ", val)
						} else {
							treeNodeListValues += "null, "
						}
					}
					treeNodeListValues = strings.TrimSuffix(treeNodeListValues, ", ")
					result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = DSAHelpers.ListToTree(new List<int?> { %s });\n", index, arg, treeNodeListValues))
				}

			case "ListNode":
				listNodeValues := value.([]interface{})
				listNodeValuesStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(listNodeValues)), ", "), "[]")
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = DSAHelpers.ListToLinkedList(new List<int> { %s });\n", index, arg, listNodeValuesStr))

			default:
				result.WriteString(fmt.Sprintf("testCase%d[\"%s\"] = %v;\n", index, arg, value))
			}
		}

		result.WriteString(fmt.Sprintf("\ntestCases.Add(testCase%d);\n", index))
	}

	return result.String()
}

func generateCsharpArgNames(args map[string]string) string {
	keys := make([]string, 0, len(args))
	for key := range args {
		keys = append(keys, fmt.Sprintf("\"%s\"", key))
	}

	result := strings.Join(keys, ", ")
	return result
}

func generateCsharpParameters(args map[string]string) string {
	var result []string

	idx := 0
	for _, argType := range args {
		formatted := fmt.Sprintf("((%s) methodArgs[%d])", argType, idx)
		result = append(result, formatted)
		idx++
	}

	finalResult := strings.Join(result, ", ")
	return finalResult
}

func generateCsharpArgs(args map[string]string) string {
	var result []string
	result = append(result, "var csharpArgs = new Dictionary<string, string> {")
	for argName, argType := range args {
		result = append(result, fmt.Sprintf("{\"%s\", \"%s\"},", argName, argType))
	}
	result = append(result, "};")
	return strings.Join(result, "\n")
}
