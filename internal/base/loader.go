package base

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const TAG_KEY = "pbac-guardian"

const (
	// PATH_SEP is the separator used in paths to exported fields.
	// Paths have format "entity:field1:field2..." where ":" separates parts.
	// Example: "source:name", "target:user:email"
	PATH_SEP = ":"

	// TAG_SEP is the separator used inside struct field tags.
	// Tags can have format "field_name,flag1,flag2" where "," separates parts.
	// Currently only the first part (field name) is used.
	TAG_SEP = ","

	TIME_MODIFIER_SEP = "|"
	MODIFIERS_PARTS   = 2
)

type Loader struct{}

func (l *Loader) parseAsPath(path any) (*Entity, []string, error) {
	if path, ok := path.(string); ok {
		splited := strings.Split(path, PATH_SEP)

		if len(splited) == 0 {
			return nil, nil, NewErrInvalidPath(path, "path is empty")
		}

		entity := Entity(splited[0])
		if !entity.IsValid() {
			return nil, nil, NewErrInvalidPath(path, fmt.Sprintf("path allocates to unknown entity: %v", entity))
		}

		if len(splited) == 1 {
			return &entity, nil, nil
		}

		return &entity, splited[1:], nil
	}
	return nil, nil, NewErrInvalidType("string", reflect.TypeOf(path).String())
}

func (l *Loader) loadFromStruct(ctx context.Context, v reflect.Value, t reflect.Type, path []string) (any, error) {
	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get(TAG_KEY)

		if tag != "" {
			flags := strings.Split(tag, TAG_SEP)
			tagValue := strings.TrimSpace(flags[0])

			if tagValue != path[0] {
				continue
			}

			if !field.CanInterface() {
				return nil, NewErrInvalidPath(path[0], "got unexpored field")
			}

			if len(path) > 1 {
				return l.loadFromSourceTargetItem(ctx, field.Interface(), path[1:])
			}

			return field.Interface(), nil
		}
	}

	return nil, NewErrInvalidPath(path[0], "param doesn't exist")
}

func (l *Loader) loadFromMap(ctx context.Context, value reflect.Value, path []string) (any, error) {
	keyValue := reflect.ValueOf(path[0])
	mapValue := value.MapIndex(keyValue)

	if !mapValue.IsValid() {
		return nil, NewErrInvalidPath(path[0], "param doesn't exist")
	}

	if !mapValue.CanInterface() {
		return nil, NewErrInvalidPath(path[0], "got unexported value")
	}

	item := mapValue.Interface()

	if len(path) > 1 {
		return l.loadFromSourceTargetItem(ctx, item, path[1:])
	}

	return item, nil
}

func (l *Loader) loadFromEnv(path string) (any, error) {
	if value, ok := os.LookupEnv(path); ok {
		return value, nil
	}

	return nil, NewErrInvalidPath(path, "env doesn't exists")
}

func (l *Loader) loadFromSourceTargetItem(ctx context.Context, data any, path []string) (any, error) {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, NewErrInvalidType(fmt.Sprintf("%v or %v or %v", reflect.Pointer.String(), reflect.Struct.String(), reflect.Map.String()), v.Kind().String())
		}
		v = v.Elem()
	}

	if path == nil || len(path) == 0 {
		return data, nil
	}

	switch v.Kind() {
	case reflect.Struct:
		t := reflect.TypeOf(data)
		return l.loadFromStruct(ctx, v, t, path)
	case reflect.Map:
		return l.loadFromMap(ctx, v, path)
	}

	return nil, NewErrInvalidPath(path[0], "param doesn't exist")
}

func (l *Loader) loadTimeModications(path []string, target time.Time) (time.Time, error) {
	if len(path) == 0 {
		return target, nil
	}

	modifiers := strings.Split(path[0], TIME_MODIFIER_SEP)
	if len(modifiers) != MODIFIERS_PARTS {
		return target, NewErrInvalidPath(path[0], "expected int|day/hour/second/month/year")
	}

	modifierValue, err := strconv.Atoi(modifiers[0])
	if err != nil {
		return target, NewErrInvalidType(reflect.Int.String(), reflect.TypeOf(modifierValue).String())
	}

	switch modifiers[1] {
	case "day":
		return target.Add(time.Hour * 24 * MODIFIERS_PARTS), nil
	case "hour":
		return target.Add(time.Hour * MODIFIERS_PARTS), nil
	case "minute":
		return target.Add(time.Minute * MODIFIERS_PARTS), nil
	case "second":
		return target.Add(time.Second * MODIFIERS_PARTS), nil
	case "millisecond":
		return target.Add(time.Millisecond * MODIFIERS_PARTS), nil
	}

	return target, NewErrInvalidType("day|hour|minute|second|millisecond", modifiers[1])
}

func (l *Loader) loadTime(path []string) (time.Time, error) {
	var t time.Time

	switch path[0] {
	case "now":
		t = time.Now()
	default:
		return t, NewErrInvalidPath(path[0], "expected now, but got invalid")
	}

	return l.loadTimeModications(path[1:], t)
}

func (l *Loader) Load(ctx context.Context, source, target, item any, path any) (any, error) {
	entity, parsedPath, err := l.parseAsPath(path)
	if err != nil {
		return path, nil
	}

	switch *entity {
	case Entity_SOURCE:
		return l.loadFromSourceTargetItem(ctx, source, parsedPath)
	case Entity_TARGET:
		return l.loadFromSourceTargetItem(ctx, target, parsedPath)
	case Entity_ENV:
		return l.loadFromEnv(parsedPath[0])
	case Entity_ITEM:
		return l.loadFromSourceTargetItem(ctx, item, parsedPath)
	case Entity_TIME:
		return l.loadTime(parsedPath)
	}

	return nil, NewErrInvalidPath(path, fmt.Sprintf("unxpected entity: %v", entity))
}

func (l *Loader) IsPath(path any) bool {
	_, _, err := l.parseAsPath(path)
	return err == nil
}
