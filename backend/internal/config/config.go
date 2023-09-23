package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func LoadEnv(v any) error {
	store, err := readEnv()
	if err != nil {
		return err
	}

	return populateConfig(v, store)
}

func readEnv() (map[string]string, error) {
	f, err := os.Open(".env")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	kv := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if value == "" {
				return nil, fmt.Errorf("reading .env: invalid value for key %s", key)
			}
			kv[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return kv, nil
}

func populateConfig(config any, kv map[string]string) error {
	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr {
		return errors.New("config must be a pointer")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("config must be a pointer to a struct")
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		tag := f.Tag.Get("env")
		if tag == "" {
			continue
		}

		v, ok := kv[tag]
		if !ok {
			return errors.New("missing env var: " + f.Tag.Get("env"))
		}

		switch f.Type.Kind() {
		case reflect.Bool:
			vb, err := strconv.ParseBool(v)
			if err != nil {
				return newConversionError(tag, f, v)
			}
			rv.Field(i).SetBool(vb)
		case reflect.Int:
			vi, err := strconv.Atoi(v)
			if err != nil {
				return newConversionError(tag, f, v)
			}
			rv.Field(i).SetInt(int64(vi))
		case reflect.String:
			rv.Field(i).SetString(v)
		default:
			return errors.New("unsupported type: " + f.Type.Kind().String())
		}
	}

	return nil
}

func newConversionError(tag string, field reflect.StructField, value string) error {
	return fmt.Errorf(
		"invalid value for %s type %s: %s",
		tag,
		field.Type.Kind().String(),
		value,
	)
}
