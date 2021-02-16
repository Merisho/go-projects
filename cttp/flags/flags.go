package flags

import (
    "errors"
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

var nonPointerError = errors.New("flags argument must be a pointer")

func Parse(str string, flags interface{}) error {
    flagsType := reflect.TypeOf(flags)
    if flagsType == nil || flagsType.Kind() != reflect.Ptr {
        return nonPointerError
    }

    flagsStructType := flagsType.Elem()
    flagsStructVal := reflect.ValueOf(flags).Elem()

    parsedFlags := parseFlagsStr(str)
    fieldTags := mapFieldTagsToNames(flagsStructType)

    if err := validateFlagsTypes(parsedFlags, fieldTags, flagsStructType); err != nil {
        return err
    }

    setFieldValues(parsedFlags, fieldTags, flagsStructType, flagsStructVal)

    return nil
}

func validateFlagsTypes(parsedFlags, fieldTags map[string]string, flagsStructType reflect.Type) error {
    for name, val := range parsedFlags {
        fieldName := fieldTags[name]
        fieldType, ok := flagsStructType.FieldByName(fieldName)
        if !ok {
            continue
        }

        switch fieldType.Type.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            _, err := strconv.ParseInt(val, 10, 64)
            if err != nil {
                return fmt.Errorf("invalid %s flag type", name)
            }
        }
    }

    return nil
}

func setFieldValues(parsedFlags, fieldTags map[string]string, flagsStructType reflect.Type, flagsStructVal reflect.Value) {
    for name, val := range parsedFlags {
        fieldName := fieldTags[name]
        fieldType, ok := flagsStructType.FieldByName(fieldName)
        if !ok {
            continue
        }

        fieldValue := flagsStructVal.FieldByName(fieldName)
        setFieldValue(fieldValue, fieldType, name, val)
    }
}

func setFieldValue(fieldValue reflect.Value, fieldType reflect.StructField, flagName, flagVal string) {
    switch fieldType.Type.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        intVal, _ := strconv.ParseInt(flagVal, 10, 64)
        fieldValue.SetInt(intVal)
    case reflect.String:
        fieldValue.SetString(flagVal)
    case reflect.Bool:
        fieldValue.SetBool(true)
    }
}

func parseFlagsStr(s string) map[string]string {
    start := strings.Index(s, "-")
    if start == -1 {
        return nil
    }

    str := s[start:]
    rawFlags := strings.Split(str, "-")

    flags := make(map[string]string)
    for _, rf := range rawFlags {
        name, val := extractNameValue(rf)
        flags[name] = val
    }

    return flags
}

func extractNameValue(rawFlag string) (name, value string) {
    firstSpace := strings.Index(rawFlag, " ")

    name = rawFlag
    if firstSpace != -1 {
        name = strings.TrimSpace(rawFlag[:firstSpace])
        value = strings.TrimSpace(rawFlag[firstSpace + 1:])
    }

    return name, value
}

func mapFieldTagsToNames(structType reflect.Type) map[string]string {
    fieldsNum := structType.NumField()

    fieldTags := make(map[string]string)
    for i := 0; i < fieldsNum; i++ {
        f := structType.Field(i)
        tag, ok := f.Tag.Lookup("flag")
        if !ok {
            continue
        }

        fieldTags[tag] = f.Name
    }

    return fieldTags
}
