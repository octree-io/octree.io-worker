package testharness

import (
	"fmt"
	"strings"
)

func CppHarness(code string, args map[string]string, testCases []map[string]interface{}, returnType string) string {
	harnessCode := fmt.Sprintf(`
#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <map>
#include <unordered_map>
#include <set>
#include <cmath>
#include <queue>
#include <stack>
#include <deque>
#include <cstdlib>
#include <climits>
#include <any>
#include <optional>

using namespace std;

std::string returnType = "%s";

struct ListNode {
    int val;
    ListNode* next;
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode* next) : val(x), next(next) {}
};

struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode* left, TreeNode* right) : val(x), left(left), right(right) {}
};

struct GraphNode {
    int val;
    std::vector<GraphNode*> neighbors;
    GraphNode(int x) : val(x) {}
    GraphNode(int x, const std::vector<GraphNode*>& neighbors) : val(x), neighbors(neighbors) {}
};

%s

TreeNode* list_to_tree(const std::vector<std::optional<int>>& lst) {
    if (lst.empty() || !lst[0].has_value()) return nullptr;

    TreeNode* root = new TreeNode(lst[0].value());
    std::queue<TreeNode*> queue;
    queue.push(root);
    size_t index = 1;

    while (!queue.empty() && index < lst.size()) {
        TreeNode* node = queue.front();
        queue.pop();

        if (index < lst.size() && lst[index].has_value()) {
            node->left = new TreeNode(lst[index].value());
            queue.push(node->left);
        }
        index++;

        if (index < lst.size() && lst[index].has_value()) {
            node->right = new TreeNode(lst[index].value());
            queue.push(node->right);
        }
        index++;
    }
    return root;
}

std::vector<std::optional<int>> tree_to_list(TreeNode* root) {
    std::vector<std::optional<int>> result;
    if (!root) return result;

    std::queue<TreeNode*> queue;
    queue.push(root);

    while (!queue.empty()) {
        TreeNode* node = queue.front();
        queue.pop();

        if (node) {
            result.push_back(node->val);
            queue.push(node->left);
            queue.push(node->right);
        } else {
            result.push_back(std::nullopt);  // Use std::nullopt to represent nullptr
        }
    }

    while (!result.empty() && !result.back().has_value()) {
        result.pop_back();
    }

    return result;
}

TreeNode* find_node_by_value(TreeNode* root, int value) {
    if (!root) return nullptr;
    if (root->val == value) return root;

    TreeNode* left_result = find_node_by_value(root->left, value);
    if (left_result) return left_result;

    return find_node_by_value(root->right, value);
}

ListNode* list_to_linked_list(const std::vector<int>& lst) {
    if (lst.empty()) return nullptr;

    ListNode* dummy = new ListNode(0);
    ListNode* tail = dummy;

    for (int val : lst) {
        tail->next = new ListNode(val);
        tail = tail->next;
    }

    return dummy->next;
}

class TestHelper {
public:
    static void printResult(const std::vector<int>& result) {
        for (int val : result) {
            std::cout << val << " ";
        }
        std::cout << std::endl;
    }

    static void printResult(int result) {
        std::cout << result << std::endl;
    }

    static void printResult(const std::string& result) {
        std::cout << result << std::endl;
    }

    static void printResult(bool result) {
        std::cout << std::boolalpha << result << std::endl;
    }

    static void printResult(const std::vector<std::vector<std::string>>& result) {
        std::cout << "[";
        for (size_t i = 0; i < result.size(); ++i) {
            const auto& vec = result[i];
            std::cout << "[";
            for (size_t j = 0; j < vec.size(); ++j) {
                std::cout << "\"" << vec[j] << "\"";
                if (j != vec.size() - 1) {
                    std::cout << ",";
                }
            }
            std::cout << "]";
            if (i != result.size() - 1) {
                std::cout << ",";
            }
        }
        std::cout << "]" << std::endl;
    }

    static void printResult(TreeNode* root) {
        if (returnType == "TreeNode-int") {
            std::cout << root->val << std::endl;
            return;
        }
        std::vector<std::optional<int>> treeList = tree_to_list(root);
        std::cout << "[";
        for (size_t i = 0; i < treeList.size(); ++i) {
            if (treeList[i].has_value()) {
                std::cout << treeList[i].value();
            } else {
                std::cout << "null";
            }
            if (i != treeList.size() - 1) {
                std::cout << ",";
            }
        }
        std::cout << "]" << std::endl;
    }

    static void printResult(ListNode* head) {
        std::cout << "[";
        ListNode* current = head;
        bool first = true;
        while (current) {
            if (!first) {
                std::cout << ",";
            }
            std::cout << current->val;
            current = current->next;
            first = false;
        }
        std::cout << "]" << std::endl;
    }
};

class TestHarness {
public:
    void run() {
        Solution solution;

        std::vector<std::map<std::string, std::any>> testCases;

        %s

        %s

        std::vector<std::string> argNames = { %s };

        for (int i = 0; i < testCases.size(); i++) {
            std::map<std::string, std::any> testCase = testCases[i];

            TreeNode* root = nullptr;

            if (testCase.find("root") != testCase.end()) {
                auto root_value = std::any_cast<TreeNode*>(testCase.at("root"));
                root = root_value;
            }

            std::vector<std::any> methodArgs;
            for (const std::string& argName : argNames) {
                auto argType = args[argName];
                auto value = testCase[argName];

                if (argType == "TreeNode*" && value.type() == typeid(int)) {
                    int node_value = std::any_cast<int>(value);
                    value = find_node_by_value(root, node_value);
                }

                methodArgs.push_back(value);
            }

            auto result = solution.solve(
                %s
            );

            TestHelper::printResult(result);
        }
    }
};

int main() {
    TestHarness testHarness;
    testHarness.run();
    return 0;
}
`, returnType, code, generateCppTestCases(args, testCases), generateCppArgs(args), generateCppArgNames(args), generateCppMethodArgs(args))

	return harnessCode
}

func generateCppTestCases(args map[string]string, testCases []map[string]interface{}) string {
	var result []string

	for i, testCase := range testCases {
		var caseLines []string
		caseLines = append(caseLines, fmt.Sprintf("std::map<std::string, std::any> testCase%d;", i))

		for arg, argType := range args {
			value := testCase[arg]

			switch argType {
			case "int[]":
				values := toIntSlice(value)
				caseLines = append(caseLines, fmt.Sprintf("testCase%d[\"%s\"] = std::vector<int>{%s};", i, arg, strings.Join(values, ", ")))
			case "int":
				caseLines = append(caseLines, fmt.Sprintf("testCase%d[\"%s\"] = %d;", i, arg, value))
			case "string":
				caseLines = append(caseLines, fmt.Sprintf("testCase%d[\"%s\"] = std::string(\"%s\");", i, arg, value))
			case "string[]":
				stringsArray := toStringSlice(value)
				caseLines = append(caseLines, fmt.Sprintf("testCase%d[\"%s\"] = std::vector<std::string>{%s};", i, arg, strings.Join(stringsArray, ", ")))
			case "TreeNode":
				switch value.(type) {
				case int:
					caseLines = append(caseLines, fmt.Sprintf("testCase%d[\"%s\"] = %d;", i, arg, value))
				default:
					treeValues := toOptionalIntSlice(value)
					caseLines = append(caseLines, fmt.Sprintf("testCase%d[\"%s\"] = list_to_tree(std::vector<std::optional<int>>{%s});", i, arg, strings.Join(treeValues, ", ")))
				}
			case "ListNode":
				listValues := toIntSlice(value)
				caseLines = append(caseLines, fmt.Sprintf("testCase%d[\"%s\"] = list_to_linked_list(std::vector<int>{%s});", i, arg, strings.Join(listValues, ", ")))
			default:
				caseLines = append(caseLines, fmt.Sprintf("testCase%d[\"%s\"] = %v;", i, arg, value))
			}
		}

		caseLines = append(caseLines, fmt.Sprintf("testCases.push_back(testCase%d);", i))
		result = append(result, strings.Join(caseLines, "\n"))
	}

	return strings.Join(result, "\n")
}

func generateCppArgs(args map[string]string) string {
	var result []string
	result = append(result, "std::map<std::string, std::string> args = {")
	for argName, argType := range args {
		cppType := getCppType(argType)
		result = append(result, fmt.Sprintf("{\"%s\", \"%s\"},", argName, cppType))
	}
	result = append(result, "};")
	return strings.Join(result, "\n")
}

func generateCppArgNames(args map[string]string) string {
	var result []string
	for arg := range args {
		result = append(result, fmt.Sprintf("\"%s\"", arg))
	}
	return strings.Join(result, ", ")
}

func generateCppMethodArgs(args map[string]string) string {
	var result []string
	index := 0
	for _, argType := range args {
		cppType := getCppType(argType)
		result = append(result, fmt.Sprintf("std::any_cast<%s>(methodArgs[%d])", cppType, index))
		index += 1
	}
	return strings.Join(result, ", ")
}

func getCppType(argType string) string {
	switch argType {
	case "int[]":
		return "std::vector<int>"
	case "int":
		return "int"
	case "string":
		return "std::string"
	case "string[]":
		return "std::vector<std::string>"
	case "string[][]":
		return "std::vector<std::vector<std::string>>"
	case "TreeNode":
		return "TreeNode*"
	case "ListNode":
		return "ListNode*"
	default:
		return argType
	}
}

func toOptionalIntSlice(value interface{}) []string {
	v := value.([]interface{})
	var result []string
	for _, num := range v {
		if num == nil {
			result = append(result, "std::nullopt")
		} else {
			result = append(result, fmt.Sprintf("%d", num))
		}
	}
	return result
}

func toIntSlice(value interface{}) []string {
	v := value.([]interface{})
	var result []string
	for _, num := range v {
		result = append(result, fmt.Sprintf("%d", num))
	}
	return result
}

func toStringSlice(value interface{}) []string {
	v := value.([]interface{})
	var result []string
	for _, str := range v {
		result = append(result, fmt.Sprintf("\"%s\"", str.(string)))
	}
	return result
}
