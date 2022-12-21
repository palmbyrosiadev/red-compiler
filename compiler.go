/*

RED - A simple, stack-based programming language

Copyright (C) 2022  The RED Authors

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var dtypes []string = []string{"int", "string", "bool", "float"}
var stack []stackVal
var symbols map[string]stackVal
var modules map[string]mod
var frommod bool = false
var module mod = mod{}
var modname string
var tempstack []stackVal

type mod struct {
	funcs   map[string]funct
	extvars map[string]stackVal
	symbols map[string]stackVal
}

type funct struct {
	name      string
	args      []string
	body      []string
	condition string
}
type stackVal struct {
	val    int
	symbol string
	sval   string
	dtype  int
	bval   bool
}

func runmod(code string) {
	for name, v := range module.extvars {
		module.symbols[name] = v
	}
	// Split the line into parts
	parts := strings.Split(code, " ")

	// Determine the operation
	op := parts[0]

	// Execute the operation
	switch op {
	case "PUSH":
		var s stackVal = stackVal{}

		// Push the value onto the temptempstack
		val, err := strconv.Atoi(parts[1])
		if err != nil {
			if strings.HasPrefix(strings.Join(parts[1:], " "), "\"") && strings.HasSuffix(strings.Join(parts[1:], " "), "\"") {
				s.sval = strings.Join(parts[1:], " ")
				s.dtype = 1
				tempstack = append(tempstack, s)
			} else if strings.HasPrefix(strings.Join(parts[1:], " "), "'") && strings.HasSuffix(strings.Join(parts[1:], " "), "'") {
				s.sval = strings.Join(parts[1:], " ")
				s.dtype = 1
				tempstack = append(tempstack, s)
			} else if parts[1] == "true" || parts[1] == "false" {
				s.dtype = 2
				if parts[1] == "true" {
					s.bval = true
				} else {
					s.bval = false
				}
				tempstack = append(tempstack, s)
			}

			/*
				fmt.Println(err)
				os.Exit(1)
			*/
		} else {
			s.val = val
			s.dtype = 0
			tempstack = append(tempstack, s)
		}

	case "ADD":
		// Pop the top two values from the tempstack and add them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]

		if val1.dtype == 1 || val2.dtype == 1 {
			fmt.Println("Cannot add strings")
			os.Exit(1)
		}

		tempstack = append(tempstack, stackVal{val: val1.val + val2.val, dtype: 0})
	case "SUB":
		// Pop the top two values from the tempstack and subtract them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]

		if val1.dtype == 1 || val2.dtype == 1 {
			fmt.Println("Cannot subtract strings")
			os.Exit(1)
		}

		tempstack = append(tempstack, stackVal{val: val1.val - val2.val, dtype: 0})
	case "MULT":
		// Pop the top two values from the tempstack and multiply them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]

		if val1.dtype == 1 || val2.dtype == 1 {
			fmt.Println("Cannot multiply strings")
			os.Exit(1)
		}

		tempstack = append(tempstack, stackVal{val: val1.val * val2.val, dtype: 0})
	case "DIV":
		// Pop the top two values from the tempstack and divide them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]

		if val2.val == 0 {
			fmt.Println("Cannot divide by zero")
			os.Exit(1)
		}

		if val1.dtype == 1 || val2.dtype == 1 {
			fmt.Println("Cannot divide strings")
			os.Exit(1)
		}

		tempstack = append(tempstack, stackVal{val: val1.val / val2.val, dtype: 0})
	case "STORE":
		// Pop the top value from the tempstack and store it in the symbol table
		val := tempstack[len(tempstack)-1]
		tempstack = tempstack[:len(tempstack)-1]
		module.symbols[parts[1]] = val
	case "LOAD":
		// Load the value from the symbol table and push it onto the tempstack
		val, ok := module.symbols[parts[1]]
		if !ok {
			fmt.Printf("Undefined symbol: %s\n", parts[1])
			os.Exit(1)
		}
		tempstack = append(tempstack, val)
	case "PRINT":
		// Pop the top value from the tempstack and print it
		val := tempstack[len(tempstack)-1]
		tempstack = tempstack[:len(tempstack)-1]
		if val.dtype == 1 {
			fmt.Println(val.sval)
		} else if val.dtype == 2 {
			fmt.Println(val.bval)
		} else {
			fmt.Println(val.val)
		}
	case "STR":
		var s stackVal = stackVal{}
		s.dtype = 1
		val := tempstack[len(tempstack)-1]
		tempstack = tempstack[:len(tempstack)-1]
		if val.dtype == 0 {
			s.sval = strconv.Itoa(val.val)
		} else if val.dtype == 2 {
			s.sval = strconv.FormatBool(val.bval)
		} else {
			s.sval = val.sval
		}
		tempstack = append(tempstack, s)
	case "INT":
		var s stackVal = stackVal{}
		s.dtype = 0
		val := tempstack[len(tempstack)-1]
		tempstack = tempstack[:len(tempstack)-1]
		if val.dtype == 1 {
			i, err := strconv.Atoi(val.sval)
			if err != nil {
				fmt.Println("Cannot convert string to int")
				os.Exit(1)
			}
			s.val = i
		} else if val.dtype == 2 {
			fmt.Println("Cannot convert bool to int")
			os.Exit(1)
		} else {
			s.val = val.val
		}
		tempstack = append(tempstack, s)
	case "BOOL":
		var s stackVal = stackVal{}
		s.dtype = 2
		val := tempstack[len(tempstack)-1]
		tempstack = tempstack[:len(tempstack)-1]
		if val.dtype == 1 {
			b, err := strconv.ParseBool(val.sval)
			if err != nil {
				fmt.Println("Cannot convert string to bool")
				os.Exit(1)
			}
			s.bval = b
		} else if val.dtype == 0 {
			fmt.Println("Cannot convert int to bool")
			os.Exit(1)
		} else {
			s.bval = val.bval
		}
		tempstack = append(tempstack, s)
	case "STRCAT":
		// Pop the top two values from the tempstack and concatenate them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]
		if !(val1.dtype == 1 || val2.dtype == 1) {
			fmt.Println("Cannot concatenate non-strings")
			os.Exit(1)
		}
		tempstack = append(tempstack, stackVal{dtype: 1, sval: val1.sval + val2.sval})
	case "EQ":
		// Pop the top two values from the tempstack and compare them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.val == val2.val})
		} else if val1.dtype == 1 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.sval == val2.sval})
		} else {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.bval == val2.bval})
		}
	case "NEQ":
		// Pop the top two values from the tempstack and compare them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.val != val2.val})
		} else if val1.dtype == 1 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.sval != val2.sval})
		} else {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.bval != val2.bval})
		}
	case "GT":
		// Pop the top two values from the tempstack and compare them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.val > val2.val})
		} else if val1.dtype == 1 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.sval > val2.sval})
		} else {
			fmt.Println("Cannot compare bools")
			os.Exit(1)
		}
	case "GTE":
		// Pop the top two values from the tempstack and compare them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.val >= val2.val})
		} else if val1.dtype == 1 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.sval >= val2.sval})
		} else {
			fmt.Println("Cannot compare bools")
			os.Exit(1)
		}
	case "LT":
		// Pop the top two values from the tempstack and compare them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.val < val2.val})
		} else if val1.dtype == 1 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.sval < val2.sval})
		} else {
			fmt.Println("Cannot compare bools")
			os.Exit(1)
		}
	case "LTE":
		// Pop the top two values from the tempstack and compare them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.val <= val2.val})
		} else if val1.dtype == 1 {
			tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.sval <= val2.sval})
		} else {
			fmt.Println("Cannot compare bools")
			os.Exit(1)
		}

	default:
		fmt.Printf("Invalid operation: %s\n", op)
		os.Exit(1)
	}
	for name := range module.extvars {
		module.extvars[name] = module.symbols[name]
	}
	modules[modname] = module

}

func run(code string) {
	// Split the line into parts
	parts := strings.Split(code, " ")

	// Determine the operation
	op := parts[0]

	// Execute the operation
	switch op {
	case "PUSH":
		var s stackVal = stackVal{}

		// Push the value onto the stack
		val, err := strconv.Atoi(parts[1])
		if err != nil {
			if strings.HasPrefix(strings.Join(parts[1:], " "), "\"") && strings.HasSuffix(strings.Join(parts[1:], " "), "\"") {
				s.sval = strings.Join(parts[1:], " ")
				s.dtype = 1
				stack = append(stack, s)
			} else if strings.HasPrefix(strings.Join(parts[1:], " "), "'") && strings.HasSuffix(strings.Join(parts[1:], " "), "'") {
				s.sval = strings.Join(parts[1:], " ")
				s.dtype = 1
				stack = append(stack, s)
			} else if parts[1] == "true" || parts[1] == "false" {
				s.dtype = 2
				if parts[1] == "true" {
					s.bval = true
				} else {
					s.bval = false
				}
				stack = append(stack, s)
			}

			/*
				fmt.Println(err)
				os.Exit(1)
			*/
		} else {
			s.val = val
			s.dtype = 0
			stack = append(stack, s)
		}

	case "ADD":
		// Pop the top two values from the stack and add them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if val1.dtype == 1 || val2.dtype == 1 {
			fmt.Println("Cannot add strings")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val + val2.val, dtype: 0})
	case "SUB":
		// Pop the top two values from the stack and subtract them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if val1.dtype == 1 || val2.dtype == 1 {
			fmt.Println("Cannot subtract strings")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val - val2.val, dtype: 0})
	case "MULT":
		// Pop the top two values from the stack and multiply them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if val1.dtype == 1 || val2.dtype == 1 {
			fmt.Println("Cannot multiply strings")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val * val2.val, dtype: 0})
	case "DIV":
		// Pop the top two values from the stack and divide them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if val2.val == 0 {
			fmt.Println("Cannot divide by zero")
			os.Exit(1)
		}

		if val1.dtype == 1 || val2.dtype == 1 {
			fmt.Println("Cannot divide strings")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val / val2.val, dtype: 0})
	case "STORE":
		// Pop the top value from the stack and store it in the symbol table
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		symbols[parts[1]] = val
	case "LOAD":
		// Load the value from the symbol table and push it onto the stack
		val, ok := symbols[parts[1]]
		if !ok {
			fmt.Printf("Undefined symbol: %s\n", parts[1])
			os.Exit(1)
		}
		stack = append(stack, val)
	case "PRINT":
		// Pop the top value from the stack and print it
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype == 1 {
			fmt.Println(val.sval)
		} else if val.dtype == 2 {
			fmt.Println(val.bval)
		} else {
			fmt.Println(val.val)
		}
	case "STR":
		var s stackVal = stackVal{}
		s.dtype = 1
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype == 0 {
			s.sval = strconv.Itoa(val.val)
		} else if val.dtype == 2 {
			s.sval = strconv.FormatBool(val.bval)
		} else {
			s.sval = val.sval
		}
		stack = append(stack, s)
	case "INT":
		var s stackVal = stackVal{}
		s.dtype = 0
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype == 1 {
			i, err := strconv.Atoi(val.sval)
			if err != nil {
				fmt.Println("Cannot convert string to int")
				os.Exit(1)
			}
			s.val = i
		} else if val.dtype == 2 {
			fmt.Println("Cannot convert bool to int")
			os.Exit(1)
		} else {
			s.val = val.val
		}
		stack = append(stack, s)
	case "BOOL":
		var s stackVal = stackVal{}
		s.dtype = 2
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype == 1 {
			b, err := strconv.ParseBool(val.sval)
			if err != nil {
				fmt.Println("Cannot convert string to bool")
				os.Exit(1)
			}
			s.bval = b
		} else if val.dtype == 0 {
			fmt.Println("Cannot convert int to bool")
			os.Exit(1)
		} else {
			s.bval = val.bval
		}
		stack = append(stack, s)
	case "STRCAT":
		// Pop the top two values from the stack and concatenate them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if !(val1.dtype == 1 || val2.dtype == 1) {
			fmt.Println("Cannot concatenate non-strings")
			os.Exit(1)
		}
		stack = append(stack, stackVal{dtype: 1, sval: val1.sval + val2.sval})
	case "EQ":
		// Pop the top two values from the stack and compare them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.val == val2.val})
		} else if val1.dtype == 1 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.sval == val2.sval})
		} else {
			stack = append(stack, stackVal{dtype: 2, bval: val1.bval == val2.bval})
		}
	case "NEQ":
		// Pop the top two values from the stack and compare them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.val != val2.val})
		} else if val1.dtype == 1 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.sval != val2.sval})
		} else {
			stack = append(stack, stackVal{dtype: 2, bval: val1.bval != val2.bval})
		}
	case "GT":
		// Pop the top two values from the stack and compare them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.val > val2.val})
		} else if val1.dtype == 1 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.sval > val2.sval})
		} else {
			fmt.Println("Cannot compare bools")
			os.Exit(1)
		}
	case "GTE":
		// Pop the top two values from the stack and compare them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.val >= val2.val})
		} else if val1.dtype == 1 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.sval >= val2.sval})
		} else {
			fmt.Println("Cannot compare bools")
			os.Exit(1)
		}
	case "LT":
		// Pop the top two values from the stack and compare them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.val < val2.val})
		} else if val1.dtype == 1 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.sval < val2.sval})
		} else {
			fmt.Println("Cannot compare bools")
			os.Exit(1)
		}
	case "LTE":
		// Pop the top two values from the stack and compare them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != val2.dtype {
			fmt.Println("Cannot compare different types")
			os.Exit(1)
		}
		if val1.dtype == 0 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.val <= val2.val})
		} else if val1.dtype == 1 {
			stack = append(stack, stackVal{dtype: 2, bval: val1.sval <= val2.sval})
		} else {
			fmt.Println("Cannot compare bools")
			os.Exit(1)
		}

	default:
		fmt.Printf("Invalid operation: %s\n", op)
		os.Exit(1)

	}
}

func main() {
	// Read the input file
	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Split the input into lines
	lines := strings.Split(string(bytes), "\n")

	// Initialize the symbol table and the stack
	modules = make(map[string]mod)
	symbols = make(map[string]stackVal)
	funcs := make(map[string]funct)
	stack = make([]stackVal, 0)
	activefuncwrite := false
	var activefunc funct = funct{}
	var running funct = funct{}
	var runningfunc bool = false

	// Iterate through the lines and compile them
	for _, line := range lines {
		// Skip empty lines
		if line == "" && !runningfunc {
			continue
		}

		// Split the line into parts
		parts := strings.Split(line, " ")

		// Determine the operation
		op := parts[0]
		if activefuncwrite == true {
			if op == "ENDFUNC" {

				activefuncwrite = false
				funcs[activefunc.name] = activefunc
				activefunc = funct{}
				continue
			} else {
				activefunc.body = append(activefunc.body, line)
				continue
			}
		}

		if runningfunc == true {
			if frommod {

				if running.condition == "" {
					for _, code := range running.body {
						runmod(code)
					}
				} else {
					for module.symbols[running.condition].bval {
						for _, code := range running.body {
							runmod(code)
						}
					}
				}
			} else {
				if running.condition == "" {
					for _, code := range running.body {
						run(code)
					}
				} else {
					for symbols[running.condition].bval {
						for _, code := range running.body {
							run(code)
						}
					}
				}
			}
			tempstack = make([]stackVal, 0)
			runningfunc = false
			running = funct{}
			continue
		}

		// Execute the operation
		switch op {
		case "PUSH":
			var s stackVal = stackVal{}

			// Push the value onto the stack
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				if strings.HasPrefix(parts[1], "\"") && strings.HasSuffix(parts[1], "\"") {
					s.sval = parts[1]
					s.dtype = 1
					stack = append(stack, s)
				} else if strings.HasPrefix(parts[1], "'") && strings.HasSuffix(parts[1], "'") {
					s.sval = parts[1]
					s.dtype = 1
					stack = append(stack, s)
				} else if parts[1] == "true" || parts[1] == "false" {
					s.dtype = 2
					if parts[1] == "true" {
						s.bval = true
					} else {
						s.bval = false
					}
					stack = append(stack, s)
				}

				/*
					fmt.Println(err)
					os.Exit(1)
				*/
			} else {
				s.val = val
				s.dtype = 0
				stack = append(stack, s)
			}

		case "ADD":
			// Pop the top two values from the stack and add them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if val1.dtype == 1 || val2.dtype == 1 {
				fmt.Println("Cannot add strings")
				os.Exit(1)
			}

			stack = append(stack, stackVal{val: val1.val + val2.val, dtype: 0})
		case "SUB":
			// Pop the top two values from the stack and subtract them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if val1.dtype == 1 || val2.dtype == 1 {
				fmt.Println("Cannot subtract strings")
				os.Exit(1)
			}

			stack = append(stack, stackVal{val: val1.val - val2.val, dtype: 0})
		case "MULT":
			// Pop the top two values from the stack and multiply them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if val1.dtype == 1 || val2.dtype == 1 {
				fmt.Println("Cannot multiply strings")
				os.Exit(1)
			}

			stack = append(stack, stackVal{val: val1.val * val2.val, dtype: 0})
		case "DIV":
			// Pop the top two values from the stack and divide them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if val2.val == 0 {
				fmt.Println("Cannot divide by zero")
				os.Exit(1)
			}

			if val1.dtype == 1 || val2.dtype == 1 {
				fmt.Println("Cannot divide strings")
				os.Exit(1)
			}

			stack = append(stack, stackVal{val: val1.val / val2.val, dtype: 0})
		case "STORE":
			// Pop the top value from the stack and store it in the symbol table
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			symbols[parts[1]] = val
		case "LOAD":
			// Load the value from the symbol table and push it onto the stack
			val, ok := symbols[parts[1]]
			if !ok {
				fmt.Printf("Undefined symbol: %s\n", parts[1])
				os.Exit(1)
			}
			stack = append(stack, val)
		case "PRINT":
			// Pop the top value from the stack and print it
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if val.dtype == 1 {
				fmt.Println(val.sval)
			} else if val.dtype == 2 {
				fmt.Println(val.bval)
			} else {
				fmt.Println(val.val)
			}
		case "STR":
			var s stackVal = stackVal{}
			s.dtype = 1
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if val.dtype == 0 {
				s.sval = strconv.Itoa(val.val)
			} else if val.dtype == 2 {
				s.sval = strconv.FormatBool(val.bval)
			} else {
				s.sval = val.sval
			}
			stack = append(stack, s)
		case "INT":
			var s stackVal = stackVal{}
			s.dtype = 0
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if val.dtype == 1 {
				i, err := strconv.Atoi(val.sval)
				if err != nil {
					fmt.Println("Cannot convert string to int")
					os.Exit(1)
				}
				s.val = i
			} else if val.dtype == 2 {
				fmt.Println("Cannot convert bool to int")
				os.Exit(1)
			} else {
				s.val = val.val
			}
			stack = append(stack, s)
		case "BOOL":
			var s stackVal = stackVal{}
			s.dtype = 2
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if val.dtype == 1 {
				b, err := strconv.ParseBool(val.sval)
				if err != nil {
					fmt.Println("Cannot convert string to bool")
					os.Exit(1)
				}
				s.bval = b
			} else if val.dtype == 0 {
				fmt.Println("Cannot convert int to bool")
				os.Exit(1)
			} else {
				s.bval = val.bval
			}
			stack = append(stack, s)
		case "STRCAT":
			// Pop the top two values from the stack and concatenate them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if !(val1.dtype == 1 || val2.dtype == 1) {
				fmt.Println("Cannot concatenate non-strings")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 1, sval: val1.sval + val2.sval})
		case "EQ":
			// Pop the top two values from the stack and compare them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if val1.dtype != val2.dtype {
				fmt.Println("Cannot compare different types")
				os.Exit(1)
			}
			if val1.dtype == 0 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.val == val2.val})
			} else if val1.dtype == 1 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.sval == val2.sval})
			} else {
				stack = append(stack, stackVal{dtype: 2, bval: val1.bval == val2.bval})
			}
		case "NEQ":
			// Pop the top two values from the stack and compare them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if val1.dtype != val2.dtype {
				fmt.Println("Cannot compare different types")
				os.Exit(1)
			}
			if val1.dtype == 0 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.val != val2.val})
			} else if val1.dtype == 1 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.sval != val2.sval})
			} else {
				stack = append(stack, stackVal{dtype: 2, bval: val1.bval != val2.bval})
			}
		case "GT":
			// Pop the top two values from the stack and compare them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if val1.dtype != val2.dtype {
				fmt.Println("Cannot compare different types")
				os.Exit(1)
			}
			if val1.dtype == 0 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.val > val2.val})
			} else if val1.dtype == 1 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.sval > val2.sval})
			} else {
				fmt.Println("Cannot compare bools")
				os.Exit(1)
			}
		case "GTE":
			// Pop the top two values from the stack and compare them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if val1.dtype != val2.dtype {
				fmt.Println("Cannot compare different types")
				os.Exit(1)
			}
			if val1.dtype == 0 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.val >= val2.val})
			} else if val1.dtype == 1 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.sval >= val2.sval})
			} else {
				fmt.Println("Cannot compare bools")
				os.Exit(1)
			}
		case "LT":
			// Pop the top two values from the stack and compare them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if val1.dtype != val2.dtype {
				fmt.Println("Cannot compare different types")
				os.Exit(1)
			}
			if val1.dtype == 0 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.val < val2.val})
			} else if val1.dtype == 1 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.sval < val2.sval})
			} else {
				fmt.Println("Cannot compare bools")
				os.Exit(1)
			}
		case "LTE":
			// Pop the top two values from the stack and compare them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if val1.dtype != val2.dtype {
				fmt.Println("Cannot compare different types")
				os.Exit(1)
			}
			if val1.dtype == 0 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.val <= val2.val})
			} else if val1.dtype == 1 {
				stack = append(stack, stackVal{dtype: 2, bval: val1.sval <= val2.sval})
			} else {
				fmt.Println("Cannot compare bools")
				os.Exit(1)
			}
		case "FUNC":
			activefuncwrite = true
			activefunc.name = parts[1]

		case "RUN":
			for i, f := range funcs {
				if i == parts[1] {
					runningfunc = true
					running = f
					break
				}
			}
			if len(parts) > 2 {
				running.condition = ""
				for l, v := range symbols {
					if l == parts[2] {
						if v.dtype == 2 {
							running.condition = l
							break
						} else {
							fmt.Printf("Cannot use %s as condition\n", parts[2])
							os.Exit(1)
						}
					}
				}
				if running.condition == "" {
					fmt.Printf("No such symbol: %s\n", parts[2])
					os.Exit(1)
				}
				if len(parts) > 3 {
					for i := 3; i < len(parts); i++ {
						running.args = append(running.args, parts[i])
					}
				}
			}
			if !runningfunc {
				fmt.Printf("No such function: %s\n", parts[1])
				os.Exit(1)
			}
		case "EXIT":
			os.Exit(0)
		case "MODSTORE":
			// Store a value in exported module variable
			if len(parts) < 3 {
				fmt.Println("Invalid syntax")
				os.Exit(1)
			}
			if parts[1] == "" || parts[2] == "" {
				fmt.Println("Invalid syntax")
				os.Exit(1)
			}
			for name, m := range modules {
				if name == parts[1] {
					val := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					m.extvars[parts[2]] = val
				}
			}

		case "MODGET":
			// Get a value from exported module variable
			if len(parts) < 3 {
				fmt.Println("Invalid syntax")
				os.Exit(1)
			}
			if parts[1] == "" || parts[2] == "" {
				fmt.Println("Invalid syntax")
				os.Exit(1)
			}
			for name, m := range modules {
				if name == parts[1] {
					for n, modu := range m.extvars {
						if n == parts[2] {
							stack = append(stack, modu)
						}
					}
				}
			}

		case "MODRUN":
			// Run a function from a module
			if len(parts) < 3 {
				fmt.Println("Invalid syntax")
				os.Exit(1)
			}
			if parts[1] == "" || parts[2] == "" {
				fmt.Println("Invalid syntax")
				os.Exit(1)
			}
			running.condition = ""
			for name, m := range modules {
				if name == parts[1] {
					for n, modu := range m.funcs {
						if n == parts[2] {
							runningfunc = true
							running = modu
							frommod = true
							module = m
							modname = name
							break
						}
					}
				}
			}
			if len(parts) > 3 {
				for l, v := range module.symbols {
					if l == parts[3] {
						if v.dtype == 2 {
							running.condition = l
							break
						} else {
							fmt.Printf("Cannot use %s as condition\n", parts[3])
							os.Exit(1)
						}
					}
				}
				if running.condition == "" {
					fmt.Printf("No such symbol: %s\n", parts[3])
					os.Exit(1)
				}
			}
		case "IMPORT":
			// Import a file
			bytes, err := ioutil.ReadFile(parts[1])
			if err != nil {
				fmt.Println("Invalid module")
				os.Exit(1)
			}
			lines := strings.Split(string(bytes), "\n")
			m := mod{funcs: make(map[string]funct), symbols: make(map[string]stackVal, 0), extvars: make(map[string]stackVal, 0)}
			active := false
			activef := funct{}
			for _, line := range lines {
				if active {
					if strings.HasPrefix(line, "ENDFUNC") {
						active = false
						m.funcs[activef.name] = activef
						activef = funct{}
						continue
					} else {
						activef.body = append(activef.body, line)
						continue
					}
				}
				if line != "" {
					partsin := strings.Split(line, " ")
					switch partsin[0] {
					case "EXPORT":
						// Export a symbol
						if len(partsin) > 2 {
							val, err := strconv.Atoi(partsin[2])
							if err != nil {
								fmt.Println("Invalid export")
								os.Exit(1)
							}
							var s stackVal = stackVal{dtype: 0, val: val}
							m.extvars[partsin[1]] = s
						} else {
							fmt.Println("Invalid export")
							os.Exit(1)
						}
					case "SET":
						// Set a symbol
						if len(partsin) > 2 {
							if partsin[2] == "true" || partsin[2] == "false" {
								if partsin[2] == "true" {
									m.symbols[partsin[1]] = stackVal{dtype: 2, bval: true}
								} else {
									m.symbols[partsin[1]] = stackVal{dtype: 2, bval: false}
								}
							} else {
								fmt.Println("SET is used for boolean values only")
								os.Exit(1)
							}
						} else {
							fmt.Println("Invalid set")
							os.Exit(1)
						}
					case "FUNC":
						active = true
						activef.name = partsin[1]

					default:
						fmt.Println("Invalid module")
						os.Exit(1)
					}

				}
			}
			modules[parts[2]] = m
		default:
			fmt.Printf("Invalid operation: %s\n", op)
			os.Exit(1)

		}
	}
}
