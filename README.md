# gocherwell
***gocherwell*** is a go-Module for communication with the REST-API of [Cherwell](https://cherwell.com "https://cherwell.com").

[![Go Reference](https://pkg.go.dev/badge/github.com/itsscb/gocherwell.svg)](https://pkg.go.dev/github.com/itsscb/gocherwell)
## Usage
### Authentication
```
var user = "TESTUSER"
var password = "p4$$w0rd"
var clientID = "0000-1111-2222-3333-4444"
var baseURI = "https://example.com/CherwellAPI/"
var auth_mode = "Internal"
var grant_type = "password"

cl := gocherwell.NewClient(
    user,
    password,
    clientID,
    baseURI,
    auth_mode,
    grant_type,
).Login()
```

### Get BusinessObjects
#### By DisplayName
Example returns the BusinessObject with the ***DisplayName*** *Configuration Item*
```
bo := cl.GetBusinessObjectByDisplayName("Configuration Item")
```
#### By BusObID
Example returns the BusinessObject with the ***BusObID*** *012345678910abcdefghijklmnop*
```
bo := cl.GetBusinessObjectByBusObID("012345678910abcdefghijklmnop")
```

### Get BusinessObjectRecords
#### By PublicID
Example returns the BusinessObjectRecord of the ***Configuration Item*** with the ***PublicID*** *NOTEBOOK001*
```
rec := bo.GetBusinessObjectRecordByPublicID(cl, "NOTEBOOK001")
```
#### By RecID
Example returns the BusinessObjectRecord of the ***Configuration Item*** with the ***RecID*** *abcdefghijklmnop012345678910*
```
rec := bo.GetBusinessObjectRecordByRecID(cl, "abcdefghijklmnop012345678910")
```
#### By Search
##### Single Record (First Hit)
Example returns the first Hit of BusinessObjectRecords with the ***AssetName*** *NOTEBOOK001* of ***Type*** *Notebook* with the ***Status*** *Active* 
```
rec := bo.SearchObjectRecord(cl, []string{
    "AssetName", "EQ", "NOTEBOOK001",
    },
    []string{
        "Type","EQ","Notebook",
    },
    []string{
        "Status","EQ","Active"
    },
)
```
##### Multiple Records
Example returns all BusinessObjectRecords of ***Type*** *Notebook* with the ***Status*** *Active* 
```
records := bo.SearchObjectRecord(cl, []string{
    []string{
        "Type","EQ","Notebook",
    },
    []string{
        "Status","EQ","Active"
    },
)
```
### BusinessObjectRecord Actions
#### New BusinessObjectRecord
This method of ***BusinessObject*** takes all given ***Field***s and creates a BusinessObjectRecord with them
```
resp := rec.SaveBusinessObjectRecord(cl)
```

#### Save BusinessObjectRecord
This method of ***BusinessObjectRecord*** goes over all ***.FieldValues*** and commits the changed fields to ***.Fields*** and sets ***Dirty*** to *True*
```
rec := bo.NewBusinessObjectRecord(cl, []gocherwell.Field{
    gocherwell.Field{
        DisplayName:    "AssetName",
        Value:          "NOTEBOOK001",
    },
    gocherwell.Field{
        DisplayName:    "AssetType",
        Value:          "Notebook",
    },
    gocherwell.Field{
        DisplayName:    "Status",
        Value:          "Active",
    },
})
```

#### Delete BusinessObjectRecord
This method of ***BusinessObjectRecord*** deletes the executing ***BusinessObjectRecord***
```
resp := rec.DeleteBusinessObjectRecord(cl)
```

#### Link BusinessObjectRecords
This method of ***BusinessObjectRecord*** links the executing and the given ***BusinessObjectRecord***


The following Example links the ***Configuration Item*** *NOTEBOOK001* to the ***Note** with the ***PublicID*** *NOTE-1234*
```
child := cl.GetBusinessObject("Note").GetBusinessObjectRecordByPublicID("NOTE-1234")
resp := rec.LinkBusinessObjectRecord(cl, child, "Configuration Item Links Note")
```

#### Unlink BusinessObjectRecords
This method of ***BusinessObjectRecord*** unlinks the executing and the given ***BusinessObjectRecord***


The following Example unlinks the ***Configuration Item*** *NOTEBOOK001* to the ***Note** with the ***PublicID*** *NOTE-1234*
```
child := cl.GetBusinessObject("Note").GetBusinessObjectRecordByPublicID("NOTE-1234")
resp := rec.UninkBusinessObjectRecord(cl, child, "Configuration Item Links Note")
```
