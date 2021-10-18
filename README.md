# gocherwell
Go module for utilization of the Cherwell API

## Info
This build was tested with Cherwell API v9.7.
See https://help.cherwell.com/bundle/cherwell_rest_api_970_help_only/page/content/system_administration/rest_api/csm_rest_api_landing_page.html

## Usage
### Create a new Instance
````
client := gocherwell.NewClient("USERNAME", "PASSWORD", "CLIENT_ID","BASE_URI","AUTH_MODE","GRANT_TYPE")
````
### Authenticate / Get Token
````
err := client.Login()
if err != nil {
    fmt.Print(err.Error())
    return err
}
````
### Get BusinessObject Summaries / Get all BusinessObjects
````
bos, err := client.GetAllBusOb()
if err != nil {
    fmt.Print(err.Error())
    return err
}
````
### Get BusinessObject by DisplayName
````
bo, err := client.GetBusOb("DISPLAYNAME_OF_BUSINESSOBJECT")
if err != nil {
    fmt.Print(err.Error())
    return err
}
````
### Get BusinessObjectRecord by PublicID
````
record, err := client.GetObRecByPublicID(bo, "PUBLICID_OF_RECORD")
if err != nil {
    fmt.Print(err.Error())
    return err
}
````
### Get BusinessObjectRecord by RecordID
````
record, err := client.GetObRecByRecID(bo, "RECORDID_OF_RECORD")
if err != nil {
    fmt.Print(err.Error())
    return err
}
````
### Get BusinessObjectTemplate
````
template, err := client.GetBusObTemplate("BUSINESSOBJECTID")
if err != nil {
    fmt.Print(err.Error())
    return err
}
````
### Get QuickSearch Results
````
results, err := client.QuickSearch("SEARCHTEXT")
if err != nil {
    fmt.Print(err.Error())
    return err
}
````
### Get Search Results
````
filter, err := client.NewFilter(bo, []gocherwell.Filter{
    {
        FieldName: "AssetName",
        Operator: "EQ",
        FieldValue: "NAME_OF_ASSET",
})
if err != nil {
    fmt.Print(err.Error())
    return err
}
record, err := client.Search(bo, filter)
if err != nil {
    fmt.Print(err.Error())
    return err
}
````
### Save(Update) / Create BusinessObjectRecord
````
record := bo.NewObRec([]gocherwell.Field{
    {
        Name: "AssetName",
        Value: "NAME_OF_NEW_ASSET",
    },
})
response, err := client.SaveObRec(record) 
if err != nil {
    fmt.Print(err.Error())
    return err
}
````
### Delete BusinessObjectRecord
````
record, err := client.DeleteObRec(record)
if err != nil {
    fmt.Print(err.Error())
    return err
}
````