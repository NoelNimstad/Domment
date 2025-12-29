package main

import "regexp"

type ExpressionTemplate struct {
	Name       string
	Expression string
}

type Expression struct {
	Name       string
	Expression *regexp.Regexp
}

var r struct {
	Expressions []Expression
}

func InitialiseRegex() {
	expr := [...]ExpressionTemplate{
		{"Function", `func (.*)\((.*)\) ?(.*?) ?{`},
		{"Struct", `type (.*?) struct {`},
		{"Variable", `var (.*?) (.*?) = (.*?)`},
	}

	for _, v := range expr {
		regex, err := regexp.Compile(v.Expression)
		check(err)
		r.Expressions = append(r.Expressions, Expression{Name: v.Name, Expression: regex})
	}
}
