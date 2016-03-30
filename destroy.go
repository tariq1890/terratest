package terratest

import (
	"fmt"

	"github.com/gruntwork-io/terratest/log"
	"github.com/gruntwork-io/terratest/terraform"
)

// Destroy both destroys all the given elements of the RandomResourceCollection and calls terraform destroy.
func Destroy(options *TerratestOptions, rand *RandomResourceCollection) (string, error) {
	logger := log.NewLogger(options.TestName)

	err := rand.DestroyResources()
	if err != nil {
		return "", fmt.Errorf("Failed to destroy random resource collection: %s", err.Error())
	}

	logger.Println("Running terraform destroy...")
	output, err := destroyHelper(options, options.GetTfStateFileName())
	if err != nil {
		return output, fmt.Errorf("Failed to terraform destroy: %s", err.Error())
	}

	return output, nil
}

// Helper function that calls terraform destroy
func destroyHelper(options *TerratestOptions, remoteStateS3ObjectName string) (string, error) {
	logger := log.NewLogger(options.TestName)
	output, err := terraform.DestroyAndGetOutput(options.TemplatePath, options.Vars, logger)
	if err != nil {
		logger.Printf(`Failed to terraform destroy.
** WARNING ** Terraform destroy has failed which means you must manually delete any resources created by the "terraform apply" run.
Test Name: %s
Terraform Template Path: %s
AWS Region: <scroll up to see it>
Remote State Location: s3://%s/%s
Official Error Message: %s
`, options.TemplatePath, options.TestName, options.TfRemoteStateS3BucketName, remoteStateS3ObjectName, err.Error())
		return output, err
	}

	return output, nil
}