package go-cherwell

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	deleteBusObRecURI        = "api/V1/deletebusinessobject/busobid/$/busobrecid/#"
	saveBusObRecURI          = "api/V1/savebusinessobject"
	getSearchResultsURI      = "api/V1/getsearchresults"
	getBusObTemplateURI      = "api/V1/getbusinessobjecttemplate"
	getQuickSearchResultsURI = "api/V1/getquicksearchresults"
	getBusObRecByRecIdURI    = "api/v1/getbusinessobject/busobid/$/busobrecid/#"
	getBusObRecByPublicIdURI = "api/v1/getbusinessobject/busobid/$/publicid/*"
	getBusObSummariesAllURI  = "api/v1/getbusinessobjectsummaries/type/All"
)

// # Cherwell structs
type Error struct {
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	HasError     bool   `json:"hasError,omitempty"`
}

// ## BusinessObject struct
type BusOb struct {
	BusObID         string  `json:"busObId,omitempty"`
	Name            string  `json:"name,omitempty"`
	DisplayName     string  `json:"displayName,omitempty"`
	FirstRecIdField string  `json:"firstRecIdField,omitempty"`
	GroupSummaries  []BusOb `json:"groupSummaries,omitempty"`
	RecIdFields     string  `json:"recIdFields,omitempty"`
	StateFieldId    string  `json:"stateFieldId,omitempty"`
	States          string  `json:"states,omitempty"`
	Group           bool    `json:"group,omitempty"`
	Lookup          bool    `json:"lookup,omitempty"`
	Major           bool    `json:"major,omitempty"`
	Supporting      bool    `json:"supporting,omitempty"`
}

// ## BusinessObjectRecord struct
type ObRec struct {
	BusObID               string        `json:"busObId,omitempty"`
	BusObRecID            string        `json:"busObRecId,omitempty"`
	BusObPublicID         string        `json:"busObPublicId,omitempty"`
	CacheKey              string        `json:"cacheKey,omitempty"`
	CacheScope            string        `json:"cacheScope,omitempty"`
	Fields                []Field       `json:"fields,omitempty"`
	Persist               bool          `json:"persist,omitempty"`
	FieldValidationErrors []interface{} `json:"fieldValidationErrors,omitempty"`
	Error
	NotificationTriggers []interface{} `json:"notificationTriggers,omitempty"`
	Links                []Link        `json:"links,omitempty"`
}

// ### Field struct of BusinessObjectRecord
type Field struct {
	Dirty       bool   `json:"dirty,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Html        string `json:"html,omitempty"`
	Name        string `json:"name,omitempty"`
	Value       string `json:"value,omitempty"`
	FieldID     string `json:"fieldId,omitempty"`
}

type Fields struct {
	Fields []Field `json:"fields,omitempty"`
	Error
}

// ## Authentication struct for Communication with the API
type Client struct {
	User          string `json:"username"`
	Password      string `json:"password,omitempty"`
	ClientID      string `json:"as:client_id"`
	BaseURI       string `json:"-"`
	Expires       string `json:".expires,omitempty"`
	Issued        string `json:".issued,omitempty"`
	Access_token  string `json:"access_token,omitempty"`
	Token_type    string `json:"token_type"`
	Refresh_token string `json:"refresh_token,omitempty"`
	Expires_in    int    `json:"expires_in,omitempty"`
	Accept        string `json:"Accept,omitempty"`
	Grant_Type    string `json:"grant_type,omitempty"`
	Auth_mode     string `json:"-"`
}

// ## QuickSearch struct
type QuickSearch struct {
	Error
	Groups []struct {
		Error
		IsBusObTarget          bool `json:"isBusObTarget,omitempty"`
		SimpleResultsListItems []SimpleResultsListItem `json:"simpleResultsListItems,omitempty"`
		SubTitle string `json:"subTitle,omitempty"`
		TargetID string `json:"targetId,omitempty"`
		Title    string `json:"title,omitempty"`
	} `json:"groups,omitempty"`
	Title string `json:"title,omitempty"`
}

type Filter struct {
	FieldID   string `json:"fieldId,omitempty"`
	FieldName string `json:"-"`
	Operator  string `json:"operator,omitempty"`
	Value     string `json:"value,omitempty"`
}

type Search struct {
	Association        string   `json:"association,omitempty"`
	BusObID            string   `json:"busObId,omitempty"`
	CustomGridDefID    string   `json:"customGridDefId,omitempty"`
	DateTimeFormatting string   `json:"dateTimeFormatting,omitempty"`
	FieldID            string   `json:"fieldId,omitempty"`
	Fields             []string `json:"fields,omitempty"`
	Filters            []Filter `json:"filters,omitempty"`
	IncludeAllFields   bool     `json:"includeAllFields,omitempty"`
	IncludeSchema      bool     `json:"includeSchema,omitempty"`
	PageNumber         int64    `json:"pageNumber,omitempty"`
	PageSize           int64    `json:"pageSize,omitempty"`
	PromptValues       []struct {
		BusObID                  string   `json:"busObId,omitempty"`
		CollectionStoreEntireRow string   `json:"collectionStoreEntireRow,omitempty"`
		CollectionValueField     string   `json:"collectionValueField,omitempty"`
		FieldID                  string   `json:"fieldId,omitempty"`
		ListReturnFieldID        string   `json:"listReturnFieldId,omitempty"`
		PromptID                 string   `json:"promptId,omitempty"`
		Value                    struct{} `json:"value,omitempty"`
		ValueIsRecID             bool     `json:"valueIsRecId,omitempty"`
	} `json:"promptValues,omitempty"`
	Scope      string `json:"scope,omitempty"`
	ScopeOwner string `json:"scopeOwner,omitempty"`
	SearchID   string `json:"searchId,omitempty"`
	SearchName string `json:"searchName,omitempty"`
	SearchText string `json:"searchText,omitempty"`
	Sorting    []struct {
		FieldID       string `json:"fieldId,omitempty"`
		SortDirection int64  `json:"sortDirection,omitempty"`
	} `json:"sorting,omitempty"`
}

type SearchResult struct {
	BusinessObjects []ObRec `json:"businessObjects,omitempty"`
	Error
	HasPrompts bool   `json:"hasPrompts,omitempty"`
	Links      []Link `json:"links,omitempty"`
	Prompts    []struct {
		AllowValuesOnly          bool     `json:"allowValuesOnly,omitempty"`
		BusObID                  string   `json:"busObId,omitempty"`
		CollectionStoreEntireRow string   `json:"collectionStoreEntireRow,omitempty"`
		CollectionValueField     string   `json:"collectionValueField,omitempty"`
		ConstraintXML            string   `json:"constraintXml,omitempty"`
		Contents                 string   `json:"contents,omitempty"`
		Default                  string   `json:"default,omitempty"`
		FieldID                  string   `json:"fieldId,omitempty"`
		IsDateRange              bool     `json:"isDateRange,omitempty"`
		ListDisplayOption        string   `json:"listDisplayOption,omitempty"`
		ListReturnFieldID        string   `json:"listReturnFieldId,omitempty"`
		MultiLine                bool     `json:"multiLine,omitempty"`
		PromptID                 string   `json:"promptId,omitempty"`
		PromptType               string   `json:"promptType,omitempty"`
		PromptTypeName           string   `json:"promptTypeName,omitempty"`
		Required                 bool     `json:"required,omitempty"`
		Text                     string   `json:"text,omitempty"`
		Value                    struct{} `json:"value,omitempty"`
		Values                   []string `json:"values,omitempty"`
	} `json:"prompts,omitempty"`
	SearchResultsFields []struct {
		Caption                   string `json:"caption,omitempty"`
		CurrencyCulture           string `json:"currencyCulture,omitempty"`
		CurrencySymbol            string `json:"currencySymbol,omitempty"`
		DecimalDigits             int64  `json:"decimalDigits,omitempty"`
		DefaultSortOrderAscending bool   `json:"defaultSortOrderAscending,omitempty"`
		DisplayName               string `json:"displayName,omitempty"`
		FieldID                   string `json:"fieldId,omitempty"`
		FieldName                 string `json:"fieldName,omitempty"`
		FullFieldID               string `json:"fullFieldId,omitempty"`
		HasDefaultSortField       bool   `json:"hasDefaultSortField,omitempty"`
		IsBinary                  bool   `json:"isBinary,omitempty"`
		IsCurrency                bool   `json:"isCurrency,omitempty"`
		IsDateTime                bool   `json:"isDateTime,omitempty"`
		IsFilterAllowed           bool   `json:"isFilterAllowed,omitempty"`
		IsLogical                 bool   `json:"isLogical,omitempty"`
		IsNumber                  bool   `json:"isNumber,omitempty"`
		IsShortDate               bool   `json:"isShortDate,omitempty"`
		IsShortTime               bool   `json:"isShortTime,omitempty"`
		IsVisible                 bool   `json:"isVisible,omitempty"`
		SortOrder                 string `json:"sortOrder,omitempty"`
		Sortable                  bool   `json:"sortable,omitempty"`
		StorageName               string `json:"storageName,omitempty"`
		WholeDigits               int64  `json:"wholeDigits,omitempty"`
	} `json:"searchResultsFields,omitempty"`
	SimpleResults struct {
		Error
		Groups []struct {
			Error
			IsBusObTarget          bool `json:"isBusObTarget,omitempty"`
			SimpleResultsListItems []SimpleResultsListItem `json:"simpleResultsListItems,omitempty"`
			SubTitle string `json:"subTitle,omitempty"`
			TargetID string `json:"targetId,omitempty"`
			Title    string `json:"title,omitempty"`
		} `json:"groups,omitempty"`
		Title string `json:"title,omitempty"`
	} `json:"simpleResults,omitempty"`
	TotalRows int64 `json:"totalRows,omitempty"`
}

type SimpleResultsListItem struct {
	ObRec
	DocRepositoryItemID string `json:"docRepositoryItemId,omitempty"`
	GalleryImage        string `json:"galleryImage,omitempty"`
	PublicID            string `json:"publicId,omitempty"`
	Scope               string `json:"scope,omitempty"`
	ScopeOwner          string `json:"scopeOwner,omitempty"`
	SubTitle            string `json:"subTitle,omitempty"`
	Text                string `json:"text,omitempty"`
	Title               string `json:"title,omitempty"`
}

type SaveResponse struct {
	ObRec
	Error
	FieldValidationErrors []struct {
		Error     string `json:"error"`
		ErrorCode string `json:"errorCode"`
		FieldID   string `json:"fieldId"`
	} `json:"fieldValidationErrors"`
	NotificationTriggers []struct {
		Key          string `json:"key"`
		SourceChange string `json:"sourceChange"`
		SourceID     string `json:"sourceId"`
		SourceType   string `json:"sourceType"`
	} `json:"notificationTriggers"`
}

// ## Link struct
type Link struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// # Methods

// ## NewClient
func NewClient(user, password, clientID, baseURI, auth_mode, grant_type string) Client {
	return Client{
		User:       user,
		Password:   password,
		ClientID:   clientID,
		BaseURI:    baseURI,
		Auth_mode:  auth_mode,
		Grant_Type: grant_type,
	}
}

// ## NewFilter
func (a *Client) NewFilter(bo BusOb, filters []Filter) ([]Filter, error) {
	filter := []Filter{}
	fs, err := a.GetBusObTemplate(bo.BusObID)
	if err != nil {
		fmt.Print(err.Error())
		return filter, err
	}
	for _, fi := range fs.Fields {
		for _, f := range filters {
			if fi.Name == f.FieldName {
	f.FieldID = fi.FieldID
	filter = append(filter, f)
			}
		}
	}
	empty := true
	for _, fi := range fs.Fields {
		for _, f := range filter {
			if fi.FieldID == f.FieldID {
	empty = false
			}
		}
	}
	if empty {
		return filter, errors.New("empty filter or misspelled field-name")
	}
	return filter, err
}

// ## Login
func (a *Client) Login() error {
	params := url.Values{}
	params.Add("grant_type", a.Grant_Type)
	params.Add("client_id", a.ClientID)
	params.Add("username", a.User)
	params.Add("password", a.Password)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", (a.BaseURI + "token?auth_mode=" + a.Auth_mode), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = UnJson(resp.Body, &a)
	return err
}

// ## Get all BusinessObjects
func (a Client) GetAllBusOb() ([]BusOb, error) {
	busObs := []BusOb{}
	params := url.Values{}
	params.Add("client_id", a.ClientID)
	body := strings.NewReader(params.Encode())
	uri := a.BaseURI + getBusObSummariesAllURI
	req, err := http.NewRequest("GET", uri, body)
	if err != nil {
		return busObs, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + a.Access_token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return busObs, err
	}
	defer resp.Body.Close()

	err = UnJson(resp.Body, &busObs)
	return busObs, err
}

// ## Get BusinessObject by DisplayName
func (a Client) GetBusOb(name string) (BusOb, error) {
	bo := BusOb{}
	bos, err := a.GetAllBusOb()
	if err != nil {
		fmt.Print(err.Error())
		return bo, err
	}
	for _, bob := range bos {
		if bob.DisplayName == name {
			bo = bob
			break
		}
		for _, b := range bob.GroupSummaries {
			if b.DisplayName == name {
	bo = b
	break
			}
		}
	}

	if bo.BusObID == "" {
		return bo, errors.New(("no business-object found: " + name))
	}
	return bo, err
}

// ## Get BusinessObjectRecord by PublicID
func (a Client) GetObRecByPublicID(bo BusOb, pid string) (ObRec, error) {
	rec := ObRec{}
	params := url.Values{}
	params.Add("client_id", a.ClientID)
	body := strings.NewReader(params.Encode())
	uri := a.BaseURI + strings.Replace(getBusObRecByPublicIdURI, "$", bo.BusObID, 1) + strings.Replace(getBusObRecByPublicIdURI, "*", pid, 1)
	req, err := http.NewRequest("GET", uri, body)
	if err != nil {
		return rec, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + a.Access_token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return rec, err
	}
	defer resp.Body.Close()

	err = UnJson(resp.Body, &rec)
	return rec, err
}

// ## Get BusinessObjectRecord by RecID
func (a Client) GetObRecByRecID(boid string, rid string) (ObRec, error) {
	rec := ObRec{}
	params := url.Values{}
	params.Add("client_id", a.ClientID)
	body := strings.NewReader(params.Encode())
	uri := a.BaseURI + strings.Replace(getBusObRecByRecIdURI, "$", boid, 1) + strings.Replace(getBusObRecByRecIdURI, "#", rid, 1)
	req, err := http.NewRequest("GET", uri, body)
	if err != nil {
		return rec, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + a.Access_token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return rec, err
	}
	defer resp.Body.Close()

	err = UnJson(resp.Body, &rec)
	return rec, err
}

// ## QuickSearch
func (a Client) QuickSearch(q string) (QuickSearch, error) {
	recs := QuickSearch{}
	data := struct {
		SearchText string `json:"searchText"`
	}{
		SearchText: q,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return recs, err
	}

	body := bytes.NewReader(payloadBytes)
	uri := a.BaseURI + getQuickSearchResultsURI
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return recs, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + a.Access_token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return recs, err
	}
	defer resp.Body.Close()

	err = UnJson(resp.Body, &recs)
	return recs, err
}

func (a Client) GetBusObTemplate(boID string) (Fields, error) {
	fields := Fields{}
	data := struct {
		BusObID         string `json:"busObId"`
		IncludeAll      bool   `json:"includeAll"`
		IncludeRequired bool   `json:"includeRequired"`
	}{
		BusObID:         boID,
		IncludeAll:      true,
		IncludeRequired: true,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return fields, err
	}

	body := bytes.NewReader(payloadBytes)
	uri := a.BaseURI + getBusObTemplateURI
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return fields, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + a.Access_token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fields, err
	}
	defer resp.Body.Close()

	err = UnJson(resp.Body, &fields)
	return fields, err
}

func (a Client) GetOpenIncidents() (SearchResult, error) {
	result := SearchResult{}
	bos, err := a.GetAllBusOb()
	if err != nil {
		return result, err
	}
	incident := BusOb{}
	for _, b := range bos {
		if b.DisplayName == "Incident" {
			incident = b
			break
		}
	}
	fs, err := a.GetBusObTemplate(incident.BusObID)
	if err != nil {
		fmt.Print(err.Error())
		return result, err
	}
	var fid string
	for _, f := range fs.Fields {
		if f.DisplayName == "Status" {
			fid = f.FieldID
		}
	}
	data := Search{
		BusObID: incident.BusObID,
		Filters: []Filter{
			{
	FieldID:  fid,
	Operator: "eq",
	Value:    "Neu",
			},
			{
	FieldID:  fid,
	Operator: "eq",
	Value:    "In Bearbeitung",
			},
		},
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	body := bytes.NewReader(payloadBytes)
	uri := a.BaseURI + getSearchResultsURI
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return result, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + a.Access_token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	err = UnJson(resp.Body, &result)
	return result, err
}

func (a Client) Search(bo BusOb, filter []Filter) (SearchResult, error) {
	result := SearchResult{}
	data := Search{
		BusObID:          bo.BusObID,
		Filters:          filter,
		IncludeAllFields: true,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", (a.BaseURI + "api/V1/getsearchresults"), body)
	if err != nil {
		return result, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + a.Access_token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	err = UnJson(resp.Body, &result)
	return result, err
}

func (a Client) SaveObRec(rec ObRec, fields []Field) (SaveResponse, error) {
	saveResp := SaveResponse{}
	template, err := a.GetBusObTemplate(rec.BusObID)
	if err != nil {
		return saveResp, err
	}

	rec.Fields = []Field{}
	for _, field := range template.Fields {
		for _, new := range fields {
			if field.Name == new.Name {
				new.Dirty = true
				new.FieldID = field.FieldID
				rec.Fields = append(rec.Fields, new)
			}
		}
	}
	rec.Persist = true
	payloadBytes, err := json.Marshal(rec)
	if err != nil {
		return saveResp, err
	}
	body := bytes.NewReader(payloadBytes)
	uri := a.BaseURI + saveBusObRecURI
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return saveResp, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + a.Access_token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return saveResp, err
	}
	defer resp.Body.Close()
	err = UnJson(resp.Body, &saveResp)
	return saveResp, err
}

func (a Client) DeleteObRec(rec ObRec) (ObRec, error) {
	respRec := ObRec{}
	uri := a.BaseURI + strings.Replace(deleteBusObRecURI, "$", rec.BusObID, 1) + strings.Replace(deleteBusObRecURI, "#", rec.BusObRecID, 1)
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return respRec, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + a.Access_token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return respRec, err
	}
	defer resp.Body.Close()
	err = UnJson(resp.Body, &respRec)
	return respRec, err
}

func (bo BusOb) NewObRec(fields []Field) ObRec {
	rec := ObRec{
		BusObID: bo.BusObID,
		Fields:  fields,
		Persist: true,
	}
	return rec
}

func UnJson(input io.ReadCloser, data interface{}) error {
	ioBody, err := ioutil.ReadAll(input)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	ioBody = bytes.TrimPrefix(ioBody, []byte("\xef\xbb\xbf"))
	body := []byte(string(ioBody[:]))
	err = json.Unmarshal([]byte(string(body)), &data)
	return err
}