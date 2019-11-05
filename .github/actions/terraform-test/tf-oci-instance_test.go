// Copyright (c) 2019 Oracle and/or its affiliates,  All rights reserved.

package test

import (
	"testing"
	"time"
	"os"
	//"github.com/gruntwork-io/terratest/modules/oci"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	//"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"gotest.tools/assert"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"fmt"
	"strings"
	"github.com/gruntwork-io/terratest/modules/retry"
	"io/ioutil"
	"crypto/tls"
)

// An example of how to test the Terraform module in examples/terraform-http-example using Terratest.
func TestOCIComputeInstanceTFModule(t *testing.T) {
	t.Parallel()

	// Specify the text the OCI Instance will return when we make HTTP requests to it.
	instanceText := "Hello World from Apache running on test-cotud"

	
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: os.Getenv("TF_ACTION_WORKING_DIR"),
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			// Oracle Cloud Infrastructure Authentication details
			"tenancy_ocid": os.Getenv("TF_VAR_tenancy_id"),
			"user_ocid": os.Getenv("TF_VAR_user_id"),
			"fingerprint": os.Getenv("TF_VAR_fingerprint"),
			"private_key_path": os.Getenv("TF_VAR_private_key_path"),

			// Region
			"region": "uk-london-1",
			"availability_domain": "uFjs:UK-LONDON-1-AD-1",
			
			// Compartment
			"compartment_ocid": "ocid1.compartment.oc1..aaaaaaaacnmuyhg2mpb3z6v6egermq47nai3jk5qaoieg3ztinqhamalealq",
			
			// Compute Instance Configurations
			"instance_display_name": "test-cotud-techflow-11",
			"source_ocid": "ocid1.image.oc1.uk-london-1.aaaaaaaa32voyikkkzfxyo4xbdmadc2dmvorfxxgdhpnk6dw64fa3l4jh7wa",
			"vcn_id": "ocid1.vcn.oc1.uk-london-1.aaaaaaaahdze4mulmvmym4xvrj22cb2zq72tm6qul2oarfnpub6sxwv77owq",
			"cidr" : "10.20.10.0/24",
			"ssh_authorized_keys": os.Getenv("TF_VAR_ssh_authorized_keys"),
			"ssh_private_key_path": os.Getenv("TF_VAR_ssh_private_key_path"),
			"instance_count": "1",
			
			// Storage Volume Configurations
			"block_storage_sizes_in_gbs": "[\"50\"]",
		},
	}
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)
	
	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	instanceURL := terraform.Output(t, terraformOptions, "instance_url")

	// Run `terraform output` to get the values of output variables
	hasPublicIP := terraform.Output(t, terraformOptions, "public_ip")

	keyPair := GetKeyPair(os.Getenv("TF_VAR_ssh_private_key_path"), os.Getenv("TF_VAR_ssh_authorized_key_path"))

	testSSHToPublicHost(t, terraformOptions, keyPair)
	
	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}
	
	// Verify that we get back a 200 OK with the expected instanceText
	http_helper.HttpGetWithRetry(t, instanceURL, &tlsConfig, 200, instanceText, maxRetries, timeBetweenRetries)

	// Verify we're getting back the outputs we expect - that the instance has been assigned with a public IP
	assert.Assert(t, len(hasPublicIP) != 0)
}

func testSSHToPublicHost(t *testing.T, terraformOptions *terraform.Options, keyPair *ssh.KeyPair) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_instance_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "opc",
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair,
		SshUserName: "opc",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := fmt.Sprintf("SSH to public host %s", publicInstanceIP)

	// Run a simple echo command on the server
	expectedText := "Hello, World"
	command := fmt.Sprintf("echo -n '%s'", expectedText)

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})

	// Run a command on the server that results in an error,
	expectedText = "Hello, World"
	command = fmt.Sprintf("echo -n '%s' && exit 1", expectedText)
	description = fmt.Sprintf("SSH to public host %s with error command", publicInstanceIP)

	// Verify that we can SSH to the Instance, run the command and see the output
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {

		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err == nil {
			return "", fmt.Errorf("Expected SSH command to return an error but got none")
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}

func GetKeyPair(privfile string, pubfile string) *ssh.KeyPair {

	// read keys from file
	_, errPriv := os.Stat(privfile)
	_, errPub := os.Stat(pubfile)
	if (errPriv == nil && errPub == nil) {
		priv, errPriv := ioutil.ReadFile(privfile)
		if errPriv != nil {
			fmt.Errorf("Failed to read priv file - %s", errPriv)
		}
		pub, errPub := ioutil.ReadFile(pubfile)
		if errPub != nil {
			fmt.Errorf("Failed to read pub file - %s", errPub)
		}
		return &ssh.KeyPair{PublicKey: string(pub), PrivateKey: string(priv)}
	}
	return nil
}