package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
    "unsafe"
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

func (stack *Stack) pop() int {
	if stack.stack_size == 0 {
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

func interpret_tokens(tokens []Token, stack *Stack, words map[string][]Token, is_word bool, 
    current_word string, variables map[string]int, constants map[string]int) {
	in_word := false
	var word_tokens []Token
    var string_buffer string
    var loop Loop
    
	for i := 0; i < len(tokens); i++ {
		val := tokens[i]
        if is_word {
        }
		if in_word {
			if val.category == semi_type {
				words[current_word] = word_tokens
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
				stack.push(value)
			case add_type:
				a := stack.pop()
				b := stack.pop()
				stack.push(b + a)
			case sub_type:
				a := stack.pop()
				b := stack.pop()
				stack.push(b - a)
			case mul_type:
				a := stack.pop()
				b := stack.pop()
				stack.push(b * a)
			case div_type:
				a := stack.pop()
				b := stack.pop()
				stack.push(b / a)
			case mod_type:
				a := stack.pop()
				b := stack.pop()
				stack.push(b % a)
			case equal_type:
				a := stack.pop()
				b := stack.pop()
				if a == b {
					stack.push(-1)
				} else {
					stack.push(0)
				}
			case less_type:
				a := stack.pop()
				b := stack.pop()
				if a > b {
					stack.push(-1)
				} else {
					stack.push(0)
				}
			case greater_type:
				a := stack.pop()
				b := stack.pop()
				if a < b {
					stack.push(-1)
				} else {
					stack.push(0)
				}
			case and_type:
				a := stack.pop()
				b := stack.pop()
				if a != 0 && b != 0 {
					stack.push(-1)
				} else {
					stack.push(0)
				}
			case or_type:
				a := stack.pop()
				b := stack.pop()
				if a != 0 || b != 0 {
					stack.push(-1)
				} else {
					stack.push(0)
				}
			case invert_type:
				a := stack.pop()
				stack.push(-a - 1)
			case drop_type:
				stack.pop()
			case dup_type:
				a := stack.peek()
				stack.push(a)
			case swap_type:
				a := stack.pop()
				b := stack.pop()
				stack.push(a)
				stack.push(b)
			case over_type:
				a := stack.stack[stack.stack_size-2]
				stack.push(a)
			case rot_type:
				a := stack.pop()
				b := stack.pop()
				c := stack.pop()
				stack.push(b)
				stack.push(a)
				stack.push(c)
			case print_type:
				a := stack.pop()
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
				a := stack.pop()
				fmt.Print(string(a))
			case cr_type:
				fmt.Println()
			case colon_type:
				in_word = true
                i++
                current_word = tokens[i].value 
			case semi_type:
				in_word = false
			case word_type:
                variable, ok := variables[val.value]
                constant, cok := constants[val.value]
                if cok {
                    stack.push(constant)
                } else if ok {
                    ptr := (int)(uintptr(unsafe.Pointer(&variable)))
                    if tokens[i + 1].category == bang_type {
                        a := stack.pop()
                        variables[val.value] = a
                    } else if tokens[i + 1].category == plus_bang_type {
                        a := stack.pop()
                        variables[val.value] = variable + a
                    } else {
                        stack.push(ptr)
                    }
                } else if in_word {
					current_word = val.value
				} else {
					interpret_tokens(words[val.value], stack, words, true, current_word, variables, constants)
				}
            case at_type:
                a := stack.pop()
                stack.push(*(*int)(unsafe.Pointer(uintptr(a))))
            case question_type:
                a := stack.pop()
                fmt.Print(*(*int)(unsafe.Pointer(uintptr(a))))
			case if_type:
				if is_word {
					a := stack.pop()
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
                if is_word {
					a := stack.pop()
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
                if is_word {
                    starting := stack.pop()
                    target := stack.pop()
                    loop.current = starting
                    loop.target = target
                } else {
                    log.Fatal("should not do")
                }
            case loop_type:
                if is_word {
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
                if is_word {
                    stack.push(loop.current)
                } else {
                    log.Fatal("should not do")
                }
            case variable_type:
                i++
                val = tokens[i]
                _, cok := constants[val.value]
                if cok {
                    delete(constants, val.value)
                }
                variables[val.value] = 0
            case constant_type:
                i++
                val = tokens[i]
                _, ok := variables[val.value]
                if ok {
                    delete(variables, val.value)
                }
                a := stack.pop()
                constants[val.value] = a
			}
		}
	}

}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Sprintf("usage: %s <file_name.forth>", os.Args[0]))
	}

	filename := os.Args[1]

	tokens := lex_file(filename)

	var stack Stack
    var current_word string
	words := make(map[string][]Token)
	variables := make(map[string]int)
	constants := make(map[string]int)

	interpret_tokens(tokens, &stack, words, false, current_word, variables, constants)
}
