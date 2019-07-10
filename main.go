package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	resourceGroupName     = "GoVMQuickstart"
	resourceGroupLocation = "eastus"
	deploymentName        = "VMDeployQuickstart"
	templateFile          = "template1.json"
	parametersFile        = "parameters1.json" // not used in this example
)

// Information loaded from the authorization file to identify the client
type clientInfo struct {
	SubscriptionID string
	VMPassword     string
}

var (
	ctx        = context.Background()
	clientData clientInfo
	authorizer autorest.Authorizer
)

// Authenticate with the Azure services using file-based authentication
func init() {
	var err error
	authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Failed to get OAuth config: %v", err)
	}

	authInfo, err := readJSON(os.Getenv("AZURE_AUTH_LOCATION"))
	clientData.SubscriptionID = (*authInfo)["subscriptionId"].(string)
	clientData.VMPassword = (*authInfo)["clientSecret"].(string)
}

func main() {
	log.Printf("Starting deployment: %s", deploymentName)
	result, err := createDeployment()
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	if result.Name != nil {
		log.Printf("Completed deployment %v: %v", deploymentName, *result.Properties.ProvisioningState)
	} else {
		log.Printf("Completed deployment %v (no data returned to SDK)", deploymentName)
	}
}

// Create the deployment
func createDeployment() (deployment resources.DeploymentExtended, err error) {
	template, err := readJSON(templateFile)
	if err != nil {
		return
	}
	// params, err := readJSON(parametersFile)
	// if err != nil {
	// 	return
	// }
	// (*params)["vm_password"] = map[string]string{
	// 	"value": clientData.VMPassword,
	// }

	deploymentsClient := resources.NewDeploymentsClient(clientData.SubscriptionID)
	deploymentsClient.Authorizer = authorizer

	location := "eastus"

	deploymentFuture, err := deploymentsClient.CreateOrUpdateAtSubscriptionScope(
		ctx,
		deploymentName,
		resources.Deployment{
			Location: &location,
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: nil,
				Mode:       resources.Incremental,
			},
		},
	)
	if err != nil {
		return
	}
	err = deploymentFuture.Future.WaitForCompletionRef(ctx, deploymentsClient.BaseClient.Client)
	if err != nil {
		return
	}
	return deploymentFuture.Result(deploymentsClient)
}

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	contents := make(map[string]interface{})
	json.Unmarshal(data, &contents)
	return &contents, nil
}
