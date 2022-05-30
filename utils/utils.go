package utils

import (
	"io"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func PrintTable[T any](data []T, w io.Writer) error {
	if len(data) == 0 {
		return errors.New("Data must contain at least one entry")
	}

	x := reflect.ValueOf(data[0])
	var headers []string
	for i := 0; i < x.Type().NumField(); i++ {
		headerName := x.Type().Field(i).Name
		headers = append(headers, strings.Join(camelcase.Split(headerName), " "))
	}

	table := tablewriter.NewWriter(w)
	table.SetHeader(headers)
	for _, entry := range data {
		e := reflect.ValueOf(entry)
		var row []string
		for i := 0; i < len(headers); i++ {
			row = append(row, e.Field(i).String())
		}
		table.Append(row)
	}

	table.Render()

	return nil
}
