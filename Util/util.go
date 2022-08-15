package Util

import (
	"math/rand"
	"reflect"
	"sync"
)

const mutexLocked = 1

func MutexLocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

func Random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func Contain(s []reflect.Value, str string) bool {
	for _, v := range s {
		if v.Interface().(string) == str {
			return true
		}
	}

	return false
}
