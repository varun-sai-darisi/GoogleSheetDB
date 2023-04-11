package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"net/http"
)

type SpreadSheet struct {
	ctx           context.Context
	service       *sheets.Service
	spreadSheetId string
	sheetName     string
}

func InitialiseSheet(ctx context.Context, service *sheets.Service, spreadSheetId string, sheetName string) *SpreadSheet {
	return &SpreadSheet{
		ctx:           ctx,
		service:       service,
		spreadSheetId: spreadSheetId,
		sheetName:     sheetName,
	}
}

func (sheet *SpreadSheet) AddDataToSpreadSheet(row *sheets.ValueRange) (*sheets.AppendValuesResponse, error) {
	response, err := sheet.service.Spreadsheets.Values.Append(sheet.spreadSheetId, sheet.sheetName, row).
		ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(sheet.ctx).Do()
	if err != nil || response.HTTPStatusCode != 200 {
		fmt.Printf("Error %v in adding entries to google sheet", err.Error())
		return nil, err
	}
	return response, nil
}

func (sheet *SpreadSheet) AddUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	row := &sheets.ValueRange{
		Values: [][]interface{}{{person.Name, person.MobileNumber}},
	}
	response, err := sheet.AddDataToSpreadSheet(row)
	if err != nil {
		fmt.Sprintf("error %v in adding details to the spread sheet", err.Error())
	}
	json.NewEncoder(w).Encode(response)
}
