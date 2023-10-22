package ap_tools

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type LineByLineHandler func(lineNumber int64, lineString string)

func ReadLineByLineAndProcess(fileName string, handler LineByLineHandler) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	cnt := int64(0)
	buf := bufio.NewReader(f)
	for {
		cnt++
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				return nil
			}
			return err
		}
		handler(cnt, line)
	}
}

func LineString2Struct(out interface{}, lineString, delimiter string) error {
	obj := reflect.ValueOf(out)
	if obj.Type().Kind() != reflect.Ptr {
		return nil
	}
	obj = obj.Elem()
	t := obj.Type()
	if t.Kind() != reflect.Struct {
		return nil
	}
	arr := strings.Split(lineString, delimiter)
	cnt := len(arr)
	if cnt < t.NumField() {
		return errors.New("number of fields does not match")
	}
	fieldIndexMap := make(map[int]int)
	for i := t.NumField() - 1; i > -1; i-- {
		tags := t.Field(i).Tag.Get("index")
		if len(tags) < 1 || tags == "-" {
			continue
		}
		idx, err := strconv.Atoi(tags)
		if err != nil {
			return err
		}
		fieldIndexMap[idx-1] = i
	}
	for i := 0; i < cnt; i++ {
		k, ok := fieldIndexMap[i]
		if !ok {
			continue
		}
		switch t.Field(k).Type.Kind() {
		case reflect.String:
			obj.Field(k).SetString(arr[i])
		case reflect.Bool:
			v, err := strconv.ParseBool(arr[i])
			if err != nil {
				return err
			}
			obj.Field(k).SetBool(v)
		case reflect.Int, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(arr[i], 10, 64)
			if err != nil {
				return err
			}
			obj.Field(k).SetInt(v)
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(arr[i], 10, 64)
			if err != nil {
				return err
			}
			obj.Field(k).SetUint(v)
		case reflect.Float32, reflect.Float64:
			v, err := strconv.ParseFloat(arr[i], 64)
			if err != nil {
				return err
			}
			obj.Field(k).SetFloat(v)
		}
	}
	return nil
}
