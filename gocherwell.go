// gocherwell is a go-Module for communication with the REST-API of Cherwell
package gocherwell

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Cherwell API URIs.
const (
	deleteBusObRecURI        = "api/V1/deletebusinessobject/busobid/$/busobrecid/#"
	saveBusObRecURI          = "api/V1/savebusinessobject"
	getSearchResultsURI      = "api/V1/getsearchresults"
	getBusObTemplateURI      = "api/V1/getbusinessobjecttemplate"
	getQuickSearchResultsURI = "api/V1/getquicksearchresults"
	getBusObRecByRecIdURI    = "api/v1/getbusinessobject/busobid/$/busobrecid/#"
	getBusObRecByPublicIdURI = "api/v1/getbusinessobject/busobid/$/publicid/*"
	getBusObSchemaURI        = "api/v1/getbusinessobjectschema/busobid/$?includerelationships=true"
	getBusObSummariesAllURI  = "api/v1/getbusinessobjectsummaries/type/All"
	getRelatedBusObURI       = "api/V1/getrelatedbusinessobject/parentbusobid/$/parentbusobrecid/#/relationshipid/?"
	linkBusObRecURI          = "api/V2/linkrelatedbusinessobject/parentbusobid/$/parentbusobrecid/#/relationshipid/?/busobid/&/busobrecid/+"
	unlinkBusObRecURI        = "api/V1/unlinkrelatedbusinessobject/parentbusobid/$/parentbusobrecid/#/relationshipid/?/busobid/&/busobrecid/+"
)

// Mapping of the Placeholders for the Cherwell API URIs.
const (
	busobid         = "$"
	busobrecid      = "#"
	busobpublicid   = "*"
	relationshipid  = "?"
	childbusobid    = "&"
	childbusobrecid = "+"
)

// Client contains the necessary Values to communicate with the Cherwell API.
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

// BusinessObject contains the Values of a Cherwell BusinessObject.
type BusinessObject struct {
	BusObID         string           `json:"busObId,omitempty"`
	Name            string           `json:"name,omitempty"`
	DisplayName     string           `json:"displayName,omitempty"`
	FirstRecIdField string           `json:"firstRecIdField,omitempty"`
	GroupSummaries  []BusinessObject `json:"groupSummaries,omitempty"`
	RecIdFields     string           `json:"recIdFields,omitempty"`
	StateFieldId    string           `json:"stateFieldId,omitempty"`
	States          string           `json:"states,omitempty"`
	Group           bool             `json:"group,omitempty"`
	Lookup          bool             `json:"lookup,omitempty"`
	Major           bool             `json:"major,omitempty"`
	Supporting      bool             `json:"supporting,omitempty"`
}

// BusinessObject contains the Values of a Cherwell BusinessObjectSchema.
type BusinessObjectSchema struct {
	BusObID string `json:"busObId"`
	Error
	FieldDefinitions []struct {
		AutoFill             bool   `json:"autoFill"`
		Calculated           bool   `json:"calculated"`
		Category             string `json:"category"`
		DecimalDigits        int64  `json:"decimalDigits"`
		Description          string `json:"description"`
		Details              string `json:"details"`
		DisplayName          string `json:"displayName"`
		Enabled              bool   `json:"enabled"`
		FieldID              string `json:"fieldId"`
		HasDate              bool   `json:"hasDate"`
		HasTime              bool   `json:"hasTime"`
		IsFullTextSearchable bool   `json:"isFullTextSearchable"`
		MaximumSize          string `json:"maximumSize"`
		Name                 string `json:"name"`
		ReadOnly             bool   `json:"readOnly"`
		Required             bool   `json:"required"`
		Type                 string `json:"type"`
		TypeLocalized        string `json:"typeLocalized"`
		Validated            bool   `json:"validated"`
		WholeDigits          int64  `json:"wholeDigits"`
	} `json:"fieldDefinitions"`
	FirstRecIDField string `json:"firstRecIdField"`
	GridDefinitions []struct {
		DisplayName string `json:"displayName"`
		GridID      string `json:"gridId"`
		Name        string `json:"name"`
	} `json:"gridDefinitions"`
	HTTPStatusCode string `json:"httpStatusCode"`
	Name           string `json:"name"`
	RecIDFields    string `json:"recIdFields"`
	Relationships  []struct {
		Cardinality      string `json:"cardinality"`
		Description      string `json:"description"`
		DisplayName      string `json:"displayName"`
		FieldDefinitions []struct {
			AutoFill             bool   `json:"autoFill"`
			Calculated           bool   `json:"calculated"`
			Category             string `json:"category"`
			DecimalDigits        int64  `json:"decimalDigits"`
			Description          string `json:"description"`
			Details              string `json:"details"`
			DisplayName          string `json:"displayName"`
			Enabled              bool   `json:"enabled"`
			FieldID              string `json:"fieldId"`
			HasDate              bool   `json:"hasDate"`
			HasTime              bool   `json:"hasTime"`
			IsFullTextSearchable bool   `json:"isFullTextSearchable"`
			MaximumSize          string `json:"maximumSize"`
			Name                 string `json:"name"`
			ReadOnly             bool   `json:"readOnly"`
			Required             bool   `json:"required"`
			Type                 string `json:"type"`
			TypeLocalized        string `json:"typeLocalized"`
			Validated            bool   `json:"validated"`
			WholeDigits          int64  `json:"wholeDigits"`
		} `json:"fieldDefinitions"`
		RelationshipID string `json:"relationshipId"`
		Target         string `json:"target"`
	} `json:"relationships"`
	StateFieldID string `json:"stateFieldId"`
	States       string `json:"states"`
}

// BusinessObject contains the Values of a Cherwell BusinessObjectTemplate.
type BusinessObjectTemplate struct {
	Error
	Fields []Field `json:"fields"`
}

// BusinessObject contains the Values of a Cherwell BusinessObjectRecord
// and is used to Unmarshal multiple HTTP-Responses of the Cherwell API.
type BusinessObjectRecord struct {
	BusObID               string        `json:"busObId,omitempty"`
	BusObRecID            string        `json:"busObRecId,omitempty"`
	BusObPublicID         string        `json:"busObPublicId,omitempty"`
	CacheKey              string        `json:"cacheKey,omitempty"`
	CacheScope            string        `json:"cacheScope,omitempty"`
	Fields                []Field       `json:"fields,omitempty"`
	Persist               bool          `json:"persist,omitempty"`
	FieldValidationErrors []interface{} `json:"fieldValidationErrors,omitempty"`
	Error
	NotificationTriggers []interface{}          `json:"notificationTriggers,omitempty"`
	Links                []Link                 `json:"links,omitempty"`
	FieldValues          map[string]interface{} `json:"-"`
}

// RelatedBusinessObjects is used to Unmarshal the HTTP-Response of the Cherwell API
// regarding related BusinessObjects.
type RelatedBusinessObjects struct {
	Error
	Links                  []Link                 `json:"links"`
	PageNumber             int64                  `json:"pageNumber"`
	PageSize               int64                  `json:"pageSize"`
	ParentBusObID          string                 `json:"parentBusObId"`
	ParentBusObPublicID    string                 `json:"parentBusObPublicId"`
	ParentBusObRecID       string                 `json:"parentBusObRecId"`
	RelatedBusinessObjects []BusinessObjectRecord `json:"relatedBusinessObjects"`
	RelationshipID         string                 `json:"relationshipId"`
	TotalRecords           int64                  `json:"totalRecords"`
}

// Search is used to Marshal the HTTP-Request to the Cherwell API
// regarding Searches.
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

// SearchResult is used to Unmarshal the HTTP-Response of the Cherwell API
// regarding Searches.
type SearchResult struct {
	BusinessObjects []BusinessObjectRecord `json:"businessObjects,omitempty"`
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
			IsBusObTarget          bool                    `json:"isBusObTarget,omitempty"`
			SimpleResultsListItems []SimpleResultsListItem `json:"simpleResultsListItems,omitempty"`
			SubTitle               string                  `json:"subTitle,omitempty"`
			TargetID               string                  `json:"targetId,omitempty"`
			Title                  string                  `json:"title,omitempty"`
		} `json:"groups,omitempty"`
		Title string `json:"title,omitempty"`
	} `json:"simpleResults,omitempty"`
	TotalRows int64 `json:"totalRows,omitempty"`
}

// SimplehResultsListItem is used to Unmarshal multiple HTTP-Responses of the Cherwell API
// regarding Searches.
// Extends SearchResult
type SimpleResultsListItem struct {
	BusinessObjectRecord
	DocRepositoryItemID string `json:"docRepositoryItemId,omitempty"`
	GalleryImage        string `json:"galleryImage,omitempty"`
	PublicID            string `json:"publicId,omitempty"`
	Scope               string `json:"scope,omitempty"`
	ScopeOwner          string `json:"scopeOwner,omitempty"`
	SubTitle            string `json:"subTitle,omitempty"`
	Text                string `json:"text,omitempty"`
	Title               string `json:"title,omitempty"`
}

// Filter is used to Marshal the HTTP-Response to the Cherwell API
// regarding Searches.
// Extends Search
type Filter struct {
	FieldID   string `json:"fieldId,omitempty"`
	FieldName string `json:"-"`
	Operator  string `json:"operator,omitempty"`
	Value     string `json:"value,omitempty"`
}

type Field struct {
	Dirty       bool   `json:"dirty,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Html        string `json:"html,omitempty"`
	Name        string `json:"name,omitempty"`
	Value       string `json:"value,omitempty"`
	FieldID     string `json:"fieldId,omitempty"`
	FullFieldID string `json:"fullFieldId"`
}

type Error struct {
	ErrorCode      string `json:"errorCode,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
	HasError       bool   `json:"hasError,omitempty"`
	HTTPStatusCode string `json:"httpStatusCode"`
}

type Link struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// NewClient returns a Pointer to an Instance of a Cherwell Client
func NewClient(user, password, clientID, baseURI, auth_mode, grant_type string) *Client {
	return &Client{
		User:       user,
		Password:   password,
		ClientID:   clientID,
		BaseURI:    baseURI,
		Auth_mode:  auth_mode,
		Grant_Type: grant_type,
	}
}

// Login authenticates with the Cherwell Server and retreives an AccessToken, saves it
// in the Instance of the Client and returns a Pointer to the Instance of the Client
func (cl *Client) Login() *Client {
	params := url.Values{}
	params.Add("grant_type", cl.Grant_Type)
	params.Add("client_id", cl.ClientID)
	params.Add("username", cl.User)
	if cl.Password != "" {
		params.Add("password", cl.Password)
	}
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", (cl.BaseURI + "token?auth_mode=" + cl.Auth_mode), body)
	if err != nil {
		fmt.Printf("\nLogin failed @ NewRequest: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("\nLogin failed @ DoRequest: %v", err)
		return nil
	}
	defer resp.Body.Close()
	err = unJson(resp.Body, &cl)
	if err != nil {
		fmt.Printf("\nLogin failed @ unJson: %v", err)
		return cl
	}
	return cl
}

// keepAlive sends the RefreshToken to the Cherwell Server and retreives a new AccessToken and returns a bool
func (cl *Client) keepAlive() bool {
	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("client_id", cl.ClientID)
	params.Add("refresh_token", cl.Refresh_token)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", (cl.BaseURI + "token"), body)
	if err != nil {
		fmt.Printf("\nLogin failed @ NewRequest: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("\nLogin failed @ DoRequest: %v", err)
		return false
	}
	defer resp.Body.Close()
	err = unJson(resp.Body, &cl)
	if err != nil {
		fmt.Printf("\nLogin failed @ unJson: %v", err)
		return false
	}
	return true
}

// validateToken checks if the AccessToken is expired or still valid and returns a bool
func (cl *Client) validateToken() bool {
	t, err := time.Parse(time.RFC1123, cl.Expires)
	if err != nil {
		fmt.Printf("\nCould not parse Timestamp of AccessToken: %v", err)
		return false
	}
	return time.Now().Before(t.Add(time.Minute * -5))
}

// request creates, enriches and submits a HTTP-Request to the Cherwell Server
// and unmarshales the HTTP-Response to a given Output-Object
func (cl *Client) request(method, uri string, input, output interface{}) {
	if !cl.validateToken() {
		if !cl.keepAlive() {
			cl = cl.Login()
		}
	}

	req := &http.Request{}
	var err error

	if input == nil {
		params := url.Values{}
		params.Add("client_id", cl.ClientID)
		body := strings.NewReader(params.Encode())
		req, err = http.NewRequest(strings.ToUpper(method), uri, body)
		if err != nil {
			fmt.Printf("\nFailed to create Request: %v\nMethod: %v\nURI: %v", err, method, uri)
			return
		}
	} else {
		payloadBytes, err := json.Marshal(input)
		if err != nil {
			fmt.Printf("\nFailed to get BusinessObjectTemplate @ Marshal: %v", err)
			return
		}
		body := bytes.NewReader(payloadBytes)

		req, err = http.NewRequest(strings.ToUpper(method), uri, body)
		if err != nil {
			fmt.Printf("\nFailed to create Request: %v\nMethod: %v\nURI: %v", err, method, uri)
			return
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", ("Bearer " + cl.Access_token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("\nFailed to send Request: %v\nMethod: %v\nURI: %v", err, method, uri)
		return
	}
	defer resp.Body.Close()

	err = unJson(resp.Body, &output)
	if err != nil {
		fmt.Printf("\nFailed to unmarshal Response: %v\nMethod: %v\nURI: %v\nResponse: %v", err, method, uri, resp)
		return
	}
}

// GetBusinessObjectByDisplayName retreives a Cherwell BusinessObject by given DisplayName and returns it
func (cl *Client) GetBusinessObjectByDisplayName(displayName string) *BusinessObject {
	res := []BusinessObject{}
	uri := cl.BaseURI + getBusObSummariesAllURI
	cl.request("get", uri, nil, &res)
	for _, b := range res {
		//fmt.Printf("\nBusOB: %v", b.DisplayName)
		if b.DisplayName == displayName {
			return &b
		}
		for _, c := range b.GroupSummaries {
			if c.DisplayName == displayName {
				return &c
			}
		}
	}

	return nil
}

// GetBusinessObjectByBusObID retreives a Cherwell BusinessObject by given BusObID and returns it
func (cl *Client) GetBusinessObjectByBusObID(busObID string) *BusinessObject {
	res := []BusinessObject{}
	uri := cl.BaseURI + getBusObSummariesAllURI

	val := make(map[string]string)
	val["busobid"] = busObID

	uri = formatURI(uri, val)
	cl.request("GET", uri, nil, &res)

	for _, b := range res {
		if b.BusObID == busObID {
			return &b
		}
		for _, c := range b.GroupSummaries {
			if c.BusObID == busObID {
				return &c
			}
		}
	}
	return nil
}

// GetBusinessObjectRecordByPublicID retreives a Cherwell BusinessObjectRecord by given PublicID and returns it
func (bo *BusinessObject) GetBusinessObjectRecordByPublicID(cl *Client, publicID string) *BusinessObjectRecord {
	if bo == nil {
		fmt.Printf("\nBusinessObject cannot be nil")
		return nil
	}
	res := BusinessObjectRecord{}
	uri := cl.BaseURI + getBusObRecByPublicIdURI

	val := make(map[string]string)
	val["busobid"] = bo.BusObID
	val["busobpublicid"] = publicID

	uri = formatURI(uri, val)
	cl.request("GET", uri, nil, &res)

	rec := *res.processFields()
	return &rec
}

// GetBusinessObjectRecordByRecID retreives a Cherwell BusinessObjectRecord by given RecID and returns it
func (bo *BusinessObject) GetBusinessObjectRecordByRecID(cl *Client, recID string) *BusinessObjectRecord {
	if bo == nil {
		fmt.Printf("\nBusinessObject cannot be nil")
		return nil
	}
	res := BusinessObjectRecord{}
	uri := cl.BaseURI + getBusObRecByRecIdURI

	val := make(map[string]string)
	val["busobid"] = bo.BusObID
	val["busobrecid"] = recID

	uri = formatURI(uri, val)

	cl.request("GET", uri, nil, &res)
	rec := *res.processFields()
	return &rec
}

// GetBusinessObjectTemplate retreives a Cherwell BusinessObjectTemplate of a given BusinessObject and returns it
func (bo *BusinessObject) GetBusinessObjectTemplate(cl *Client) *BusinessObjectTemplate {
	if bo == nil {
		fmt.Printf("\nBusinessObject cannot be nil")
		return nil
	}
	res := BusinessObjectTemplate{}
	uri := cl.BaseURI + getBusObTemplateURI
	query := struct {
		BusObID         string `json:"busObId"`
		IncludeAll      bool   `json:"includeAll"`
		IncludeRequired bool   `json:"includeRequired"`
	}{
		BusObID:         bo.BusObID,
		IncludeAll:      true,
		IncludeRequired: true,
	}
	cl.request("POST", uri, &query, &res)
	return &res
}

// SearchBusinessObjectRecord retreives a Cherwell BusinessObjectRecord by Search-Request with Filters and returns it
func (bo *BusinessObject) SearchBusinessObjectRecord(cl *Client, filters ...[]string) *BusinessObjectRecord {
	if bo == nil {
		fmt.Printf("\nBusinessObject cannot be nil")
		return nil
	}
	res := SearchResult{}
	uri := cl.BaseURI + getSearchResultsURI
	filter := []Filter{}
	fields := bo.GetBusinessObjectTemplate(cl)
	for _, f := range fields.Fields {

		for _, fi := range filters {
			if len(fi) != 3 {
				fmt.Printf("\nFilter invalid!\nWant: ['FieldDisplayName','Operator','Value']\nGot:%v", f)
				return nil
			}
			if f.DisplayName == fi[0] {
				filter = append(filter, Filter{
					FieldID:  f.FieldID,
					Operator: fi[1],
					Value:    fi[2],
				})
			}
		}
	}

	query := Search{
		BusObID:          bo.BusObID,
		Filters:          filter,
		IncludeAllFields: true,
	}
	cl.request("POST", uri, &query, &res)
	rec := res.BusinessObjects[0].processFields()
	return rec
}

// SearchMultipleBusinessObjectRecord retreives multiple Cherwell BusinessObjectRecords by Search-Request with Filters and returns them
func (bo *BusinessObject) SearchMultipleBusinessObjectRecords(cl *Client, filters ...[]string) *[]BusinessObjectRecord {
	if bo == nil {
		fmt.Printf("\nBusinessObject cannot be nil")
		return nil
	}
	res := SearchResult{}
	uri := cl.BaseURI + getSearchResultsURI
	filter := []Filter{}
	fields := bo.GetBusinessObjectTemplate(cl)
	for _, f := range fields.Fields {

		for _, fi := range filters {
			if len(fi) != 3 {
				fmt.Printf("\nFilter invalid!\nWant: ['FieldDisplayName','Operator','Value']\nGot:%v", f)
				return nil
			}
			if f.DisplayName == fi[0] {
				filter = append(filter, Filter{
					FieldID:  f.FieldID,
					Operator: fi[1],
					Value:    fi[2],
				})
			}
		}
	}

	query := Search{
		BusObID:          bo.BusObID,
		Filters:          filter,
		IncludeAllFields: true,
	}
	cl.request("POST", uri, &query, &res)
	var records []BusinessObjectRecord
	for _, r := range res.BusinessObjects {
		records = append(records, *r.processFields())
	}
	return &records
}

// GetBusinessObjectSchema retreives a Cherwell BusinessObjectSchema of a given BusinessObject and returns it
func (bo *BusinessObject) GetBusinessObjectSchema(cl *Client) *BusinessObjectSchema {
	if bo == nil {
		fmt.Printf("\nBusinessObject cannot be nil")
		return nil
	}
	var schema BusinessObjectSchema
	uri := cl.BaseURI + strings.Replace(getBusObSchemaURI, "$", bo.BusObID, 1)
	cl.request("GET", uri, nil, &schema)
	return &schema
}

// processFields enriches a BusinessObjectRecord with FieldValues
// to make access to the Values of Fields easier and returns it
func (rec *BusinessObjectRecord) processFields() *BusinessObjectRecord {
	fields := make(map[string]interface{})
	for _, f := range rec.Fields {
		fields[f.DisplayName] = f.Value
	}
	rec.FieldValues = fields
	return rec
}

// SaveBusinessObjectRecord commits the Changes in FieldValues to Fields and saves the Cherwell BusinessObjectRecord and returns it
func (rec *BusinessObjectRecord) SaveBusinessObjectRecord(cl *Client) *BusinessObjectRecord {
	if rec == nil {
		fmt.Printf("\nBusinessObjectRecord cannot be nil")
		return nil
	}
	saveResp := BusinessObjectRecord{}
	uri := cl.BaseURI + saveBusObRecURI

	rec.Persist = true

	for i, f := range rec.Fields {
		if f.Value != rec.FieldValues[f.DisplayName] {
			rec.Fields[i].Value = fmt.Sprint(rec.FieldValues[f.DisplayName])
			rec.Fields[i].Dirty = true
			fmt.Printf("\nChanged: %v", rec.Fields[i])
		}
	}
	cl.request("post", uri, &rec, &saveResp)
	return &saveResp
}

// DeleteBusinessObjectRecord deletes a Cherwell BusinessObjectRecord and returns the Response
func (rec *BusinessObjectRecord) DeleteBusinesObjectRecord(cl *Client) *BusinessObjectRecord {
	if rec == nil {
		fmt.Printf("\nBusinessObjectRecord cannot be nil")
		return nil
	}
	res := BusinessObjectRecord{}
	uri := cl.BaseURI + deleteBusObRecURI
	val := make(map[string]string)
	val["busobid"] = rec.BusObID
	val["busobrecid"] = rec.BusObRecID
	uri = formatURI(uri, val)
	cl.request("DELETE", uri, nil, &res)
	return &res
}

// GetRelatedBusinessObjects retreives all Cherwell BusinessObjectRecords, by Name of the Relationship,
// related to a given BusinessObjectRecord and returns them
func (rec *BusinessObjectRecord) GetRelatedBusinessObjects(cl *Client, relationshipName string) *[]BusinessObjectRecord {
	if rec == nil {
		fmt.Printf("\nBusinessObjectRecord cannot be nil")
		return nil
	}
	res := RelatedBusinessObjects{}
	uri := cl.BaseURI + getRelatedBusObURI
	relID := cl.GetBusinessObjectByBusObID(rec.BusObID).GetBusinessObjectSchema(cl).GetRelationshipID(relationshipName)

	val := make(map[string]string)
	val["busobid"] = rec.BusObID
	val["busobrecid"] = rec.BusObRecID
	val["relationshipid"] = relID
	uri = formatURI(uri, val)
	cl.request("GET", uri, nil, &res)

	records := []BusinessObjectRecord{}

	for _, r := range res.RelatedBusinessObjects {
		records = append(records, *r.processFields())
	}
	return &records
}

// LinkBusinessObjectRecord links the Cherwell BusinessObjectRecord to a given Child BusinessObjectRecord
func (rec *BusinessObjectRecord) LinkBusinessObjectRecord(cl *Client, childRec *BusinessObjectRecord, relationshipName string) *Error {
	if rec == nil {
		fmt.Printf("\nBusinessObjectRecord cannot be nil")
		return nil
	}
	res := Error{}
	busOb := cl.GetBusinessObjectByBusObID(rec.BusObID)
	if busOb == nil {
		fmt.Printf("\nNot found: BusinessObject: %v", rec.BusObID)
		return nil
	}

	sch := busOb.GetBusinessObjectSchema(cl)
	if sch == nil {
		fmt.Printf("\nNot found: BusinessObjectSchema: %v", busOb.DisplayName)
		return nil
	}
	relID := sch.GetRelationshipID(relationshipName)
	uri := cl.BaseURI + linkBusObRecURI

	val := make(map[string]string)
	val["busobid"] = rec.BusObID
	val["busobrecid"] = rec.BusObRecID
	val["relationshipid"] = relID
	val["childbusobid"] = childRec.BusObID
	val["childbusobrecid"] = childRec.BusObRecID

	uri = formatURI(uri, val)
	cl.request("GET", uri, nil, &res)
	return &res
}

// UnlinkBusinessObjectRecord unlinks the Cherwell BusinessObjectRecord from a given Child BusinessObjectRecord
func (rec *BusinessObjectRecord) UnlinkBusinessObjectRecord(cl *Client, childRec *BusinessObjectRecord, relationshipName string) *Error {
	if rec == nil {
		fmt.Printf("\nBusinessObjectRecord cannot be nil")
		return nil
	}
	res := Error{}
	busOb := cl.GetBusinessObjectByBusObID(rec.BusObID)
	if busOb == nil {
		fmt.Printf("\nNot found: BusinessObject: %v", rec.BusObID)
		return nil
	}

	sch := busOb.GetBusinessObjectSchema(cl)
	if sch == nil {
		fmt.Printf("\nNot found: BusinessObjectSchema: %v", busOb.DisplayName)
		return nil
	}
	relID := sch.GetRelationshipID(relationshipName)
	uri := cl.BaseURI + unlinkBusObRecURI

	val := make(map[string]string)
	val["busobid"] = rec.BusObID
	val["busobrecid"] = rec.BusObRecID
	val["relationshipid"] = relID
	val["childbusobid"] = childRec.BusObID
	val["childbusobrecid"] = childRec.BusObRecID

	uri = formatURI(uri, val)
	cl.request("DELETE", uri, nil, &res)
	return &res
}

// GetRelationshipID links the Cherwell BusinessObjectRecord to a given Child BusinessObjectRecord
func (sch *BusinessObjectSchema) GetRelationshipID(relationshipName string) string {
	if sch == nil {
		fmt.Printf("\nBusinessObjectSchema cannot be nil")
		return ""
	}
	var relID string
	for _, r := range sch.Relationships {
		if r.DisplayName == relationshipName {
			relID = r.RelationshipID
			return relID
		}
	}
	return relID
}

// GetTeamMembers retreives all Contacts linked to a Cherwell BusinessObjectRecord of the
// BusinessObject "OrganizationalUnit" with the Type "Team" and the given Name as Name
func (cl *Client) GetTeamMembers(teamName string) *[]BusinessObjectRecord {
	var teamRecID string
	bo := cl.GetBusinessObjectByDisplayName("Organisationseinheit")
	teams := bo.SearchMultipleBusinessObjectRecords(cl, []string{
		"Typ", "EQ", "Team"},
		[]string{"Voller Name", "EQ", teamName},
	)

	for _, t := range *teams {
		if t.BusObPublicID == teamName {
			teamRecID = t.BusObRecID
		}
		fmt.Printf("\nTeam: %v | RecID: %v", t.BusObPublicID, t.BusObRecID)
	}

	members := bo.GetBusinessObjectRecordByRecID(cl, teamRecID).GetRelatedBusinessObjects(cl, "Organisation Unit Links Contacts Member")
	return members
}

// unJson unmarshals a given io.ReadCloser to a given interface
func unJson(input io.ReadCloser, data interface{}) error {
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

// formatURI replaces the Placeholders in a given URI with the given Values
func formatURI(uri string, values map[string]string) string {
	if val, ok := values["busobid"]; ok {
		uri = strings.Replace(uri, busobid, val, 1)
	}
	if val, ok := values["busobrecid"]; ok {
		uri = strings.Replace(uri, busobrecid, val, 1)
	}
	if val, ok := values["busobpublicid"]; ok {
		uri = strings.Replace(uri, busobpublicid, val, 1)
	}
	if val, ok := values["childbusobid"]; ok {
		uri = strings.Replace(uri, childbusobid, val, 1)
	}
	if val, ok := values["childbusobrecid"]; ok {
		uri = strings.Replace(uri, childbusobrecid, val, 1)
	}
	if val, ok := values["relationshipid"]; ok {
		uri = strings.Replace(uri, relationshipid, val, 1)
	}
	return uri
}
