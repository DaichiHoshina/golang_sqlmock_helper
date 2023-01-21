package main

import (
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"time"
)

type ShippingSlipEntry struct {
	ID                 int64
	PurchaseRequestID  string
	ShippingSlipID     int64
	CustomerID         string
	ConsignorName      string
	PrintStartLocation uint32
	LabelKind          uint32
}

type (
	GetShippingSlipEntryOption = func(s *ShippingSlipEntry)
)

func GetShippingSlipEntry(opts ...GetShippingSlipEntryOption) *ShippingSlipEntry {
	s := &ShippingSlipEntry{
		ID:                 6,
		PurchaseRequestID:  "05ce08c4-1361-41f0-b1f2-611a15d1475f",
		ShippingSlipID:     4,
		CustomerID:         "gid://shopify/Customer/5941121417270",
		ConsignorName:      "渋谷倉庫",
		PrintStartLocation: 2,
		LabelKind:          1,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func localtime() time.Time {
	loc, _ := time.LoadLocation("Local")
	return time.Date(2022, 10, 24, 12, 0, 0, 0, loc)
}

func ModelToRows(param interface{}) *sqlmock.Rows {
	dstType := reflect.TypeOf(param)
	dstValue := reflect.ValueOf(param)

	var allValues [][]driver.Value

	allValues = append(allValues, valuesFromModel(dstType, dstValue))

	rows := sqlmock.NewRows(columnsFromModelType(dstType))
	for _, row := range allValues {
		rows.AddRow(row...)
	}

	return rows
}

func valuesFromModel(dstType reflect.Type, dstValue reflect.Value) []driver.Value {
	var values []driver.Value
	for j := 0; j < dstValue.NumField(); j++ {
		field := dstType.Field(j)
		value := dstValue.FieldByName(field.Name)
		valueInterface := value.Interface()
		values = append(values, valueToDriverValue(valueInterface))
	}
	return values
}

func valueToDriverValue(valueInterface interface{}) driver.Value {
	switch value := valueInterface.(type) {
	case int, int64, int32, int16, int8:
		return value
	case uint, uint64, uint32, uint16, uint8:
		return value
	case float32:
		return value
	case float64:
		return value
	case string:
		return value
	case bool:
		if value {
			return "true"
		}
		return "false"
	case time.Time:
		return localtime()
	default:
		return nil
	}
}

func columnsFromModelType(dstType reflect.Type) []string {
	var columns []string
	for i := 0; i < dstType.NumField(); i++ {
		columns = append(columns, dstType.Field(i).Name)
	}
	return columns
}

func main() {
	result := GetShippingSlipEntry()

	rows := ModelToRows(*result)
	fmt.Println(rows)
}
