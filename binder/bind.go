// DefaultBinder 参考Echo中的实现
// 做了注释和少量调整

/*
**********************
向Echo开发团队致敬!
向Echo开发者致敬!
向Echo 致敬！
**********************
*/

package binder

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/twiglab/twig"
)

const (
	defBinderID = "_twig_default_binder_"
)

type BindUnmarshaler interface {
	UnmarshalParam(param string) error
}

type DefaultBinder struct{}

func New() *DefaultBinder {
	return &DefaultBinder{}
}

func (b *DefaultBinder) ID() string {
	return defBinderID
}

func (b *DefaultBinder) Name() string {
	return defBinderID
}

func (b *DefaultBinder) Bind(i interface{}, c twig.Ctx) (err error) {
	req := c.Req()
	if req.ContentLength == 0 {
		if req.Method == http.MethodGet || req.Method == http.MethodDelete {
			if err = b.bindData(i, c.QueryParams(), "query"); err != nil {
				return twig.NewHttpError(http.StatusBadRequest, err.Error()).SetInternal(err)
			}
			return
		}
		return twig.NewHttpError(http.StatusBadRequest, "Request body can't be empty")
	}
	ctype := req.Header.Get(twig.HeaderContentType)
	switch {
	case strings.HasPrefix(ctype, twig.MIMEApplicationJSON):
		if err = json.NewDecoder(req.Body).Decode(i); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return twig.NewHttpError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
			} else if se, ok := err.(*json.SyntaxError); ok {
				return twig.NewHttpError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
			} else {
				return twig.NewHttpError(http.StatusBadRequest, err.Error()).SetInternal(err)
			}
			return twig.NewHttpError(http.StatusBadRequest, err.Error())
		}
	case strings.HasPrefix(ctype, twig.MIMEApplicationXML), strings.HasPrefix(ctype, twig.MIMETextXML):
		if err = xml.NewDecoder(req.Body).Decode(i); err != nil {
			if ute, ok := err.(*xml.UnsupportedTypeError); ok {
				return twig.NewHttpError(http.StatusBadRequest, fmt.Sprintf("Unsupported type error: type=%v, error=%v", ute.Type, ute.Error())).SetInternal(err)
			} else if se, ok := err.(*xml.SyntaxError); ok {
				return twig.NewHttpError(http.StatusBadRequest, fmt.Sprintf("Syntax error: line=%v, error=%v", se.Line, se.Error())).SetInternal(err)
			} else {
				return twig.NewHttpError(http.StatusBadRequest, err.Error()).SetInternal(err)
			}
			return twig.NewHttpError(http.StatusBadRequest, err.Error())
		}
	case strings.HasPrefix(ctype, twig.MIMEApplicationForm), strings.HasPrefix(ctype, twig.MIMEMultipartForm):
		params, err := c.FormParams()
		if err != nil {
			return twig.NewHttpError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
		if err = b.bindData(i, params, "form"); err != nil {
			return twig.NewHttpError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
	default:
		return twig.ErrUnsupportedMediaType
	}
	return
}

func (b *DefaultBinder) bindData(ptr interface{}, data map[string][]string, tag string) error {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()

	if typ.Kind() != reflect.Struct {
		return errors.New("binding element must be a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}
		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get(tag)

		if inputFieldName == "" {
			inputFieldName = typeField.Name
			// If tag is nil, we inspect if the field is a struct.
			if _, ok := bindUnmarshaler(structField); !ok && structFieldKind == reflect.Struct {
				if err := b.bindData(structField.Addr().Interface(), data, tag); err != nil {
					return err
				}
				continue
			}
		}

		inputValue, exists := data[inputFieldName]
		if !exists {
			// Go json.Unmarshal supports case insensitive binding.  However the
			// url params are bound case sensitive which is inconsistent.  To
			// fix this we must check all of the map values in a
			// case-insensitive search.
			inputFieldName = strings.ToLower(inputFieldName)
			for k, v := range data {
				if strings.ToLower(k) == inputFieldName {
					inputValue = v
					exists = true
					break
				}
			}
		}

		if !exists {
			continue
		}

		// Call this first, in case we're dealing with an alias to an array type
		if ok, err := unmarshalField(typeField.Type.Kind(), inputValue[0], structField); ok {
			if err != nil {
				return err
			}
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for j := 0; j < numElems; j++ {
				if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
			return err

		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	// But also call it here, in case we're dealing with an array of BindUnmarshalers
	if ok, err := unmarshalField(valueKind, val, structField); ok {
		return err
	}

	switch valueKind {
	case reflect.Ptr:
		return setWithProperType(structField.Elem().Kind(), val, structField.Elem())
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}

func unmarshalField(valueKind reflect.Kind, val string, field reflect.Value) (bool, error) {
	switch valueKind {
	case reflect.Ptr:
		return unmarshalFieldPtr(val, field)
	default:
		return unmarshalFieldNonPtr(val, field)
	}
}

// bindUnmarshaler attempts to unmarshal a reflect.Value into a BindUnmarshaler
func bindUnmarshaler(field reflect.Value) (BindUnmarshaler, bool) {
	ptr := reflect.New(field.Type())
	if ptr.CanInterface() {
		iface := ptr.Interface()
		if unmarshaler, ok := iface.(BindUnmarshaler); ok {
			return unmarshaler, ok
		}
	}
	return nil, false
}

func unmarshalFieldNonPtr(value string, field reflect.Value) (bool, error) {
	if unmarshaler, ok := bindUnmarshaler(field); ok {
		err := unmarshaler.UnmarshalParam(value)
		field.Set(reflect.ValueOf(unmarshaler).Elem())
		return true, err
	}
	return false, nil
}

func unmarshalFieldPtr(value string, field reflect.Value) (bool, error) {
	if field.IsNil() {
		// Initialize the pointer to a nil value
		field.Set(reflect.New(field.Type().Elem()))
	}
	return unmarshalFieldNonPtr(value, field.Elem())
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func Bind(i interface{}, c twig.Ctx) error {
	binder, ok := GetBinder(defBinderID, c)
	if !ok {
		return fmt.Errorf("not found binder id(%s)", defBinderID)
	}
	return binder.Bind(i, c)
}
