// Package templating is used to allow the text files to be dynamic. The intention is to
// allow things like configuration files or static files be be a bit dynamic.

package templating

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/template"
)

type OptionalString struct {
	ptr *string
}

var funcMap = template.FuncMap{
	"env":      env,
	"default":  defaultValue,
	"required": required,
}

func (s OptionalString) String() string {
	if s.ptr == nil {
		return ""
	}
	return *s.ptr
}

func env(key string) OptionalString {
	value, ok := os.LookupEnv(key)
	if !ok {
		return OptionalString{nil}
	}
	return OptionalString{&value}
}

func defaultValue(args ...interface{}) (string, error) {
	for _, arg := range args {
		if arg == nil {
			continue
		}
		switch v := arg.(type) {
		case string:
			return v, nil
		case *string:
			if v != nil {
				return *v, nil
			}
		case OptionalString:
			if v.ptr != nil {
				return *v.ptr, nil
			}
		default:
			return "", fmt.Errorf("Default: unsupported type '%T'", v)
		}
	}

	return "", errors.New("Default: all arguments are nil")
}

func required(arg interface{}) (string, error) {
	if arg == nil {
		return "", errors.New("Required argument is missing")
	}

	switch value := arg.(type) {
	case string:
		return value, nil
	case *string:
		if value != nil {
			return *value, nil
		}
	case OptionalString:
		if value.ptr != nil {
			return *value.ptr, nil
		}
	default:
		return "", fmt.Errorf("Requires: unsupported type '%T'", value)
	}
	return "", nil
}

// GenerateTemplate will action all the functions on the configuration file
func GenerateTemplate(source []byte) ([]byte, error) {
	tplt, err := template.New("file").Funcs(funcMap).Parse(string(source))
	if err != nil {
		return nil, fmt.Errorf("failed to create template. Error: %s", err)
	}

	var buffer bytes.Buffer
	if err = tplt.Execute(&buffer, nil); err != nil {
		return nil, fmt.Errorf("failed to transform template. Error: %s", err)
	}
	return buffer.Bytes(), nil
}
