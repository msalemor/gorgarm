# Deploy a Subscription Scope Template in Azure using Go

This Go code deploys a Subscription scope ARM template

## Requirements

- create a service principal and export it to a file using the Azure CLI (this can be done from the Cloud Shell in the portal)
```
az ad sp create-for-rbac --sdk-auth > azure.auth
```
- set the AZURE_AUTH_LOCATION session variable to the sp file
```
SET AZURE_AUTH_LOCATION=c:\go\azure.auth
```
or
```
export AZURE_AUTH_LOCATION=/user/home/go/azure.auth
```
- Run the application
```
go run main.go
```

## Reference

This code comes mainly from:

- https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-qs-vm
