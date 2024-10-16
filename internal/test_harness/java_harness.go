package testharness

import (
	"fmt"
	"strings"

	"octree.io-worker/internal/utils"
)

func JavaHarness(code string, args map[string]string, testCases []map[string]interface{}, returnType string) string {
	harnessCode :=
		fmt.Sprintf(`import java.util.*;
import java.lang.*;
import java.io.*;
import java.util.stream.*;

class Globals {
    public static final String returnType = "%s";
}

class ListNode {
    int val;
    ListNode next;

    ListNode() {}

    ListNode(int val) {
        this.val = val;
    }

    ListNode(int val, ListNode next) {
        this.val = val;
        this.next = next;
    }
}

class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;

    TreeNode() {}

    TreeNode(int val) {
        this.val = val;
    }

    TreeNode(int val, TreeNode left, TreeNode right) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

class GraphNode {
    int val;
    List<GraphNode> neighbors;

    GraphNode() {
        neighbors = new ArrayList<>();
    }

    GraphNode(int val) {
        this.val = val;
        this.neighbors = new ArrayList<>();
    }

    GraphNode(int val, List<GraphNode> neighbors) {
        this.val = val;
        this.neighbors = neighbors;
    }
}

// Code
%s

class DSAHelpers {
    public static TreeNode listToTree(List<Integer> lst) {
        if (lst.size() == 0 || lst.get(0) == null) return null;

        TreeNode root = new TreeNode(lst.get(0));
        Queue<TreeNode> queue = new LinkedList<>();
        queue.add(root);
        int index = 1;

        while (!queue.isEmpty() && index < lst.size()) {
            TreeNode node = queue.poll();

            if (index < lst.size() && lst.get(index) != null) {
                node.left = new TreeNode(lst.get(index));
                queue.add(node.left);
            }
            index++;

            if (index < lst.size() && lst.get(index) != null) {
                node.right = new TreeNode(lst.get(index));
                queue.add(node.right);
            }
            index++;
        }

        return root;
    }

    public static TreeNode listToTree(Integer[] lst) {
        if (lst.length == 0 || lst[0] == null) return null;

        TreeNode root = new TreeNode(lst[0]);
        Queue<TreeNode> queue = new LinkedList<>();
        queue.add(root);
        int index = 1;

        while (!queue.isEmpty() && index < lst.length) {
            TreeNode node = queue.poll();

            if (index < lst.length && lst[index] != null) {
                node.left = new TreeNode(lst[index]);
                queue.add(node.left);
            }
            index++;

            if (index < lst.length && lst[index] != null) {
                node.right = new TreeNode(lst[index]);
                queue.add(node.right);
            }
            index++;
        }

        return root;
    }


    public static List<Integer> treeToList(TreeNode root) {
        List<Integer> result = new ArrayList<>();
        if (root == null) return result;

        Queue<TreeNode> queue = new LinkedList<>();
        queue.add(root);

        while (!queue.isEmpty()) {
            TreeNode node = queue.poll();

            if (node != null) {
                result.add(node.val);
                queue.add(node.left);
                queue.add(node.right);
            } else {
                result.add(null);
            }
        }

        // Remove trailing nulls
        while (!result.isEmpty() && result.get(result.size() - 1) == null) {
            result.remove(result.size() - 1);
        }

        return result;
    }

    public static TreeNode findNodeByValue(TreeNode root, int value) {
        if (root == null) return null;
        if (root.val == value) return root;

        TreeNode leftResult = findNodeByValue(root.left, value);
        if (leftResult != null) return leftResult;

        return findNodeByValue(root.right, value);
    }

    public static ListNode listToLinkedList(List<Integer> lst) {
        if (lst.isEmpty()) return null;

        ListNode dummy = new ListNode(0);
        ListNode tail = dummy;

        for (int val : lst) {
            tail.next = new ListNode(val);
            tail = tail.next;
        }

        return dummy.next;
    }

    public static ListNode listToLinkedList(Integer[] lst) {
        if (lst.length == 0) return null;

        ListNode dummy = new ListNode(0);
        ListNode tail = dummy;

        for (int val : lst) {
            tail.next = new ListNode(val);
            tail = tail.next;
        }

        return dummy.next;
    }

    public static ListNode listToLinkedList(int[] lst) {
        if (lst.length == 0) return null;

        ListNode dummy = new ListNode(0);
        ListNode tail = dummy;

        for (int val : lst) {
            tail.next = new ListNode(val);
            tail = tail.next;
        }

        return dummy.next;
    }
}

class TestHelper {
    public static void printResult(Object result) {
        if (result instanceof int[]) {
            System.out.println(Arrays.toString((int[]) result));
        } else if (result instanceof List) {
            System.out.println(listToString((List<?>) result));
        } else if (result instanceof Integer || result instanceof String || result instanceof Boolean) {
            System.out.println(result);
        } else if (result instanceof ListNode) {
            printListNode((ListNode) result);
        } else if (result instanceof TreeNode) {
            printTreeNode((TreeNode) result);
        } else {
            System.out.println(result);
        }
    }

    private static void printTreeNode(TreeNode root) {
        if (root == null) {
            System.out.println("null");
            return;
        }

        if (Globals.returnType.equals("TreeNode-int")) {
            System.out.println(root.val);
            return;
        }

        List<String> result = new ArrayList<>();
        Queue<TreeNode> queue = new LinkedList<>();
        queue.add(root);

        while (!queue.isEmpty()) {
            TreeNode node = queue.poll();
            if (node == null) {
                result.add("null");
            } else {
                result.add(String.valueOf(node.val));
                queue.add(node.left);
                queue.add(node.right);
            }
        }

        // Remove trailing "null" values
        while (result.size() > 0 && result.get(result.size() - 1).equals("null")) {
            result.remove(result.size() - 1);
        }

        System.out.println("[" + String.join(", ", result) + "]");
    }

    private static void printListNode(ListNode head) {
        List<Integer> result = new ArrayList<>();
        ListNode current = head;

        while (current != null) {
            result.add(current.val);
            current = current.next;
        }

        System.out.println("[" + result.stream().map(String::valueOf).collect(Collectors.joining(", ")) + "]");
    }

    private static String listToString(List<?> list) {
        StringBuilder sb = new StringBuilder("[");
        for (int i = 0; i < list.size(); i++) {
            Object item = list.get(i);

            if (item == null) {
                sb.append("null");  // Handle null values
            } else if (item instanceof List) {
                sb.append(listToString((List<?>) item));
            } else if (item instanceof String) {
                sb.append("\"").append(item.toString()).append("\"");  // Add quotes around strings
            } else {
                sb.append(item.toString());
            }

            if (i != list.size() - 1) {
                sb.append(", ");
            }
        }
        sb.append("]");
        return sb.toString();
    }
}

class TestHarness {
    public static void main(String[] args) {
        Solution solution = new Solution();

        List<Map<String, Object>> testCases = new ArrayList<>();

        // Test cases
        %s

        // Java args
        %s

        String[] argNames = new String[] { %s };

        for (int i = 0; i < testCases.size(); i++) {
            Map<String, Object> testCase = testCases.get(i);
            Object[] methodArgs = new Object[argNames.length];

            TreeNode root = null;
            if (testCase.containsKey("root")) {
                root = (TreeNode) testCase.get("root");
            }

            for (int j = 0; j < argNames.length; j++) {
                String argName = argNames[j];
                String argType = javaArgs.get(argName);
                Object value = testCase.get(argName);

                if ("TreeNode".equals(argType) && value instanceof Integer) {
                    int nodeValue = (Integer) value;
                    value = DSAHelpers.findNodeByValue(root, nodeValue);
                }

                methodArgs[j] = value;
            }

            Object result = solution.solve(%s);

            TestHelper.printResult(result);
        }
    }
}
  `, returnType, code, generateJavaTestCases(testCases, args), generateJavaArgs(args), generateJavaArgNames(args), generateJavaMethodArgs(args))

	return harnessCode
}

func generateJavaTestCases(testCases []map[string]interface{}, args map[string]string) string {
	var result strings.Builder

	for index, testCase := range testCases {
		result.WriteString(fmt.Sprintf("\nMap<String, Object> testCase%d = new HashMap<>();\n", index))

		for arg, argType := range args {
			value := testCase[arg]

			switch argType {
			case "int[]":
				intArray := value.([]int)
				intArrayValues := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(intArray)), ", "), "[]")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", new int[]{%s});\n", index, arg, intArrayValues))

			case "int[][]":
				int2DArray := value.([][]int)
				int2DArrayValues := ""
				for _, arr := range int2DArray {
					int2DArrayValues += fmt.Sprintf("new int[]{%s}, ", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]"))
				}
				int2DArrayValues = strings.TrimSuffix(int2DArrayValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", new int[][]{%s});\n", index, arg, int2DArrayValues))

			case "List<int[]>":
				listOfIntArray := value.([][]int)
				listValues := ""
				for _, arr := range listOfIntArray {
					listValues += fmt.Sprintf("new int[]{%s}, ", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]"))
				}
				listValues = strings.TrimSuffix(listValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", Arrays.asList(%s));\n", index, arg, listValues))

			case "List<List<int>>":
				listOfLists := value.([][]int)
				listValues := ""
				for _, list := range listOfLists {
					listValues += fmt.Sprintf("Arrays.asList(%s), ", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), ", "), "[]"))
				}
				listValues = strings.TrimSuffix(listValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", Arrays.asList(%s));\n", index, arg, listValues))

			case "List<string[]>":
				listOfStringArray := value.([][]string)
				listValues := ""
				for _, arr := range listOfStringArray {
					listValues += fmt.Sprintf("new String[]{%s}, ", strings.Join(arr, ", "))
				}
				listValues = strings.TrimSuffix(listValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", Arrays.asList(%s));\n", index, arg, listValues))

			case "string[][]":
				string2DArray := value.([][]string)
				string2DArrayValues := ""
				for _, arr := range string2DArray {
					string2DArrayValues += fmt.Sprintf("new String[]{%s}, ", strings.Join(arr, ", "))
				}
				string2DArrayValues = strings.TrimSuffix(string2DArrayValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", new String[][]{%s});\n", index, arg, string2DArrayValues))

			case "List<List<string>>":
				listOfStringLists := value.([][]string)
				listValues := ""
				for _, list := range listOfStringLists {
					listValues += fmt.Sprintf("Arrays.asList(%s), ", strings.Join(list, ", "))
				}
				listValues = strings.TrimSuffix(listValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", Arrays.asList(%s));\n", index, arg, listValues))

			case "int":
				intValue := value.(int)
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", %d);\n", index, arg, intValue))

			case "string":
				stringValue := value.(string)
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", \"%s\");\n", index, arg, stringValue))

			case "TreeNode":
				switch value.(type) {
				case int:
					result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", %d);\n", index, arg, value))
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
					result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", DSAHelpers.listToTree(new Integer[]{%s}));\n", index, arg, treeNodeListValues))
				}

			case "ListNode":
				listNodeValues := value.([]interface{})
				listNodeValuesStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(listNodeValues)), ", "), "[]")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", DSAHelpers.listToLinkedList(new int[]{%s}));\n", index, arg, listNodeValuesStr))

			case "string[]":
				stringInterfaceArray := value.([]interface{})
				var stringArray []string
				for _, v := range stringInterfaceArray {
					stringArray = append(stringArray, v.(string))
				}
				stringArrayValues := strings.Join(stringArray, "\", \"")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", Arrays.asList(\"%s\"));\n", index, arg, stringArrayValues))

			case "char[]":
				charInterfaceArray := value.([]interface{})
				var charArray []string
				for _, v := range charInterfaceArray {
					charArray = append(charArray, v.(string))
				}
				charArrayValues := strings.Join(charArray, "\", \"")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", Arrays.asList(%s));\n", index, arg, charArrayValues))

			case "char[][]":
				char2DInterfaceArray := value.([]interface{})
				var char2DArrayValues string
				for _, arr := range char2DInterfaceArray {
					charArray := arr.([]interface{})
					var charArrayStrings []string
					for _, v := range charArray {
						charArrayStrings = append(charArrayStrings, v.(string))
					}
					char2DArrayValues += fmt.Sprintf("Arrays.asList(%s), ", strings.Join(charArrayStrings, "\", \""))
				}
				char2DArrayValues = strings.TrimSuffix(char2DArrayValues, ", ")
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", Arrays.asList(%s));\n", index, arg, char2DArrayValues))

			default:
				result.WriteString(fmt.Sprintf("testCase%d.put(\"%s\", %v);\n", index, arg, value))
			}
		}

		result.WriteString(fmt.Sprintf("\ntestCases.add(testCase%d);\n", index))
	}

	return result.String()
}

func generateJavaArgNames(args map[string]string) string {
	keys := make([]string, 0, len(args))
	for key := range args {
		keys = append(keys, fmt.Sprintf("\"%s\"", key))
	}

	result := strings.Join(keys, ", ")
	return result
}

func generateJavaMethodArgs(args map[string]string) string {
	var result []string
	index := 0
	for _, arg := range args {
		argType := utils.TypeMappings["java"][arg]
		result = append(result, fmt.Sprintf("((%s) methodArgs[%d])", argType, index))
		index += 1
	}
	return strings.Join(result, ", ")
}

func generateJavaArgs(args map[string]string) string {
	var result []string
	result = append(result, "Map<String, String> javaArgs = new HashMap<>();")
	for argName, argType := range args {
		result = append(result, fmt.Sprintf("javaArgs.put(\"%s\", \"%s\");", argName, argType))
	}
	return strings.Join(result, "\n")
}
