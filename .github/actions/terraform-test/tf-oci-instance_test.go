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
)

// An example of how to test the Terraform module in examples/terraform-http-example using Terratest.
func TestOCIComputeInstanceTFModule(t *testing.T) {
	t.Parallel()
	// A unique ID we can use to namespace resources so we don't clash with anything already in the OCI account or
	// tests running in parallel
	//uniqueID := random.UniqueId()
	// Give this OCI Instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the OCI account.
	//instanceName := fmt.Sprintf("terratest-http-example-%s", uniqueID)
	
	// Specify the text the OCI Instance will return when we make HTTP requests to it.
	instanceText := "Hello World from Apache running on test-cotud"
	
	// Pick an OCI region
	//ociRegion := "uk-london-1"
	
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "/Users/cotudor/Work/Tasks-FY20/SDF/SDF-20-EvaluateModules/CoreModules/ComputeInstanceTFModule/terraform-oci-compute-instance/examples/instance_default",
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			// Oracle Cloud Infrastructure Authentication details
			"tenancy_ocid": "ocid1.tenancy.oc1..aaaaaaaaa3qmjxr43tjexx75r6gwk6vjw22ermohbw2vbxyhczksgjir7xdq",
			"user_ocid": "ocid1.user.oc1..aaaaaaaaoj3k4zp24w2aryx7z46gzhw2yxf2lrvcj3om6cikwb6k4c2p74sq",
			"fingerprint":  "25:93:69:40:2f:5b:d2:25:0e:eb:f3:41:ea:cb:18:02",
			"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAx7CWDt6MCoM4JMQbvrTaRGlHcBJjtJ94iGxJqylkVjeNit11\nSgQcV0vq+evL0GsW52TptNPlFvboZC8thGEtGQxe4NT70OHkjvnbJBiFA+5055iN\nMSkLZ5Sex0+BRwLIfOay15WqcFlgFK75Oy1uO0EY/73HVXQFttZGqZc4Sc2b85Pg\nbelSv8cn/WGGp1Wqcp+m47C7zTeYoIYyecWLr/mF6dTFhiqQwTU6301bmlsAPpjF\nZBo0411LH68lRo5Rkzkf6ulxZB1EXkNP16kyOgRj8pp+nxpMQt5kdaJPYYsUFXhM\nnOs0lvQ0TpDAJZ8xfOUAlmBFHoSUWx1x07XsxwIDAQABAoIBAQCDuhlC9curo6QF\nnNgwlVzmhAx8AaWEIS6Bz+1I26U0urDEShJ1IQERFSOMed+ZPQt+2TLR9nXJEFd8\nvyJnkTMOsvpjlhwHTvxW1LuatBIK5eJ4ZEm/hDPhwEh9chT1owBeZFDVpgUo18xp\n2ILQZsewjuDO04IK/N5IHlG+zeJzy+NpEg9vtPN1A3aDomSMqy+DrvuTFpH6CNtM\nNrDTqDccQbnF0cZSXD1L6mw4kWol8jAsSyGDf7uOOSyNT57ky80mRDBu9CuA1N9j\npnJCZaWBrQXpKco5lBLSYXJ1JddzvSSim8QV+rYd0+004ad7qojigBQEmGLOAT81\nhZkSQCwBAoGBAPKlBVV5+AyPuf9QG7QaAceMwHWIhjxZemJ6o7QARKFc2Dc1CYQO\nAutZV4JXGhrQLV+wuHY2JG5AXyVOxolWsW14BYxx0+Oq6U0p6OOetYi3LGiLCT4N\nwuAp4lMLEPfPNht7SFfq7ZannaaREzhr4LQXKlGMEszZsZ+/KATyDwJxAoGBANKu\nT609cDaXL/TnfOpWqV5bXTUiNjdZDT05S62Ynvr5KkXBa97pE0cuZ7cwGgjxC5Pe\nDteGcbUe/Eikx8FVvqfuYRvKuhAoPdFZ1miR3sWDHOlFjkuJxaGvUbBDlBwcKHx9\nAOxAruhCRmG5/dYAIJ5KlEYJF5xsZz9SmQ7A6Q63AoGAEuaZ2MOsd8YGVgX2cnwI\nIXQsVbtxwWey6dLlx5KxxeQGj55ZBGlW/uAxudxxEx+eOTL523NyOQhYoP5W5sHT\nBlTwEbWYLMbWb4VRN9HYEDM8iVQzPxsxT+bTU2asRrFkZJWg2ABby315AU2RsrZs\nhXq8eCeyGzTl6iyowGHem3ECgYBdfVCKFctnzitPyDGcY5yA7JYt7+KTKQdA9d3p\nSOKziEID9lMB9ffCDIultMi40w5KLa30YgqvTvKw4b5qwrv6FUQuawWqCdF0xyLo\nAGMUzpvTwDPmvVpf50aeqz5cQvqMU4RHUmTLWC2XTEuh3SicVYf6lCpQFaKzbNnS\nvDQvfQKBgQCAqS2DriooFdQpcblAyN2twGBQbGUYw6FUA0N9B8jMzxoq+B3P6jku\nDmhhVv0pf0pgwlmWqNk6IkSEhc3MfWkeC1c8NJFpGrL9PgI/5SDRFH5PJStot0Ng\nAh7rumeduRbfH8pr/u3k3yKP767pEALoVW3/M/myGAHliXn+wqC/3Q==\n-----END RSA PRIVATE KEY-----",
			
			// Region
			"region": "uk-london-1",
			"availability_domain": "uFjs:UK-LONDON-1-AD-1",
			
			// Compartment
			"compartment_ocid": "ocid1.compartment.oc1..aaaaaaaacnmuyhg2mpb3z6v6egermq47nai3jk5qaoieg3ztinqhamalealq",
			
			// Compute Instance Configurations
			"instance_display_name": "test-cotud",
			"source_ocid": "ocid1.image.oc1.uk-london-1.aaaaaaaa32voyikkkzfxyo4xbdmadc2dmvorfxxgdhpnk6dw64fa3l4jh7wa",
			"vcn_id": "ocid1.vcn.oc1.uk-london-1.aaaaaaaahdze4mulmvmym4xvrj22cb2zq72tm6qul2oarfnpub6sxwv77owq",
			"cidr" : "10.20.10.0/24",
			"ssh_authorized_keys": "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEAnb7/OPSHIDalG5nQyOJi+5xkKvtGfRr7YOeVvdPxljrbvdzBdJ5+kNXljFI5cy/DZxpfKy7vBoxhSo2Un87eZkcd/TsqSEb92Pxmf4n+W9RXBW4WWxsncv9xQqyUVsK6qIN5ZwmXoXDutz2xfmNFsFBQYwYaPax8wsPJWxVj/9icO4wI+oiWrdDepLPVWYz+8mubhCO/S51IVl1agNT6s1sZqDVKTLTRNeLsMlkHJoVyQhjgWH1PNYxHgJbQR+JxDnp/H8YX7cUrrqjaYU3vCwCGCkFK2Q3RXGcE5QcyOUqzPFyREkNUjuAAE+kYlBMjXLm+b/b6tn0okIAsKlfHCQ== rsa-key-20170117cotudor-ro-mac@COTUDOR-mac",
			"ssh_private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEoQIBAAKCAQEAnb7/OPSHIDalG5nQyOJi+5xkKvtGfRr7YOeVvdPxljrbvdzB\ndJ5+kNXljFI5cy/DZxpfKy7vBoxhSo2Un87eZkcd/TsqSEb92Pxmf4n+W9RXBW4W\nWxsncv9xQqyUVsK6qIN5ZwmXoXDutz2xfmNFsFBQYwYaPax8wsPJWxVj/9icO4wI\n+oiWrdDepLPVWYz+8mubhCO/S51IVl1agNT6s1sZqDVKTLTRNeLsMlkHJoVyQhjg\nWH1PNYxHgJbQR+JxDnp/H8YX7cUrrqjaYU3vCwCGCkFK2Q3RXGcE5QcyOUqzPFyR\nEkNUjuAAE+kYlBMjXLm+b/b6tn0okIAsKlfHCQIBJQKCAQEAkPSzLWsUYspuNQni\nc2g/SBMrnR5Axf0eWQwovEY38dU4oKFX0vKCJDos4c8EXAJgiEG/PHBRREmlgsdK\nTalWvtmRLeNXSVX+BamqS688wxYmc7FE+cXs5jbWx6WBZHuWxF0jc3CZLJFKECra\ndCPfLGV6TSgzyfhypSKdWp4I3UIkrpCAskg7SwOKY7gdzrU+jFUIRSCwps2Scnvg\n3IYqie+fewqukclWi9WhFBFKhPUU3uMzYvZkyljTTn9258VUp5VfBLUm6LIsHiB6\nVSTVfV2twknI+6U7gg8dDbtNltJm8Kg339En8K8+cHsD/lBXJgTSog0PbINbl4Me\nq0H4FQKBgQDTdE3H7+hYdb8oKmXn938X5vs4hoytJkLI/QfSDxbnvflx3mo1RAB5\nrWLbWmjHnrW8vaZ+8MJUmkElb05tYHobTOlaM9keMwRfEpH8gycW/tKOqHRSAlrX\nk4IUmCdeVi8JaULpCl8g7gND/UVjDYBWMSqfIs3ynTXgwgiyVC3iuwKBgQC++jeW\n3zOFOUKQ0Y0YPVe5kNedwWv3jfNB3L+qVJV7QSgu8RJIAKSIaMxK9s5fYIjEc8Zt\nP8z2Va48Hb2szRGXv+cQV75tgUozMAKU+fc7A8B+Dl6PECTPnT3nRXXhVZEMoAg+\nDMvihFpDbw7HxKTIr1yyjcB3UF36uKUFo1MLCwKBgHgDvXF6U3B6LjlkLAAyhmeD\nGPaRjh0VtzPN4dgWZvI7ZBAyIJrFuxSgrbrEnFWfRI23v1zN1hRXjMI4QUT/Z+X7\nOFXKZnjshfDFWcarTYmXjEMhVsbDEPblBKPnp6Q+wMAm/HZtq50Rd3mdlhWf4Q5T\nQbREL7M2ognxljuzPKNHAoGAUpW3LHwx9GvJwhVtckQKQmgl41qPjaUq7Axujth3\n/fKpl8Ixaz6MVqnbzWPO3SLTW98JsrPOQQJ027nVeyg+9YNq1qJ73FOVtUUxjIfE\n2z/kiYmsWYpwyHtZChAy+ah2E0wfPW1RP1vUAXwiETJwxXxDwtWDqTeCle77QLVU\nV80CgYAXGHPu1mY5hc63CxXirt1rh16ylxzAhbe2o5N5dxHqzp+qhdfnifxKup4d\nT9YDlZiJqJmTkuuqcd3FFvFLmvY9B3Ykwm1vQAdDLZ1mDk0AE4GSJYeZxnK2OX5Q\nSqeQ7RwHm9XM7y+gTGpY9oErDCVuC+WGDde+syHIrN6FYJdRgw==\n-----END RSA PRIVATE KEY-----",
			//"assign_public_ip": "true",
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

	keyPair := GetKeyPair("/Users/cotudor/my_ssh_keys/cos_key.openssh", "/Users/cotudor/my_ssh_keys/cos_key.pub")
    //keyPair := GetKeyPair("/Users/cotudor/my_ssh_keys/SantaClara_private_OpenSSH_key.prv", "/Users/cotudor/my_ssh_keys/SantaClara.pub")

	testSSHToPublicHost(t, terraformOptions, keyPair)
	
	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	
	// Verify that we get back a 200 OK with the expected instanceText
	http_helper.HttpGetWithRetry(t, instanceURL, 200, instanceText, maxRetries, timeBetweenRetries)

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