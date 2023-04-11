package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	credentials, err := base64.StdEncoding.DecodeString(os.Getenv("SHEET_JSON_KEY"))
	if err != nil {
		fmt.Printf("Error %v in accessing creds of sheet api", err.Error())
		return
	}

	config, err := google.JWTConfigFromJSON(credentials, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		fmt.Printf("error %v in fetching config from jwt token", err.Error())
		return
	}

	client := config.Client(ctx)

	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		fmt.Printf("error %v in  creating service for the google sheets", err.Error())
		return
	}

	sheetId := 0
	spreadSheetId := "1ER-wfyojMbA15BzPx24DoUaNrmnOj69yOi95sCUQ8Xo"

	sheet, err := service.Spreadsheets.Get(spreadSheetId).Fields("sheets(properties(sheetId,title))").Do()

	if err != nil {
		fmt.Printf("error %v in accessing sheets ", err.Error())
		return
	}

	sheetName := ""
	for _, v := range sheet.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetId) {
			sheetName = prop.Title
			break
		}
	}

	sheetService := InitialiseSheet(ctx, service, spreadSheetId, sheetName)

	router := mux.NewRouter()

	router.HandleFunc("/add-details", sheetService.AddUserDetails).Methods(http.MethodPost)
	http.ListenAndServe(":8000", router)
}
