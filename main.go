package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	push_type = iota
	add_type
	sub_type
	mul_type
	div_type
	mod_type
	equal_type
	less_type
	greater_type
	and_type
	or_type
	invert_type
	drop_type
	dup_type
	swap_type
	over_type
	rot_type
	print_type
    printstr_type
	emit_type
	cr_type
    quote_type
	colon_type
	semi_type
	word_type
	if_type
    else_type
	then_type
    do_type
    loop_type
    i_type
    variable_type
    constant_type
    at_type
    bang_type
    plus_bang_type
    question_type
    cells_type
    allot_type
)

type Token struct {
	value    string
	category int
}

type Stack struct {
	stack      []int
	stack_size int
}

type Loop struct {
    current int
    target int
}

func handle_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (stack *Stack) push(value int) {
	stack.stack = append(stack.stack, value)
	stack.stack_size++
}

func (stack *Stack)pop() int {
	if stack.stack_size <= 0 {
		log.Fatal("error: stack underflow")
	}

	stack.stack_size--
	result := stack.stack[stack.stack_size]
	stack.stack = stack.stack[:stack.stack_size]
	return result
}

func (stack *Stack) peek() int {
	if stack.stack_size < 0 {
		log.Fatal("error: stack underflow in peek")
	}

	result := stack.stack[stack.stack_size-1]
	return result
}

func is_digit(val string) bool {
	if _, err := strconv.Atoi(val); err == nil {
		return true
	}
	return false
}

func lex_file(filename string) []Token {
	bytes, err := os.ReadFile(filename)

	handle_error(err)

	contents := strings.Fields(string(bytes))

	tokens := []Token{}

	for _, val := range contents {
		var token Token
		switch val {
		case "+":
			token.category = add_type
		case "-":
			token.category = sub_type
		case "*":
			token.category = mul_type
		case "/":
			token.category = div_type
		case "mod":
			token.category = mod_type
		case "=":
			token.category = equal_type
		case "<":
			token.category = less_type
		case ">":
			token.category = greater_type
		case "and":
			token.category = and_type
		case "or":
			token.category = or_type
		case "invert":
			token.category = invert_type
		case "drop":
			token.category = drop_type
		case "dup":
			token.category = dup_type
		case "swap":
			token.category = swap_type
		case "over":
			token.category = over_type
		case "rot":
			token.category = rot_type
		case ".":
			token.category = print_type
        case ".\"":
            token.category = printstr_type
        case "\"":
            token.category = quote_type
		case "emit":
			token.category = emit_type
		case "cr":
			token.category = cr_type
		case ":":
			token.category = colon_type
		case ";":
			token.category = semi_type
		case "if":
			token.category = if_type
		case "else":
			token.category = else_type
		case "then":
			token.category = then_type
        case "do":
            token.category = do_type
        case "loop":
            token.category = loop_type
        case "i":
            token.category = i_type
        case "variable":
            token.category = variable_type
        case "constant":
            token.category = constant_type
        case "@":
            token.category = at_type
        case "!":
            token.category = bang_type
        case "+!":
            token.category = plus_bang_type
        case "?":
            token.category = question_type
        case "allot":
            token.category = allot_type
        case "cells":
            token.category = cells_type
		default:
			if is_digit(val) {
				token.category = push_type
				token.value = val
			} else {
				token.category = word_type
				token.value = val
			}
		}
        token.value = val
		tokens = append(tokens, token)
	}

	return tokens
}

type Gorth struct {
    stack Stack
    words map[string][]Token
    is_word bool
    current_word string
    variables map[string]int
    constants map[string]int
    memory []int
}

func (gorth *Gorth) interpret_tokens(tokens []Token) {
    in_word := false
	var word_tokens []Token
    var string_buffer string
    var loop Loop
    
	for i := 0; i < len(tokens); i++ {
		val := tokens[i]
		if in_word {
			if val.category == semi_type {
				gorth.words[gorth.current_word] = word_tokens
                word_tokens = nil
				in_word = false
			} else if val.category == word_type {
                    word_tokens = append(word_tokens, val)
            } else if val.category == printstr_type {
                for index := 1; val.category != quote_type; i += index {
                    val = tokens[i]
                    word_tokens = append(word_tokens, val)
                }
                i--
            } else {
				word_tokens = append(word_tokens, val)
			}
		} else {
			switch val.category {
			case push_type:
				value, err := strconv.Atoi(val.value)
				handle_error(err)
				gorth.stack.push(value)
			case add_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				gorth.stack.push(b + a)
			case sub_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				gorth.stack.push(b - a)
			case mul_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				gorth.stack.push(b * a)
			case div_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				gorth.stack.push(b / a)
			case mod_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				gorth.stack.push(b % a)
			case equal_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				if a == b {
					gorth.stack.push(-1)
				} else {
					gorth.stack.push(0)
				}
			case less_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				if a > b {
					gorth.stack.push(-1)
				} else {
					gorth.stack.push(0)
				}
			case greater_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				if a < b {
					gorth.stack.push(-1)
				} else {
					gorth.stack.push(0)
				}
			case and_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				if a != 0 && b != 0 {
					gorth.stack.push(-1)
				} else {
					gorth.stack.push(0)
				}
			case or_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				if a != 0 || b != 0 {
					gorth.stack.push(-1)
				} else {
					gorth.stack.push(0)
				}
			case invert_type:
				a := gorth.stack.pop()
				gorth.stack.push(-a - 1)
			case drop_type:
				gorth.stack.pop()
			case dup_type:
				a := gorth.stack.peek()
				gorth.stack.push(a)
			case swap_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				gorth.stack.push(a)
				gorth.stack.push(b)
			case over_type:
				a := gorth.stack.stack[gorth.stack.stack_size-2]
				gorth.stack.push(a)
			case rot_type:
				a := gorth.stack.pop()
				b := gorth.stack.pop()
				c := gorth.stack.pop()
				gorth.stack.push(b)
				gorth.stack.push(a)
				gorth.stack.push(c)
			case print_type:
				a := gorth.stack.pop()
				fmt.Print(a)
            case printstr_type:
                for index := 1; val.category != quote_type; index++ {
                    val = tokens[i]
                    if val.category != printstr_type && val.category != quote_type {
                        if index == 2 {
                            string_buffer += val.value
                        } else {
                            string_buffer += " " + val.value
                        }
                    }
                    i++
                }
                i--
                fmt.Print(string_buffer)
                string_buffer = ""
            case quote_type:
                log.Fatal("error in quote")
			case emit_type:
				a := gorth.stack.pop()
				fmt.Print(string(a))
			case cr_type:
				fmt.Println()
			case colon_type:
				in_word = true
                i++
                gorth.current_word = tokens[i].value 
			case semi_type:
				in_word = false
			case word_type:
                ptr, ok := gorth.variables[val.value]
                constant, cok := gorth.constants[val.value]
                if cok {
                    gorth.stack.push(constant)
                } else if ok {
                    variable := gorth.memory[ptr] 
                    if i + 1 >= len(tokens) {
                        gorth.stack.push(ptr)
                    } else {
                        if tokens[i + 1].category == bang_type {
                            a := gorth.stack.pop()
                            gorth.memory[ptr] = a
                        } else if tokens[i + 1].category == plus_bang_type {
                            a := gorth.stack.pop()
                            gorth.memory[ptr] = variable + a
                        } else {
                            gorth.stack.push(ptr)
                        }
                    }
                } else if in_word {
					gorth.current_word = val.value
				} else {
                    gorth.is_word = true
					gorth.interpret_tokens(gorth.words[val.value])
				}
            case at_type:
                a := gorth.stack.pop()
                gorth.stack.push(gorth.memory[a])
            case question_type:
                a := gorth.stack.pop()
                fmt.Print(gorth.memory[a])
			case if_type:
				if gorth.is_word {
					a := gorth.stack.pop()
					if a == 0 {
						for index := 0; val.category != then_type && val.category != else_type; index++ {
							val = tokens[i]
							i++
						}
                        i--
					}
				} else {
					log.Fatal("error if")
				}
            case else_type:
                if gorth.is_word {
					a := gorth.stack.pop()
                    if a != 0 {
                        for index := 0; val.category != then_type; index++ {
                            val = tokens[i]
                            i++
                        }
                    }
                } else {
                    log.Fatal("should not else")
                }
			case then_type:
            case do_type:
                if gorth.is_word {
                    starting := gorth.stack.pop()
                    target := gorth.stack.pop()
                    loop.current = starting
                    loop.target = target
                } else {
                    log.Fatal("should not do")
                }
            case loop_type:
                if gorth.is_word {
                    if loop.current < loop.target - 1 {
                        for index := 1; val.category != do_type; i -= index {
                            val = tokens[i] 
                        }
                        loop.current += 1
                    }
                    i++
                } else {
                    log.Fatal("should not do")
                }
            case i_type:
                if gorth.is_word {
                    gorth.stack.push(loop.current)
                } else {
                    log.Fatal("should not do")
                }
            case variable_type:
                i++
                val = tokens[i]
                _, cok := gorth.constants[val.value]
                if cok {
                    delete(gorth.constants, val.value)
                }
                gorth.variables[val.value] = len(gorth.memory) 
                gorth.memory = append(gorth.memory, 0)
            case constant_type:
                i++
                val = tokens[i]
                _, ok := gorth.variables[val.value]
                if ok {
                    delete(gorth.variables, val.value)
                }
                a := gorth.stack.pop()
                gorth.constants[val.value] = a
			}
		}
	}

}

func main() {
	if len(os.Args) < 2 {
		//log.Fatal(fmt.Sprintf("usage: %s <file_name.forth>", os.Args[0]))
	}

	//filename := os.Args[1]

    filename := "test.forth"
	tokens := lex_file(filename)

    var gorth Gorth
    var stack Stack

	words := make(map[string][]Token)
	variables := make(map[string]int)
	constants := make(map[string]int)
    var memory []int
    gorth.stack = stack
    gorth.words = words
    gorth.variables = variables
    gorth.constants = constants
    gorth.is_word = false
    gorth.memory = memory

	gorth.interpret_tokens(tokens)
}
