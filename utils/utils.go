package utils

import (
	"bytes"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// Contains checks if the slice/array contains the value
func Contains(slice, val interface{}) bool {
	sliceVal := reflect.ValueOf(slice)
	for i := 0; i < sliceVal.Len(); i++ {
		if val == sliceVal.Index(i).Interface() {
			return true
		}
	}
	return false
}

// Select random value from slice/array
func RandomVal(slice interface{}) interface{} {
	slc := reflect.ValueOf(slice)
	// Seed the generator using the current time
	rand.Seed(time.Now().UnixNano())
	val := slc.Index(rand.Intn(slc.Len())).Interface()
	return val
}

// Inserts a string every n characters
func InsertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune('-')
		}
	}
	return buffer.String()
}

func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func Nth(i int) string {
	if i%10 == 1 && i != 11 {
		return strconv.Itoa(i) + "st"
	}
	if i%10 == 2 && i != 12 {
		return strconv.Itoa(i) + "nd"
	}
	if i%10 == 3 && i != 13 {
		return strconv.Itoa(i) + "rd"
	}
	return strconv.Itoa(i) + "th"
}
