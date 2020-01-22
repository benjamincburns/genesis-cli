package util

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

const indentStr = "    "

func pprintln(subj string, attr ...color.Attribute) {
	if viper.GetBool("no-colors") {
		fmt.Print(subj)
		return
	}
	fmt.Println(color.New(attr...).Sprint(subj))
}

func pprint(subj string, attr ...color.Attribute) {
	if viper.GetBool("no-colors") {
		fmt.Print(subj)
		return
	}
	fmt.Print(color.New(attr...).Sprint(subj))
}

func Printf(format string, a ...interface{}) {
	Print(fmt.Sprintf(format, a...))
}

func Print(i interface{}) {
	pprintln(fmt.Sprint(i), color.FgHiWhite)
}

func PlainPrintf(format string, a ...interface{}) {
	PlainPrint(fmt.Sprintf(format, a...))
}

func PlainPrint(i interface{}) {
	fmt.Print(i)
}

func AuthPrintf(format string, a ...interface{}) {
	AuthPrint(fmt.Sprintf(format, a...))
}

func AuthPrint(i interface{}) {
	pprintln(fmt.Sprint(i), color.FgCyan)
}

func PrintS(depth int, v interface{}) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += indentStr
	}
	if strger, ok := v.(fmt.Stringer); ok {
		pprintln(indent + strger.String())
		return
	}
	rv := reflect.ValueOf(v)
	t := rv.Type()
	if t.Kind() == reflect.Ptr && !rv.IsNil() {
		PrintS(depth, rv.Elem().Interface())
		return
	}
	if t.Kind() == reflect.Map {
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			PrintKV(depth, k.Interface(), v.Interface())
		}
		return
	}
	if t.Kind() == reflect.Slice {
		rv := reflect.ValueOf(v)
		for i := 0; i < rv.Len(); i++ {
			PrintS(depth+1, rv.Index(i).Interface())
		}
		return
	}
	if t.Kind() != reflect.Struct {
		pprintln(fmt.Sprintf(indent+"%+v", v))
		return
	}

	//reflect.ValueOf(v).Field(i).Kind()
	for i := 0; i < t.NumField(); i++ {
		if !rv.Field(i).CanInterface() {
			continue
		}
		field := t.Field(i)
		name := field.Name
		if field.Tag.Get("genesis") != "" {
			name = field.Tag.Get("genesis")
		}
		PrintKV(depth, name, rv.Field(i).Interface())
	}
}

func PrintKV(depth int, k interface{}, v interface{}) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += indentStr
	}
	if strger, ok := v.(fmt.Stringer); ok {
		pprint(indent+fmt.Sprint(k)+": ", color.FgYellow)
		pprintln(indent + strger.String())
		return
	}
	rv := reflect.ValueOf(v)
	t := rv.Type()
	if t.Kind() == reflect.Ptr && !rv.IsNil() {
		PrintS(depth, rv.Elem().Interface())
		return
	}
	if t.Kind() == reflect.Slice {
		pprintln(indent+fmt.Sprint(k)+": ", color.FgYellow)
		rv := reflect.ValueOf(v)
		for i := 0; i < rv.Len(); i++ {
			PrintS(depth+1, rv.Index(i).Interface())
		}
		return
	}
	if t.Kind() == reflect.Map {
		pprintln(indent+fmt.Sprint(k)+": ", color.FgYellow)
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			PrintKV(depth+1, k.Interface(), v.Interface())
		}
		return
	}
	if t.Kind() != reflect.Struct {
		pprint(indent+fmt.Sprint(k)+": ", color.FgYellow)
		pprintln(fmt.Sprintf("%+v", v))
		return
	}

	pprintln(indent+fmt.Sprint(k)+": ", color.FgYellow)
	//reflect.ValueOf(v).Field(i).Kind()
	for i := 0; i < t.NumField(); i++ {
		if !rv.Field(i).CanInterface() {
			continue
		}
		field := t.Field(i)
		name := field.Name
		if field.Tag.Get("genesis") != "" {
			name = field.Tag.Get("genesis")
		}
		PrintKV(depth+1, name, rv.Field(i).Interface())
	}
}
