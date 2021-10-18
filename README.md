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
