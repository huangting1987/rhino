package core

import (
	"fmt"
	"reflect"
)

//Object constructor
type ObjectCreator map[int]reflect.Type

func (hset ObjectCreator) Register(kind int, value interface{}) { //obj=Test{}
	hset[kind] = reflect.TypeOf(value)
}

func (hset ObjectCreator) UnRegister(kind int) {
	delete(hset, kind)
}

func (hset ObjectCreator) New(kind int) (val interface{}, err error) { //ptr
	if field, ok := hset[kind]; ok {
		val = reflect.New(field).Interface()
	} else {
		err = fmt.Errorf("can't find class kind=%d", kind)
	}
	return
}

/*
* type of
 */
func Typeof(p interface{}) string {
	if p == nil {
		return "<nil>"
	}
	return reflect.TypeOf(p).String()
}

/*
*size of
 */
func Sizeof(v interface{}) int {
	return int(reflect.TypeOf(v).Size())
}

/*
*所有引用
 */
func SizeAllof(data interface{}) int {
	var npm = &typeStruct{make(map[uintptr]bool), 0}
	num := npm.sizeof(reflect.ValueOf(data))
	return num
}

//Including type size 包括type的大小
func SizeTypeof(data interface{}) int {
	var npm = &typeStruct{make(map[uintptr]bool), 0}
	num := npm.sizeof(reflect.ValueOf(data))
	return num + npm.exNum
}

//class
type typeStruct struct {
	npm   map[uintptr]bool
	exNum int
}

func (s *typeStruct) sizeof(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Map:
		sum := 0
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapkey := keys[i]
			num := s.sizeof(mapkey)
			if num < 0 {
				return -1
			}
			sum += num
			num = s.sizeof(v.MapIndex(mapkey))
			if num < 0 {
				return -1
			}
			sum += num
		}
		s.exNum += int(v.Type().Size())
		return sum

	case reflect.Slice:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			num := s.sizeof(v.Index(i))
			if num < 0 {
				return -1
			}
			sum += num
		}
		s.exNum += int(v.Type().Size())
		return sum

	case reflect.Array:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			num := s.sizeof(v.Index(i))
			if num < 0 {
				return -1
			}
			sum += num
		}
		return sum

	case reflect.String:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			num := s.sizeof(v.Index(i))
			if num < 0 {
				return -1
			}
			sum += num
		}
		s.exNum += int(v.Type().Size())
		return sum

	case reflect.Ptr:
		s.exNum += int(v.Type().Size())
		if v.IsNil() {
			return 0
		}
		if _, ok := s.npm[v.Pointer()]; ok {
			return 0
		} else {
			s.npm[v.Pointer()] = true
		}
		return s.sizeof(v.Elem())

	case reflect.Interface:
		s.exNum += int(v.Type().Size())
		if v.IsNil() {
			return 0
		}
		return s.sizeof(v.Elem())

	case reflect.Uintptr: //Don't think it's Pointer
		return int(v.Type().Size())

	case reflect.UnsafePointer: //Don't think it's Pointer
		return int(v.Type().Size())

	case reflect.Struct:
		sum := 0
		for i, n := 0, v.NumField(); i < n; i++ {
			if v.Type().Field(i).Tag.Get("ss") == "-" {
				continue
			}
			num := s.sizeof(v.Field(i))
			if num < 0 {
				return -1
			}
			sum += num
		}
		return sum

	case reflect.Func, reflect.Chan:
		s.exNum += int(v.Type().Size())
		if v.IsNil() {
			return 0
		}
		return 0 //Temporary non handling func,chan.
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int:
		return int(v.Type().Size())
	case reflect.Bool:
		return int(v.Type().Size())
	default:
		fmt.Println("t.Kind() no found:", v.Kind())
	}
	return -1
}
