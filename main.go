package main

import (
	"fmt"
	"reflect"

	"github.com/gocql/gocql"
)

type ClientOsHostname struct {
	Date     string                    `cql:"date" json:"date"`
	Deviceid gocql.UUID                `cql:"deviceid" json:"deviceid"`
	Epoch    int64                     `cql:"epoch" json:"epoch"`
	Clients  []ClientOsHoatnameClients `cql:"clients" json:"clients"`
}

type ClientOsHoatnameClients struct {
	DetectionMethod *string `cql:"detection_method" json:"detection_method"`
	Hostname        *string `cql:"hostname" json:"hostname"`
	Mac             *string `cql:"mac" json:"mac"`
	Os              *string `cql:"os" json:"os"`
}

type FieldInfo struct {
	Name  string
	Index int
	Type  reflect.Type
}

func ReadStruct(structMap map[string]FieldInfo, st interface{}) {
	sT := reflect.TypeOf(st)
	sV := reflect.ValueOf(st)

	// Return if not struct or pointer to struct.
	if sT.Kind() == reflect.Ptr {
		sT = sT.Elem()
	}
	if sT.Kind() != reflect.Struct {
		return
	}
	fmt.Printf("-- type: %v, %v, %v\n", sV.Type().Name(), sT.Kind().String(), sT.String())

	for i := 0; i < sT.NumField(); i++ {
		field := sT.Field(i)
		value := sV.Field(i)
		fmt.Printf("== %+v, %+v, %+v\n", field, field.Type, value)

		fmt.Printf("##%v\t,%v\t,[%v/%v] %v, %v\n", i, field.Tag.Get("json"), field.Type, field.Type.Kind(), field.Name, field)

		switch field.Type.Kind() {
		case reflect.Struct:
			ReadStruct(structMap, value.Interface())
		case reflect.Slice:
			fmt.Println(reflect.TypeOf(field.Type).Elem())

			for j := 0; j < value.Len(); j++ {
				ReadStruct(structMap, value.Index(j).Interface())
			}
			fallthrough
		case reflect.Ptr:
			structMap[field.Tag.Get("json")] = FieldInfo{
				Name:  field.Name,
				Index: field.Index[0],
				Type:  field.Type.Elem(),
			}
		default:
			structMap[field.Tag.Get("json")] = FieldInfo{
				Name:  field.Name,
				Index: field.Index[0],
				Type:  field.Type,
			}
		}
	}
}

func main() {
	structMap := make(map[string]FieldInfo)
	ReadStruct(structMap, ClientOsHostname{Clients: []ClientOsHoatnameClients{ClientOsHoatnameClients{}}})
	// ReadStruct1(reflect.TypeOf(ClientOsHostname{}))'
	fmt.Printf("===(%v)===\n%+v", len(structMap), structMap)
}
