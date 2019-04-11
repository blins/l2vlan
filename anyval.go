package main

import (
	"strconv"
	"fmt"
)

/*
Чтобы проще было работать со значениями которые передаются вот так неопределенно
 */
type AnyVal struct{
	Value interface{}
}

/*
Возвращает строковое представление.
Работает если значение имеет типы:
 - int
 - string
 - fmt.Stringer
 */
func (v AnyVal) String() string {
	switch tv := v.Value.(type) {
	case int:
		return strconv.Itoa(tv)
	case string:
		return tv
	case fmt.Stringer:
		return tv.String()
	}
	return ""
}

/*
Возвращает целочисленное представление.
Работает если значение имеет типы:
 - int
 - string
 - fmt.Stringer
*/
func (v AnyVal) Int() int {
	switch tv := v.Value.(type) {
	case int:
		return tv
	case string:
		i, _ := strconv.Atoi(tv)
		return i
	case fmt.Stringer:
		i,_ := strconv.Atoi(tv.String())
		return i
	}
	return 0
}

