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
    drop_type
    dup_type
    swap_type
    over_type
    rot_type
    print_type
    emit_type
    cr_type
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

func main(){
    bytes, err := os.ReadFile("test.forth")

    handle_error(err)

    contents := strings.Fields(string(bytes))

    tokens := []Token{}

    for _, val := range contents {
        var token Token
        if val == "+" {
            token.category = add_type
        } else if val == "-"{
            token.category = sub_type
        } else if val == "." {
            token.category = print_type
        } else if val == "*" {
            token.category = mul_type
        } else if val == "/" {
            token.category = div_type
        } else if val == "drop" {
            token.category = drop_type
        } else if val == "dup" {
            token.category = dup_type
        } else if val == "swap" {
            token.category = swap_type
        } else if val == "over" {
            token.category = over_type
        } else if val == "rot" {
            token.category = rot_type
        } else if val == "emit" {
            token.category = emit_type
        } else if val == "cr" {
            token.category = cr_type
        } else if is_digit(val) {
            token.category = push_type
            token.value = val
        } else {
            log.Fatal("error")
        }
        tokens = append(tokens, token)
    }

    var stack Stack

    for _, val := range tokens {
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
        }
    }
}
