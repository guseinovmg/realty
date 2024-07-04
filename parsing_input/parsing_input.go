package parsing_input

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func Parse(r *http.Request, m any) error {
	//var m MyStruct

	// Try to parse input from the request body as JSON
	if r.Body != nil {
		ct := r.Header.Get("Content-Type")
		if ct != "" && strings.HasPrefix(ct, "application/json") {
			return json.NewDecoder(r.Body).Decode(&m)
		}
	}

	// If parsing from the request body failed or there was no body, try to parse from PostForm
	err := r.ParseForm()
	if err != nil {
		return err
	}

	t := reflect.TypeOf(m)
	if t.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		return errors.New("not a struct")
	}
	v := reflect.ValueOf(m)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		postTag := field.Tag.Get("post")
		postValue := r.PostFormValue(postTag)
		fieldValue := v.Elem().Field(i)
		if fieldValue.Kind() == reflect.String {
			fieldValue.SetString(postValue)
		}
		if fieldValue.Kind() == reflect.Float64 {
			float, err := strconv.ParseFloat(postValue, 64)
			if err != nil {
				return err
			}
			fieldValue.SetFloat(float)
		}
		if fieldValue.Kind() == reflect.Int64 {
			integer, err := strconv.ParseInt(postValue, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(integer)
		}
		if fieldValue.Kind() == reflect.Int8 {
			integer, err := strconv.ParseInt(postValue, 10, 8)
			if err != nil {
				return err
			}
			fieldValue.SetInt(integer)
		}
		fmt.Printf("Field %s has post tag %q\n", field.Name, postTag)
	}

	return nil
}
