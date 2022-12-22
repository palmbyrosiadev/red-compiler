/*

RED - A simple, stack-based programming language

Copyright (C) 2022  The RED Authors

*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var stack []stackVal
var symbols map[string]stackVal
var modules map[string]mod
var frommod bool = false
var module mod = mod{}
var modname string
var tempstack []stackVal
var comment = false

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
	val    float64
	symbol string
	sval   string
	dtype  int
	bval   bool
	list   []stackVal
}

func defimports() {
	j, err := os.Open("built-in/util.kr")

	if err != nil {
		fmt.Println("Invalid keyword file")
		os.Exit(1)
	}

	defer j.Close()

	byteValue, _ := ioutil.ReadAll(j)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	var util keymod = keymod{cases: make(map[string][]interface{})}
	for _, v := range result["main"].([]interface{}) {
		util.cases[v.(map[string]interface{})["case"].(string)] = v.(map[string]interface{})["code"].([]interface{})
	}
	keymods[result["prefix"].(string)] = util
}

func runmod(code string) {
	if code == "\n" || code == "" {
		return
	}
	for name, v := range module.extvars {
		module.symbols[name] = v
	}
	// Split the line into parts
	code = strings.ReplaceAll(code, "	", "")
	code = strings.ReplaceAll(code, "    ", "")

	// Split the line into parts
	parts := strings.Split(code, " ")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
		parts[i] = strings.ReplaceAll(parts[i], "	", "")
	}

	// Determine the operation
	op := parts[0]

	// Execute the operation
	switch op {
	case "PUSH":
		var s stackVal = stackVal{}

		// Push the value onto the temptempstack
		val, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			if strings.HasPrefix(strings.Join(parts[1:], " "), "\"") && strings.HasSuffix(strings.Join(parts[1:], " "), "\"") {
				s.sval = strings.Join(parts[1:], " ")[1 : len(strings.Join(parts[1:], " "))-1]
				s.dtype = 1
				tempstack = append(tempstack, s)
			} else if strings.HasPrefix(strings.Join(parts[1:], " "), "'") && strings.HasSuffix(strings.Join(parts[1:], " "), "'") {
				s.sval = strings.Join(parts[1:], " ")[1 : len(strings.Join(parts[1:], " "))-1]
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

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
			os.Exit(1)
		}

		tempstack = append(tempstack, stackVal{val: val1.val + val2.val, dtype: 0})
	case "SUB":
		// Pop the top two values from the tempstack and subtract them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
			os.Exit(1)
		}

		tempstack = append(tempstack, stackVal{val: val1.val - val2.val, dtype: 0})
	case "MULT":
		// Pop the top two values from the tempstack and multiply them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
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

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
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
		if len(parts) > 2 {
			if val.dtype != 4 {
				fmt.Println("Cannot index non-array")
				os.Exit(1)
			}
			index, err := strconv.Atoi(parts[2])
			if err != nil {
				valt, ok := module.symbols[parts[2]]
				if !ok {
					fmt.Printf("Undefined symbol: %s\n", parts[2])
					os.Exit(1)
				} else if valt.dtype != 0 {
					fmt.Println("Index of array must be number")
					os.Exit(1)
				} else {
					index = int(valt.val)
				}
			}
			if index >= len(val.list) {
				fmt.Println("Index out of bounds")
				os.Exit(1)
			} else {
				val = val.list[index]
			}
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
		} else if val.dtype == 0 {
			fmt.Println(val.val)
		} else {
			fmt.Println("Cannot print element")
		}
	case "STR":
		var s stackVal = stackVal{}
		s.dtype = 1
		val := tempstack[len(tempstack)-1]
		tempstack = tempstack[:len(tempstack)-1]
		if val.dtype == 0 {
			s.sval = strconv.FormatFloat(val.val, 'f', -1, 64)
		} else if val.dtype == 2 {
			s.sval = strconv.FormatBool(val.bval)
		} else {
			s.sval = val.sval
		}
		tempstack = append(tempstack, s)
	case "FLOAT":
		var s stackVal = stackVal{}
		s.dtype = 0
		val := tempstack[len(tempstack)-1]
		tempstack = tempstack[:len(tempstack)-1]
		if val.dtype == 1 {
			i, err := strconv.ParseFloat(val.sval, 64)
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
	case "NOT":
		// Pop the top value from the stack and negate it
		val := tempstack[len(tempstack)-1]
		stack = tempstack[:len(tempstack)-1]
		if val.dtype != 2 {
			fmt.Println("Cannot negate non-bool")
			os.Exit(1)
		}
		tempstack = append(tempstack, stackVal{dtype: 2, bval: !val.bval})
	case "AND":
		// Pop the top two values from the stack and AND them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]
		if val1.dtype != 2 || val2.dtype != 2 {
			fmt.Println("Cannot AND non-bools")
			os.Exit(1)
		}
		tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.bval && val2.bval})
	case "OR":
		// Pop the top two values from the stack and OR them
		val1 := tempstack[len(tempstack)-1]
		val2 := tempstack[len(tempstack)-2]
		tempstack = tempstack[:len(tempstack)-2]
		if val1.dtype != 2 || val2.dtype != 2 {
			fmt.Println("Cannot OR non-bools")
			os.Exit(1)
		}
		tempstack = append(tempstack, stackVal{dtype: 2, bval: val1.bval || val2.bval})
	case "DELAYST":
		// Delay a certain amount of miliseconds
		val := tempstack[len(tempstack)-1]
		tempstack = tempstack[:len(tempstack)-1]
		if val.dtype != 0 {
			fmt.Println("Cannot delay non-int")
			os.Exit(1)
		}
		time.Sleep(time.Duration(val.val) * time.Millisecond)
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
				val := tempstack[len(tempstack)-1]
				tempstack = tempstack[:len(tempstack)-1]
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
						tempstack = append(tempstack, modu)
					}
				}
			}
		}
	case "CLEAR":
		// Clear stack
		tempstack = make([]stackVal, 0)
	case "MAKEARRAY":
		// Make an array
		var s stackVal = stackVal{dtype: 4, list: tempstack}
		tempstack = make([]stackVal, 0)
		tempstack = append(tempstack, s)
	case "SPLIT":
		// Split a string
		if len(parts) > 1 {
			if len(tempstack) > 0 {
				if tempstack[len(tempstack)-1].dtype == 1 {
					split := strings.Split(tempstack[len(tempstack)-1].sval, parts[1])
					tempstack = tempstack[:len(tempstack)-1]
					var s stackVal = stackVal{dtype: 4, list: make([]stackVal, 0)}
					for _, v := range split {
						s.list = append(s.list, stackVal{dtype: 1, sval: v})
					}
					tempstack = append(tempstack, s)
				} else {
					fmt.Println("Cannot split non-string")
					os.Exit(1)
				}
			} else {
				fmt.Println("tempstack is empty")
				os.Exit(1)
			}
		} else {
			fmt.Println("Invalid split")
			os.Exit(1)
		}
	case "JOIN":
		// Join a string
		if len(tempstack) > 0 {
			if tempstack[len(tempstack)-1].dtype == 4 {
				join := ""
				for _, v := range tempstack[len(tempstack)-1].list {
					if v.dtype == 1 {
						join += v.sval
						/*
							if index != len(tempstack[len(tempstack)-1].list)-1 {
								if !strings.HasSuffix(parts, "JOIN") {
									join += parts[1]
								}
							}*/
					} else {
						fmt.Println("Cannot join non-string")
						os.Exit(1)
					}
				}
				tempstack = tempstack[:len(tempstack)-1]
				var s stackVal = stackVal{dtype: 1, sval: join}
				tempstack = append(tempstack, s)
			} else {
				fmt.Println("Cannot join non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("tempstack is empty")
			os.Exit(1)
		}
	case "APPEND":
		// Append to an array
		if len(tempstack) > 1 {
			if tempstack[len(tempstack)-1].dtype == 4 {
				if tempstack[len(tempstack)-2].dtype == 4 {
					tempstack[len(tempstack)-2].list = append(tempstack[len(tempstack)-2].list, tempstack[len(tempstack)-1].list...)
					tempstack = tempstack[:len(tempstack)-1]
				} else {
					tempstack[len(tempstack)-2].list = append(tempstack[len(tempstack)-2].list, tempstack[len(tempstack)-1])
					tempstack = tempstack[:len(tempstack)-1]
				}
			} else {
				fmt.Println("Cannot append non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("tempstack is empty")
			os.Exit(1)
		}
	case "LEN":
		// Get the length of an array
		if len(tempstack) > 0 {
			if tempstack[len(tempstack)-1].dtype == 4 {
				tempstack = append(tempstack, stackVal{dtype: 0, val: float64(len(tempstack[len(tempstack)-1].list))})
			} else {
				fmt.Println("Cannot get length of non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("tempstack is empty")
			os.Exit(1)
		}
	case "REMOVE":
		// Remove an item from an array
		if len(tempstack) > 1 {
			if tempstack[len(tempstack)-2].dtype == 4 {
				if tempstack[len(tempstack)-1].dtype == 0 {
					if int(tempstack[len(tempstack)-1].val) < len(tempstack[len(tempstack)-2].list) {
						var i int = int(tempstack[len(tempstack)-1].val)
						var s stackVal = stackVal{dtype: 4, list: make([]stackVal, 0)}
						s.list = append(tempstack[len(tempstack)-2].list[:i], tempstack[len(tempstack)-2].list[i+1:]...)
						tempstack = tempstack[:len(tempstack)-2]
						tempstack = append(tempstack, s)
					} else {
						fmt.Println("Index out of range")
						os.Exit(1)
					}
				} else {
					fmt.Println("Cannot remove non-integer")
					os.Exit(1)
				}
			} else {
				fmt.Println("Cannot remove from non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("tempstack is empty")
			os.Exit(1)
		}
	case "RANDINT":
		// Generate a random integer
		if len(parts) > 2 {
			min, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid minimum")
				os.Exit(1)
			}
			max, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid maximum")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 0, val: float64(rand.Intn(max-min) + min)})
		} else {
			fmt.Println("Missing minimum and maximum")
			os.Exit(1)
		}
	case "RANDFLOAT":
		// Generate a random float
		if len(parts) > 2 {
			min, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				fmt.Println("Invalid minimum")
				os.Exit(1)
			}
			max, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				fmt.Println("Invalid maximum")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 0, sval: strconv.FormatFloat(rand.Float64()*(max-min)+min, 'f', -1, 64)})
		} else {
			fmt.Println("Missing minimum and maximum")
			os.Exit(1)
		}
	case "COMM":
		// Comment
	case "MCOMM":
		// Multi-line comment
		comment = true
	case "/*":
		// Multi-line comment
		comment = true
	case "//":
		// Comment
	default:
		var cont = false
		for s, v := range keymods {
			if op == s {
				for i, c := range v.cases {
					if !(len(parts) > 1) {
						fmt.Println("Invalid operation")
						os.Exit(1)
					}
					if i == parts[1] {
						var args = parts[2:]
						for n, a := range args {
							num, err := strconv.ParseFloat(a, 64)
							var s stackVal = stackVal{}
							if err != nil {
								if a == "true" {
									s = stackVal{dtype: 2, bval: true}
								} else if a == "false" {
									s = stackVal{dtype: 2, bval: false}
								} else if strings.HasPrefix(a, "\"") && strings.HasSuffix(a, "\"") {
									s = stackVal{dtype: 1, sval: a[1 : len(a)-1]}
								} else if strings.HasPrefix(a, "'") && strings.HasSuffix(a, "'") {
									s = stackVal{dtype: 1, sval: a[1 : len(a)-1]}
								} else {
									fmt.Println("Invalid argument")
									os.Exit(1)
								}
							} else {
								s = stackVal{dtype: 0, val: num}
							}

							tempsymbols["term"+strconv.Itoa(n)] = s
						}
						for _, o := range c {
							var line string = o.(string)
							runs(line)
						}
						cont = true
					}
				}
			}
		}
		tempsymbols = make(map[string]stackVal)

		if !cont {
			fmt.Printf("Invalid operation: %s\n", op)
			os.Exit(1)
		}
	}
	for name := range module.extvars {
		module.extvars[name] = module.symbols[name]
	}
	modules[modname] = module

}

func run(code string) {
	if code == "\n" || code == "" {
		return
	}
	// Split the line into parts
	code = strings.ReplaceAll(code, "	", "")
	code = strings.ReplaceAll(code, "    ", "")

	// Split the line into parts
	parts := strings.Split(code, " ")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
		parts[i] = strings.ReplaceAll(parts[i], "	", "")
	}

	// Determine the operation
	op := parts[0]

	// Execute the operation
	switch op {
	case "PUSH":
		var s stackVal = stackVal{}

		// Push the value onto the stack
		val, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			if strings.HasPrefix(strings.Join(parts[1:], " "), "\"") && strings.HasSuffix(strings.Join(parts[1:], " "), "\"") {
				s.sval = strings.Join(parts[1:], " ")[1 : len(strings.Join(parts[1:], " "))-1]
				s.dtype = 1
				stack = append(stack, s)
			} else if strings.HasPrefix(strings.Join(parts[1:], " "), "'") && strings.HasSuffix(strings.Join(parts[1:], " "), "'") {
				s.sval = strings.Join(parts[1:], " ")[1 : len(strings.Join(parts[1:], " "))-1]
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

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val + val2.val, dtype: 0})
	case "SUB":
		// Pop the top two values from the stack and subtract them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val - val2.val, dtype: 0})
	case "MULT":
		// Pop the top two values from the stack and multiply them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val * val2.val, dtype: 0})
	case "DIV":
		// Pop the top two values from the stack and divide them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
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
		if len(parts) > 2 {
			if val.dtype != 4 {
				fmt.Println("Cannot index non-array")
				os.Exit(1)
			}
			index, err := strconv.Atoi(parts[2])
			if err != nil {
				valt, ok := symbols[parts[2]]
				if !ok {
					fmt.Printf("Undefined symbol: %s\n", parts[1])
					os.Exit(1)
				} else if valt.dtype != 0 {
					fmt.Println("Index of array must be number")
					os.Exit(1)
				} else {
					index = int(valt.val)
				}
			}
			if index >= len(val.list) {
				fmt.Println("Index out of bounds")
				os.Exit(1)
			} else {
				val = val.list[index]
			}
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
		} else if val.dtype == 0 {
			fmt.Println(val.val)
		} else {
			fmt.Println("Cannot print element")
		}
	case "STR":
		var s stackVal = stackVal{}
		s.dtype = 1
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype == 0 {
			s.sval = strconv.FormatFloat(val.val, 'f', -1, 64)
		} else if val.dtype == 2 {
			s.sval = strconv.FormatBool(val.bval)
		} else {
			s.sval = val.sval
		}
		stack = append(stack, s)
	case "FLOAT":
		var s stackVal = stackVal{}
		s.dtype = 0
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype == 1 {
			i, err := strconv.ParseFloat(val.sval, 64)
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
	case "NOT":
		// Pop the top value from the stack and negate it
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype != 2 {
			fmt.Println("Cannot negate non-bool")
			os.Exit(1)
		}
		stack = append(stack, stackVal{dtype: 2, bval: !val.bval})
	case "AND":
		// Pop the top two values from the stack and AND them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != 2 || val2.dtype != 2 {
			fmt.Println("Cannot AND non-bools")
			os.Exit(1)
		}
		stack = append(stack, stackVal{dtype: 2, bval: val1.bval && val2.bval})
	case "OR":
		// Pop the top two values from the stack and OR them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != 2 || val2.dtype != 2 {
			fmt.Println("Cannot OR non-bools")
			os.Exit(1)
		}
		stack = append(stack, stackVal{dtype: 2, bval: val1.bval || val2.bval})
	case "DELAYST":
		// Delay a certain amount of miliseconds
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype != 0 {
			fmt.Println("Cannot delay non-int")
			os.Exit(1)
		}
		time.Sleep(time.Duration(val.val) * time.Millisecond)
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
	case "CLEAR":
		// Clear stack
		stack = make([]stackVal, 0)
	case "MAKEARRAY":
		// Make an array
		var s stackVal = stackVal{dtype: 4, list: stack}
		stack = make([]stackVal, 0)
		stack = append(stack, s)
	case "SPLIT":
		// Split a string
		if len(parts) > 1 {
			if len(stack) > 0 {
				if stack[len(stack)-1].dtype == 1 {
					split := strings.Split(stack[len(stack)-1].sval, parts[1])
					stack = stack[:len(stack)-1]
					var s stackVal = stackVal{dtype: 4, list: make([]stackVal, 0)}
					for _, v := range split {
						s.list = append(s.list, stackVal{dtype: 1, sval: v})
					}
					stack = append(stack, s)
				} else {
					fmt.Println("Cannot split non-string")
					os.Exit(1)
				}
			} else {
				fmt.Println("Stack is empty")
				os.Exit(1)
			}
		} else {
			fmt.Println("Invalid split")
			os.Exit(1)
		}
	case "JOIN":
		// Join a string
		if len(stack) > 0 {
			if stack[len(stack)-1].dtype == 4 {
				join := ""
				for _, v := range stack[len(stack)-1].list {
					if v.dtype == 1 {
						join += v.sval
						/*
							if index != len(stack[len(stack)-1].list)-1 {
								if !strings.HasSuffix(parts, "JOIN") {
									join += parts[1]
								}
							}*/
					} else {
						fmt.Println("Cannot join non-string")
						os.Exit(1)
					}
				}
				stack = stack[:len(stack)-1]
				var s stackVal = stackVal{dtype: 1, sval: join}
				stack = append(stack, s)
			} else {
				fmt.Println("Cannot join non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("Stack is empty")
			os.Exit(1)
		}
	case "APPEND":
		// Append to an array
		if len(stack) > 1 {
			if stack[len(stack)-1].dtype == 4 {
				if stack[len(stack)-2].dtype == 4 {
					stack[len(stack)-2].list = append(stack[len(stack)-2].list, stack[len(stack)-1].list...)
					stack = stack[:len(stack)-1]
				} else {
					stack[len(stack)-2].list = append(stack[len(stack)-2].list, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
			} else {
				fmt.Println("Cannot append non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("Stack is empty")
			os.Exit(1)
		}
	case "LEN":
		// Get the length of an array
		if len(stack) > 0 {
			if stack[len(stack)-1].dtype == 4 {
				stack = append(stack, stackVal{dtype: 0, val: float64(len(stack[len(stack)-1].list))})
			} else {
				fmt.Println("Cannot get length of non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("Stack is empty")
			os.Exit(1)
		}
	case "REMOVE":
		// Remove an item from an array
		if len(stack) > 1 {
			if stack[len(stack)-2].dtype == 4 {
				if stack[len(stack)-1].dtype == 0 {
					if int(stack[len(stack)-1].val) < len(stack[len(stack)-2].list) {
						var i int = int(stack[len(stack)-1].val)
						var s stackVal = stackVal{dtype: 4, list: make([]stackVal, 0)}
						s.list = append(stack[len(stack)-2].list[:i], stack[len(stack)-2].list[i+1:]...)
						stack = stack[:len(stack)-2]
						stack = append(stack, s)
					} else {
						fmt.Println("Index out of range")
						os.Exit(1)
					}
				} else {
					fmt.Println("Cannot remove non-integer")
					os.Exit(1)
				}
			} else {
				fmt.Println("Cannot remove from non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("Stack is empty")
			os.Exit(1)
		}
	case "RANDINT":
		// Generate a random integer
		if len(parts) > 1 {
			min, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid minimum")
				os.Exit(1)
			}
			max, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid maximum")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 0, val: float64(rand.Intn(max-min) + min)})
		} else {
			fmt.Println("Missing minimum and maximum")
			os.Exit(1)
		}
	case "RANDFLOAT":
		// Generate a random float
		if len(parts) > 1 {
			min, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				fmt.Println("Invalid minimum")
				os.Exit(1)
			}
			max, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				fmt.Println("Invalid maximum")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 0, sval: strconv.FormatFloat(rand.Float64()*(max-min)+min, 'f', -1, 64)})
		} else {
			fmt.Println("Missing minimum and maximum")
			os.Exit(1)
		}
	case "COMM":
		// Comment
	case "MCOMM":
		// Multi-line comment
		comment = true
	case "/*":
		// Multi-line comment
		comment = true
	case "//":
		// Comment
	default:
		var cont = false
		for s, v := range keymods {
			if op == s {
				for i, c := range v.cases {
					if !(len(parts) > 1) {
						fmt.Println("Invalid operation")
						os.Exit(1)
					}
					if i == parts[1] {
						var args = parts[2:]
						for n, a := range args {
							num, err := strconv.ParseFloat(a, 64)
							var s stackVal = stackVal{}
							if err != nil {
								if a == "true" {
									s = stackVal{dtype: 2, bval: true}
								} else if a == "false" {
									s = stackVal{dtype: 2, bval: false}
								} else if strings.HasPrefix(a, "\"") && strings.HasSuffix(a, "\"") {
									s = stackVal{dtype: 1, sval: a[1 : len(a)-1]}
								} else if strings.HasPrefix(a, "'") && strings.HasSuffix(a, "'") {
									s = stackVal{dtype: 1, sval: a[1 : len(a)-1]}
								} else {
									fmt.Println("Invalid argument")
									os.Exit(1)
								}
							} else {
								s = stackVal{dtype: 0, val: num}
							}

							tempsymbols["term"+strconv.Itoa(n)] = s
						}
						for _, o := range c {
							var line string = o.(string)
							runs(line)
						}
						cont = true
					}
				}
			}
		}
		tempsymbols = make(map[string]stackVal)

		if !cont {
			fmt.Printf("Invalid operation: %s\n", op)
			os.Exit(1)
		}
	}
}

func runs(code string) {
	if code == "\n" || code == "" {
		return
	}
	// Split the line into parts
	code = strings.ReplaceAll(code, "	", "")
	code = strings.ReplaceAll(code, "    ", "")

	// Split the line into parts
	parts := strings.Split(code, " ")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
		parts[i] = strings.ReplaceAll(parts[i], "	", "")
	}

	// Determine the operation
	op := parts[0]

	// Execute the operation
	switch op {
	case "PUSH":
		var s stackVal = stackVal{}

		// Push the value onto the stack
		val, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			if strings.HasPrefix(strings.Join(parts[1:], " "), "\"") && strings.HasSuffix(strings.Join(parts[1:], " "), "\"") {
				s.sval = strings.Join(parts[1:], " ")[1 : len(strings.Join(parts[1:], " "))-1]
				s.dtype = 1
				stack = append(stack, s)
			} else if strings.HasPrefix(strings.Join(parts[1:], " "), "'") && strings.HasSuffix(strings.Join(parts[1:], " "), "'") {
				s.sval = strings.Join(parts[1:], " ")[1 : len(strings.Join(parts[1:], " "))-1]
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

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val + val2.val, dtype: 0})
	case "SUB":
		// Pop the top two values from the stack and subtract them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val - val2.val, dtype: 0})
	case "MULT":
		// Pop the top two values from the stack and multiply them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
			os.Exit(1)
		}

		stack = append(stack, stackVal{val: val1.val * val2.val, dtype: 0})
	case "DIV":
		// Pop the top two values from the stack and divide them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		if !(val1.dtype == 0 || val2.dtype == 0) {
			fmt.Println("Cannot operate non-numbers")
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
		if tempsymbols[parts[1]].dtype != 1 {
			fmt.Println("Symbol is not a string")
			os.Exit(1)
		}
		symbols[tempsymbols[parts[1]].sval] = val
	case "LOADARG":
		val, ok := tempsymbols[parts[1]]
		if !ok {
			fmt.Printf("Undefined symbol: %s\n", parts[1])
			os.Exit(1)
		}
		stack = append(stack, val)
	case "LOAD":
		// Load the value from the symbol table and push it onto the stack
		if tempsymbols[parts[1]].dtype != 1 {
			fmt.Println("Symbol is not a string")
			os.Exit(1)
		}

		val, ok := symbols[tempsymbols[parts[1]].sval]
		if !ok {
			fmt.Printf("Undefined symbol: %s\n", parts[1])
			os.Exit(1)
		}
		if len(parts) > 2 {
			if val.dtype != 4 {
				fmt.Println("Cannot index non-array")
				os.Exit(1)
			}
			index, err := strconv.Atoi(parts[2])
			if err != nil {
				valt, ok := symbols[tempsymbols[parts[2]].sval]
				if !ok {
					fmt.Printf("Undefined symbol: %s\n", parts[1])
					os.Exit(1)
				} else if valt.dtype != 0 {
					fmt.Println("Index of array must be number")
					os.Exit(1)
				} else {
					index = int(valt.val)
				}
			}
			if index >= len(val.list) {
				fmt.Println("Index out of bounds")
				os.Exit(1)
			} else {
				val = val.list[index]
			}
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
		} else if val.dtype == 0 {
			fmt.Println(val.val)
		} else {
			fmt.Println("Cannot print element")
		}
	case "STR":
		var s stackVal = stackVal{}
		s.dtype = 1
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype == 0 {
			s.sval = strconv.FormatFloat(val.val, 'f', -1, 64)
		} else if val.dtype == 2 {
			s.sval = strconv.FormatBool(val.bval)
		} else {
			s.sval = val.sval
		}
		stack = append(stack, s)
	case "FLOAT":
		var s stackVal = stackVal{}
		s.dtype = 0
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype == 1 {
			i, err := strconv.ParseFloat(val.sval, 64)
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
	case "NOT":
		// Pop the top value from the stack and negate it
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype != 2 {
			fmt.Println("Cannot negate non-bool")
			os.Exit(1)
		}
		stack = append(stack, stackVal{dtype: 2, bval: !val.bval})
	case "AND":
		// Pop the top two values from the stack and AND them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != 2 || val2.dtype != 2 {
			fmt.Println("Cannot AND non-bools")
			os.Exit(1)
		}
		stack = append(stack, stackVal{dtype: 2, bval: val1.bval && val2.bval})
	case "OR":
		// Pop the top two values from the stack and OR them
		val1 := stack[len(stack)-1]
		val2 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if val1.dtype != 2 || val2.dtype != 2 {
			fmt.Println("Cannot OR non-bools")
			os.Exit(1)
		}
		stack = append(stack, stackVal{dtype: 2, bval: val1.bval || val2.bval})
	case "DELAYST":
		// Delay a certain amount of miliseconds
		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if val.dtype != 0 {
			fmt.Println("Cannot delay non-int")
			os.Exit(1)
		}
		time.Sleep(time.Duration(val.val) * time.Millisecond)
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
	case "CLEAR":
		// Clear stack
		stack = make([]stackVal, 0)
	case "MAKEARRAY":
		// Make an array
		var s stackVal = stackVal{dtype: 4, list: stack}
		stack = make([]stackVal, 0)
		stack = append(stack, s)
	case "SPLIT":
		// Split a string
		if len(parts) > 1 {
			if len(stack) > 0 {
				if stack[len(stack)-1].dtype == 1 {
					split := strings.Split(stack[len(stack)-1].sval, parts[1])
					stack = stack[:len(stack)-1]
					var s stackVal = stackVal{dtype: 4, list: make([]stackVal, 0)}
					for _, v := range split {
						s.list = append(s.list, stackVal{dtype: 1, sval: v})
					}
					stack = append(stack, s)
				} else {
					fmt.Println("Cannot split non-string")
					os.Exit(1)
				}
			} else {
				fmt.Println("Stack is empty")
				os.Exit(1)
			}
		} else {
			fmt.Println("Invalid split")
			os.Exit(1)
		}
	case "JOIN":
		// Join a string
		if len(stack) > 0 {
			if stack[len(stack)-1].dtype == 4 {
				join := ""
				for _, v := range stack[len(stack)-1].list {
					if v.dtype == 1 {
						join += v.sval
						/*
							if index != len(stack[len(stack)-1].list)-1 {
								if !strings.HasSuffix(parts, "JOIN") {
									join += parts[1]
								}
							}*/
					} else {
						fmt.Println("Cannot join non-string")
						os.Exit(1)
					}
				}
				stack = stack[:len(stack)-1]
				var s stackVal = stackVal{dtype: 1, sval: join}
				stack = append(stack, s)
			} else {
				fmt.Println("Cannot join non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("Stack is empty")
			os.Exit(1)
		}
	case "APPEND":
		// Append to an array
		if len(stack) > 1 {
			if stack[len(stack)-1].dtype == 4 {
				if stack[len(stack)-2].dtype == 4 {
					stack[len(stack)-2].list = append(stack[len(stack)-2].list, stack[len(stack)-1].list...)
					stack = stack[:len(stack)-1]
				} else {
					stack[len(stack)-2].list = append(stack[len(stack)-2].list, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
			} else {
				fmt.Println("Cannot append non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("Stack is empty")
			os.Exit(1)
		}
	case "LEN":
		// Get the length of an array
		if len(stack) > 0 {
			if stack[len(stack)-1].dtype == 4 {
				stack = append(stack, stackVal{dtype: 0, val: float64(len(stack[len(stack)-1].list))})
			} else {
				fmt.Println("Cannot get length of non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("Stack is empty")
			os.Exit(1)
		}
	case "REMOVE":
		// Remove an item from an array
		if len(stack) > 1 {
			if stack[len(stack)-2].dtype == 4 {
				if stack[len(stack)-1].dtype == 0 {
					if int(stack[len(stack)-1].val) < len(stack[len(stack)-2].list) {
						var i int = int(stack[len(stack)-1].val)
						var s stackVal = stackVal{dtype: 4, list: make([]stackVal, 0)}
						s.list = append(stack[len(stack)-2].list[:i], stack[len(stack)-2].list[i+1:]...)
						stack = stack[:len(stack)-2]
						stack = append(stack, s)
					} else {
						fmt.Println("Index out of range")
						os.Exit(1)
					}
				} else {
					fmt.Println("Cannot remove non-integer")
					os.Exit(1)
				}
			} else {
				fmt.Println("Cannot remove from non-array")
				os.Exit(1)
			}
		} else {
			fmt.Println("Stack is empty")
			os.Exit(1)
		}
	case "RANDINT":
		// Generate a random integer
		if len(parts) > 1 {
			min, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid minimum")
				os.Exit(1)
			}
			max, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid maximum")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 0, val: float64(rand.Intn(max-min) + min)})
		} else {
			fmt.Println("Missing minimum and maximum")
			os.Exit(1)
		}
	case "RANDFLOAT":
		// Generate a random float
		if len(parts) > 1 {
			min, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				fmt.Println("Invalid minimum")
				os.Exit(1)
			}
			max, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				fmt.Println("Invalid maximum")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 0, sval: strconv.FormatFloat(rand.Float64()*(max-min)+min, 'f', -1, 64)})
		} else {
			fmt.Println("Missing minimum and maximum")
			os.Exit(1)
		}
	case "COMM":
		// Comment
	case "MCOMM":
		// Multi-line comment
		comment = true
	case "/*":
		// Multi-line comment
		comment = true
	case "//":
		// Comment
	default:
		var cont = false
		for s, v := range keymods {
			if op == s {
				for i, c := range v.cases {
					if !(len(parts) > 1) {
						fmt.Println("Invalid operation")
						os.Exit(1)
					}
					if i == parts[1] {
						var args = parts[2:]
						for n, a := range args {
							num, err := strconv.ParseFloat(a, 64)
							var s stackVal = stackVal{}
							if err != nil {
								if a == "true" {
									s = stackVal{dtype: 2, bval: true}
								} else if a == "false" {
									s = stackVal{dtype: 2, bval: false}
								} else if strings.HasPrefix(a, "\"") && strings.HasSuffix(a, "\"") {
									s = stackVal{dtype: 1, sval: a[1 : len(a)-1]}
								} else if strings.HasPrefix(a, "'") && strings.HasSuffix(a, "'") {
									s = stackVal{dtype: 1, sval: a[1 : len(a)-1]}
								} else {
									fmt.Println("Invalid argument")
									os.Exit(1)
								}
							} else {
								s = stackVal{dtype: 0, val: num}
							}

							tempsymbols["term"+strconv.Itoa(n)] = s
						}
						for _, o := range c {
							var line string = o.(string)
							runs(line)
						}
						cont = true
					}
				}
			}
		}
		tempsymbols = make(map[string]stackVal)

		if !cont {
			fmt.Printf("Invalid operation: %s\n", op)
			os.Exit(1)
		}

	}
}

type keymod struct {
	cases map[string][]interface{}
}

var keymods map[string]keymod = make(map[string]keymod)
var tempsymbols = make(map[string]stackVal)

func main() {
	rand.Seed(time.Now().UnixNano())

	defimports()

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
		if line == "" {
			continue
		}

		line = strings.ReplaceAll(line, "	", "")
		line = strings.ReplaceAll(line, "    ", "")

		// Split the line into parts
		parts := strings.Split(line, " ")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
			parts[i] = strings.ReplaceAll(parts[i], "	", "")
		}
		// Determine the operation
		op := parts[0]
		if comment {
			if op == "ENDCOMM" || op == "*/" {
				comment = false
			}
			continue
		}
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
		}

		// Execute the operation
		switch op {
		case "KEYPORT":
			// import .kr file with keywords
			if !(len(parts) == 2) {
				fmt.Println("Invalid keyword call")
				os.Exit(1)
			} else {
				j, err := os.Open(parts[1])

				if err != nil {
					fmt.Println("Invalid keyword file")
					os.Exit(1)
				}

				defer j.Close()

				byteValue, _ := ioutil.ReadAll(j)

				var result map[string]interface{}
				json.Unmarshal([]byte(byteValue), &result)
				var util keymod = keymod{cases: make(map[string][]interface{})}
				for _, v := range result["main"].([]interface{}) {
					util.cases[v.(map[string]interface{})["case"].(string)] = v.(map[string]interface{})["code"].([]interface{})
				}
				keymods[result["prefix"].(string)] = util

			}

		case "PUSH":
			var s stackVal = stackVal{}

			// Push the value onto the stack
			val, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				if strings.HasPrefix(strings.Join(parts[1:], " "), "\"") && strings.HasSuffix(strings.Join(parts[1:], " "), "\"") {
					s.sval = strings.Join(parts[1:], " ")[1 : len(strings.Join(parts[1:], " "))-1]
					s.dtype = 1
					stack = append(stack, s)
				} else if strings.HasPrefix(strings.Join(parts[1:], " "), "'") && strings.HasSuffix(strings.Join(parts[1:], " "), "'") {
					s.sval = strings.Join(parts[1:], " ")[1 : len(strings.Join(parts[1:], " "))-1]
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

			if !(val1.dtype == 0 || val2.dtype == 0) {
				fmt.Println("Cannot operate non-numbers")
				os.Exit(1)
			}

			stack = append(stack, stackVal{val: val1.val + val2.val, dtype: 0})
		case "SUB":
			// Pop the top two values from the stack and subtract them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if !(val1.dtype == 0 || val2.dtype == 0) {
				fmt.Println("Cannot operate non-numbers")
				os.Exit(1)
			}

			stack = append(stack, stackVal{val: val1.val - val2.val, dtype: 0})
		case "MULT":
			// Pop the top two values from the stack and multiply them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if !(val1.dtype == 0 || val2.dtype == 0) {
				fmt.Println("Cannot operate non-numbers")
				os.Exit(1)
			}

			stack = append(stack, stackVal{val: val1.val * val2.val, dtype: 0})
		case "DIV":
			// Pop the top two values from the stack and divide them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if !(val1.dtype == 0 || val2.dtype == 0) {
				fmt.Println("Cannot operate non-numbers")
				os.Exit(1)
			}

			if !(val1.dtype == 0 || val2.dtype == 0) {
				fmt.Println("Cannot divide non-numbers")
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
			if len(parts) > 2 {
				if val.dtype != 4 {
					fmt.Println("Cannot index non-array")
					os.Exit(1)
				}
				index, err := strconv.Atoi(parts[2])
				if err != nil {
					valt, ok := symbols[parts[2]]
					if !ok {
						fmt.Printf("Undefined symbol: %s\n", parts[1])
						os.Exit(1)
					} else if valt.dtype != 0 {
						fmt.Println("Index of array must be number")
						os.Exit(1)
					} else {
						index = int(valt.val)
					}
				}
				if index >= len(val.list) {
					fmt.Println("Index out of bounds")
					os.Exit(1)
				} else {
					val = val.list[index]
				}
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
			} else if val.dtype == 0 {
				fmt.Println(val.val)
			} else {
				fmt.Println("Cannot print element")
			}
		case "STR":
			var s stackVal = stackVal{}
			s.dtype = 1
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if val.dtype == 0 {
				s.sval = strconv.FormatFloat(val.val, 'f', -1, 64)
			} else if val.dtype == 2 {
				s.sval = strconv.FormatBool(val.bval)
			} else {
				s.sval = val.sval
			}
			stack = append(stack, s)
		case "FLOAT":
			var s stackVal = stackVal{}
			s.dtype = 0
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if val.dtype == 1 {
				i, err := strconv.ParseFloat(val.sval, 64)
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
		case "NOT":
			// Pop the top value from the stack and negate it
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if val.dtype != 2 {
				fmt.Println("Cannot negate non-bool")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 2, bval: !val.bval})
		case "AND":
			// Pop the top two values from the stack and AND them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if val1.dtype != 2 || val2.dtype != 2 {
				fmt.Println("Cannot AND non-bools")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 2, bval: val1.bval && val2.bval})
		case "OR":
			// Pop the top two values from the stack and OR them
			val1 := stack[len(stack)-1]
			val2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if val1.dtype != 2 || val2.dtype != 2 {
				fmt.Println("Cannot OR non-bools")
				os.Exit(1)
			}
			stack = append(stack, stackVal{dtype: 2, bval: val1.bval || val2.bval})
		case "DELAYST":
			// Delay a certain amount of miliseconds
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if val.dtype != 0 {
				fmt.Println("Cannot delay non-int")
				os.Exit(1)
			}
			time.Sleep(time.Duration(val.val) * time.Millisecond)
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
		case "CLEAR":
			// Clear stack
			stack = make([]stackVal, 0)
		case "MAKEARRAY":
			// Make an array
			var s stackVal = stackVal{dtype: 4, list: stack}
			stack = make([]stackVal, 0)
			stack = append(stack, s)
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
							val, err := strconv.ParseFloat(partsin[2], 64)
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
					case "EXARR":
						// Export an array
						if len(partsin) > 1 {
							var s stackVal = stackVal{dtype: 4, list: stack}
							m.extvars[partsin[1]] = s
						} else {
							fmt.Println("Invalid export")
							os.Exit(1)
						}
					case "COMM":
						// Comment
					case "MCOMM":
						// Multi-line comment
						comment = true
					case "/*":
						// Multi-line comment
						comment = true
					case "//":
						// Comment
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
		case "SPLIT":
			// Split a string
			if len(parts) > 1 {
				if len(stack) > 0 {
					if stack[len(stack)-1].dtype == 1 {
						split := strings.Split(stack[len(stack)-1].sval, parts[1])
						stack = stack[:len(stack)-1]
						var s stackVal = stackVal{dtype: 4, list: make([]stackVal, 0)}
						for _, v := range split {
							s.list = append(s.list, stackVal{dtype: 1, sval: v})
						}
						stack = append(stack, s)
					} else {
						fmt.Println("Cannot split non-string")
						os.Exit(1)
					}
				} else {
					fmt.Println("Stack is empty")
					os.Exit(1)
				}
			} else {
				fmt.Println("Invalid split")
				os.Exit(1)
			}
		case "JOIN":
			// Join a string
			if len(stack) > 0 {
				if stack[len(stack)-1].dtype == 4 {
					join := ""
					for _, v := range stack[len(stack)-1].list {
						if v.dtype == 1 {
							join += v.sval
							/*
								if index != len(stack[len(stack)-1].list)-1 {
									if !strings.HasSuffix(parts, "JOIN") {
										join += parts[1]
									}
								}*/
						} else {
							fmt.Println("Cannot join non-string")
							os.Exit(1)
						}
					}
					stack = stack[:len(stack)-1]
					var s stackVal = stackVal{dtype: 1, sval: join}
					stack = append(stack, s)
				} else {
					fmt.Println("Cannot join non-array")
					os.Exit(1)
				}
			} else {
				fmt.Println("Stack is empty")
				os.Exit(1)
			}
		case "APPEND":
			// Append to an array
			if len(stack) > 1 {
				if stack[len(stack)-1].dtype == 4 {
					if stack[len(stack)-2].dtype == 4 {
						stack[len(stack)-2].list = append(stack[len(stack)-2].list, stack[len(stack)-1].list...)
						stack = stack[:len(stack)-1]
					} else {
						stack[len(stack)-2].list = append(stack[len(stack)-2].list, stack[len(stack)-1])
						stack = stack[:len(stack)-1]
					}
				} else {
					fmt.Println("Cannot append non-array")
					os.Exit(1)
				}
			} else {
				fmt.Println("Stack is empty")
				os.Exit(1)
			}
		case "LEN":
			// Get the length of an array
			if len(stack) > 0 {
				if stack[len(stack)-1].dtype == 4 {
					stack = append(stack, stackVal{dtype: 0, val: float64(len(stack[len(stack)-1].list))})
				} else {
					fmt.Println("Cannot get length of non-array")
					os.Exit(1)
				}
			} else {
				fmt.Println("Stack is empty")
				os.Exit(1)
			}
		case "REMOVE":
			// Remove an item from an array
			if len(stack) > 1 {
				if stack[len(stack)-2].dtype == 4 {
					if stack[len(stack)-1].dtype == 0 {
						if int(stack[len(stack)-1].val) < len(stack[len(stack)-2].list) {
							var i int = int(stack[len(stack)-1].val)
							var s stackVal = stackVal{dtype: 4, list: make([]stackVal, 0)}
							s.list = append(stack[len(stack)-2].list[:i], stack[len(stack)-2].list[i+1:]...)
							stack = stack[:len(stack)-2]
							stack = append(stack, s)
						} else {
							fmt.Println("Index out of range")
							os.Exit(1)
						}
					} else {
						fmt.Println("Cannot remove non-integer")
						os.Exit(1)
					}
				} else {
					fmt.Println("Cannot remove from non-array")
					os.Exit(1)
				}
			} else {
				fmt.Println("Stack is empty")
				os.Exit(1)
			}
		case "COMM":
			// Comment
		case "MCOMM":
			// Multi-line comment
			comment = true
		case "/*":
			// Multi-line comment
			comment = true
		case "//":
			// Comment
		case "RANDINT":
			// Generate a random integer
			if len(parts) > 1 {
				min, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println("Invalid minimum")
					os.Exit(1)
				}
				max, err := strconv.Atoi(parts[2])
				if err != nil {
					fmt.Println("Invalid maximum")
					os.Exit(1)
				}
				stack = append(stack, stackVal{dtype: 0, val: float64(rand.Intn(max-min) + min)})
			} else {
				fmt.Println("Missing minimum and maximum")
				os.Exit(1)
			}
		case "RANDFLOAT":
			// Generate a random float
			if len(parts) > 1 {
				min, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					fmt.Println("Invalid minimum")
					os.Exit(1)
				}
				max, err := strconv.ParseFloat(parts[2], 64)
				if err != nil {
					fmt.Println("Invalid maximum")
					os.Exit(1)
				}
				stack = append(stack, stackVal{dtype: 0, sval: strconv.FormatFloat(rand.Float64()*(max-min)+min, 'f', -1, 64)})
			} else {
				fmt.Println("Missing minimum and maximum")
				os.Exit(1)
			}

		default:
			var cont = false
			for s, v := range keymods {
				if op == s {
					for i, c := range v.cases {
						if !(len(parts) > 1) {
							fmt.Println("Invalid operation")
							os.Exit(1)
						}
						if i == parts[1] {
							var args = parts[2:]
							for n, a := range args {
								num, err := strconv.ParseFloat(a, 64)
								var s stackVal = stackVal{}
								if err != nil {
									if a == "true" {
										s = stackVal{dtype: 2, bval: true}
									} else if a == "false" {
										s = stackVal{dtype: 2, bval: false}
									} else if strings.HasPrefix(a, "\"") && strings.HasSuffix(a, "\"") {
										s = stackVal{dtype: 1, sval: a[1 : len(a)-1]}
									} else if strings.HasPrefix(a, "'") && strings.HasSuffix(a, "'") {
										s = stackVal{dtype: 1, sval: a[1 : len(a)-1]}
									} else {
										fmt.Println("Invalid argument")
										os.Exit(1)
									}
								} else {
									s = stackVal{dtype: 0, val: num}
								}

								tempsymbols["term"+strconv.Itoa(n)] = s
							}
							for _, o := range c {
								var line string = o.(string)
								runs(line)
							}
							cont = true
						}
					}
				}
			}
			tempsymbols = make(map[string]stackVal)

			if cont {
				continue
			}
			fmt.Printf("Invalid operation: %s\n", op)
			os.Exit(1)

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
	}
}
