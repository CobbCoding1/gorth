package main

import (
    "fmt"
    "strings"
    "strconv"
    "os"
    "log"
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
    emit_type
    cr_type
    colon_type
    semi_type
    word_type
    if_type
    then_type
)

type Token struct {
    value string
    category int
}

type Stack struct {
    stack []int
    stack_size int
}

func handle_error(err error){
    if err != nil {
        log.Fatal(err)
    }
}

func (stack *Stack)push(value int){
    stack.stack = append(stack.stack, value) 
    stack.stack_size++
}

func (stack *Stack)pop()(int){
    if(stack.stack_size == 0){
        log.Fatal("error: stack underflow")
    }

    stack.stack_size--
    result := stack.stack[stack.stack_size]
    stack.stack = stack.stack[:stack.stack_size]
    return result
}

func (stack *Stack)peek()(int){
    if(stack.stack_size < 0){
        log.Fatal("error: stack underflow in peek")
    }

    result := stack.stack[stack.stack_size - 1]
    return result
}

func is_digit(val string)(bool){
    if _, err := strconv.Atoi(val); err == nil {
        return true 
    }
    return false 
}

func lex_file(filename string)([]Token) {
    bytes, err := os.ReadFile(filename)

    handle_error(err)

    contents := strings.Fields(string(bytes))

    tokens := []Token{}

    for _, val := range contents {
        var token Token
        switch(val) {
        case "+":
            token.category = add_type
        case "-":
            token.category = sub_type
        case ".":
            token.category = print_type
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
        case "then":
            token.category = then_type
        default:
            if is_digit(val) {
                token.category = push_type
                token.value = val
            } else {
                token.category = word_type
                token.value = val
            }
        }
        tokens = append(tokens, token)
    }

    return tokens
}

func interpret_tokens(tokens []Token, stack *Stack, words map[string][]Token, is_word bool) {
    in_word := false
    var current_word string 
    var word_tokens []Token
    
    for i := 0; i < len(tokens); i++ {
        val := tokens[i]
        if in_word {
            if val.category == semi_type {
                words[current_word] = word_tokens
                in_word = false
            } else if val.category == word_type {
                current_word = val.value
            } else {
                word_tokens = append(word_tokens, val) 
            }
        } else {
            switch(val.category) {
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
                    a := stack.stack[stack.stack_size - 2]
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
                case emit_type:
                    a := stack.pop()
                    fmt.Print(string(a))
                case cr_type:
                    fmt.Println()
                case colon_type:
                    in_word = true
                case semi_type:
                    in_word = false
                case word_type:
                    if in_word {
                        current_word = val.value
                    } else {
                        interpret_tokens(words[val.value], stack, words, true) 
                    }
                case if_type:
                    if is_word {
                        a := stack.peek()
                        if a != 0 {
                            
                        } else {
                            for index := 0; val.category != then_type; index++ {
                                val = tokens[i]
                                i++
                            }
                        }
                    } else {
                        log.Fatal("error if")
                    }
                case then_type:
            }
        }
    }

}

func main(){
    if len(os.Args) < 2 {
        log.Fatal(fmt.Sprintf("usage: %s <file_name.forth>", os.Args[0])) 
    }

    filename := os.Args[1]

    tokens := lex_file(filename)

    var stack Stack
    words := make(map[string][]Token)

    interpret_tokens(tokens, &stack, words, false)
}
