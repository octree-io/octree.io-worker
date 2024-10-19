package testharness

import (
	"fmt"
	"log"

	utils "octree.io-worker/internal/utils/converters"
)

func RubyHarness(code string, args map[string]string, testCases []map[string]interface{}, returnType string) string {
	rubyArgs, err := utils.JsonToRuby(args)
	if err != nil {
		log.Fatal("Error converting JSON to Ruby")
	}

	rubyTestCases, err := utils.JsonToRuby(testCases)
	if err != nil {
		log.Fatal("Error converting JSON to Ruby")
	}

	rubyCode := fmt.Sprintf(`class ListNode
    attr_accessor :val, :next

    def initialize(val = 0, nxt = nil)
        @val = val
        @next = nxt
    end
end

class TreeNode
    attr_accessor :val, :left, :right

    def initialize(val = 0, left = nil, right = nil)
        @val = val
        @left = left
        @right = right
    end
end

class GraphNode
    attr_accessor :val, :neighbors

    def initialize(val = 0, neighbors = [])
        @val = val
        @neighbors = neighbors
    end
end

%s

def list_to_tree(lst)
    return nil if lst.empty? || lst[0].nil?

    root = TreeNode.new(lst[0])
    queue = [root]
    index = 1

    while !queue.empty? && index < lst.length
        node = queue.shift

        if index < lst.length && !lst[index].nil?
            node.left = TreeNode.new(lst[index])
            queue.push(node.left)
        end
        index += 1

        if index < lst.length && !lst[index].nil?
            node.right = TreeNode.new(lst[index])
            queue.push(node.right)
        end
        index += 1
    end

    root
end

def tree_to_list(root)
    return [] if root.nil?
    result = []
    queue = [root]

    while !queue.empty?
        node = queue.shift
        if node
            result << node.val
            queue.push(node.left)
            queue.push(node.right)
        else
            result << nil
        end
    end

    while result[-1].nil?
        result.pop
    end

    result
end

def find_node_by_value(root, value)
    return nil if root.nil?
    return root if root.val == value

    left_result = find_node_by_value(root.left, value)
    return left_result unless left_result.nil?

    find_node_by_value(root.right, value)
end

def list_to_linked_list(lst)
    dummy = ListNode.new(0)
    tail = dummy

    lst.each do |val|
        tail.next = ListNode.new(val)
        tail = tail.next
    end

    dummy.next
end

def linked_list_to_list(ll)
    lst = []
    current = ll
    while current
        lst << current.val
        current = current.next
    end
    lst
end

def custom_print(result)
    if result.nil?
        puts "null"
    else
        puts result
    end
end

def run_test_cases
    ruby_args = %s
    test_cases = %s
    return_type = "%s"

    test_cases.each do |test_case|
        arg_names = ruby_args.keys
        method_args = []
        root = nil

        if test_case.key?("root")
            root = list_to_tree(test_case["root"])
            method_args << root
        end

        arg_names.each do |arg|
            arg_type = ruby_args[arg]
            value = test_case[arg]

            if arg_type == 'TreeNode'
                if value.is_a?(Array)
                    value = list_to_tree(value)
                elsif value.is_a?(Integer)
                    value = find_node_by_value(root, value)
                end
            elsif arg_type == "ListNode"
                value = list_to_linked_list(value) if value.is_a?(Array)
            end

            method_args << value unless arg == "root"
        end

        result = solve(*method_args)

        if result.is_a?(TreeNode)
					if return_type == "TreeNode-int"
						puts result.val
					else
						puts "[" + tree_to_list(result).join(",") + "]"
					end
				elsif result.is_a?(ListNode)
					puts "[" + linked_list_to_list(result).join(",") + "]"
				elsif result.is_a?(Array)
					# Check if it's a nested list (Array of Arrays)
					if result.all? { |el| el.is_a?(Array) }
						# Handle nested arrays
						nested_result = result.map do |arr|
							arr.map { |item| item.is_a?(String) ? "\"#{item}\"" : item.to_s }.join(",")
						end
						puts "[" + nested_result.map { |arr| "[#{arr}]" }.join(", ") + "]"
					elsif result.all? { |el| el.is_a?(String) }
						# Handle array of strings
						puts "[" + result.map { |el| "\"#{el}\"" }.join(", ") + "]"
					else
						# Handle regular arrays (with possible strings)
						formatted_result = result.map { |el| el.is_a?(String) ? "\"#{el}\"" : el.to_s }
						puts "[" + formatted_result.join(", ") + "]"
					end
				else
					if result.nil? && (return_type == "ListNode" || return_type == "TreeNode")
						puts "[]"
					elsif result.is_a?(String)
						# Wrap the string in quotes
						puts "\"#{result}\""
					else
						custom_print(result)
					end
				end
    end
end

run_test_cases
`, code, rubyArgs, rubyTestCases, returnType)

	return rubyCode
}
