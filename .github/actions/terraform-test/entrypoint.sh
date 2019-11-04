#!/bin/bash

apt-get update
apt install -y build-essential unzip go-dep

#Installing terraform and upgrading code to Terraform 12
wget -q https://releases.hashicorp.com/terraform/0.12.10/terraform_0.12.10_linux_amd64.zip
unzip terraform_0.12.10_linux_amd64.zip -d /usr/bin
cd ${GITHUB_WORKSPACE}/examples/instance_default
echo ${GITHUB_WORKSPACE}
terraform init && terraform 0.12upgrade -yes


#Set up environment for running terratest in go
mkdir -p $HOME/go/src/terratest/test
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOROOT/bin:$GOPATH/bin:/usr/bin:$PATH
mv /tf-oci-instance_test.go $HOME/go/src/terratest/test
cd $HOME/go/src/terratest/test

#Install terratest and dependencies
cat  << EOF > Gopkg.toml
 [[constraint]]
  name = "github.com/gruntwork-io/terratest"
  version = "0.22.2"
 [prune]
  go-tests = true
  unused-packages = true
EOF

dep ensure

#Set up environment to run the terraform code
#echo "${TF_VAR_private_key}" > ${GITHUB_WORKSPACE}/oci.pem
#export TF_VAR_private_key_path=${GITHUB_WORKSPACE}/oci.pem
export TF_ACTION_WORKING_DIR=${GITHUB_WORKSPACE}/examples/instance_default
TERRAFORM_VERSION="$(terraform --version)"
echo "------->>>Terraform version = ${TERRAFORM_VERSION}"
echo "------->>> TF_VAR_tenancy_id= ${TF_VAR_tenancy_id}"
echo "------->>> TF_VAR_user_id= ${TF_VAR_user_id}"
echo "------->>> TF_VAR_fingerprint= ${TF_VAR_fingerprint}"
echo "------->>> TF_VAR_private_key= ${TF_VAR_private_key}"
echo "------->>> TF_VAR_ssh_authorized_keys= ${TF_VAR_ssh_authorized_keys}"
echo "------->>> TF_VAR_ssh_private_key= ${TF_VAR_ssh_private_key}"
ENV="$(env)"
echo "------->>> ENV= ${ENV}"

export TF_VAR_tenancy_id='ocid1.tenancy.oc1..aaaaaaaaa3qmjxr43tjexx75r6gwk6vjw22ermohbw2vbxyhczksgjir7xdq'

export TF_VAR_user_id='ocid1.user.oc1..aaaaaaaaoj3k4zp24w2aryx7z46gzhw2yxf2lrvcj3om6cikwb6k4c2p74sq'

export TF_VAR_fingerprint='25:93:69:40:2f:5b:d2:25:0e:eb:f3:41:ea:cb:18:02'

export TF_VAR_private_key='-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAx7CWDt6MCoM4JMQbvrTaRGlHcBJjtJ94iGxJqylkVjeNit11\nSgQcV0vq+evL0GsW52TptNPlFvboZC8thGEtGQxe4NT70OHkjvnbJBiFA+5055iN\nMSkLZ5Sex0+BRwLIfOay15WqcFlgFK75Oy1uO0EY/73HVXQFttZGqZc4Sc2b85Pg\nbelSv8cn/WGGp1Wqcp+m47C7zTeYoIYyecWLr/mF6dTFhiqQwTU6301bmlsAPpjF\nZBo0411LH68lRo5Rkzkf6ulxZB1EXkNP16kyOgRj8pp+nxpMQt5kdaJPYYsUFXhM\nnOs0lvQ0TpDAJZ8xfOUAlmBFHoSUWx1x07XsxwIDAQABAoIBAQCDuhlC9curo6QF\nnNgwlVzmhAx8AaWEIS6Bz+1I26U0urDEShJ1IQERFSOMed+ZPQt+2TLR9nXJEFd8\nvyJnkTMOsvpjlhwHTvxW1LuatBIK5eJ4ZEm/hDPhwEh9chT1owBeZFDVpgUo18xp\n2ILQZsewjuDO04IK/N5IHlG+zeJzy+NpEg9vtPN1A3aDomSMqy+DrvuTFpH6CNtM\nNrDTqDccQbnF0cZSXD1L6mw4kWol8jAsSyGDf7uOOSyNT57ky80mRDBu9CuA1N9j\npnJCZaWBrQXpKco5lBLSYXJ1JddzvSSim8QV+rYd0+004ad7qojigBQEmGLOAT81\nhZkSQCwBAoGBAPKlBVV5+AyPuf9QG7QaAceMwHWIhjxZemJ6o7QARKFc2Dc1CYQO\nAutZV4JXGhrQLV+wuHY2JG5AXyVOxolWsW14BYxx0+Oq6U0p6OOetYi3LGiLCT4N\nwuAp4lMLEPfPNht7SFfq7ZannaaREzhr4LQXKlGMEszZsZ+/KATyDwJxAoGBANKu\nT609cDaXL/TnfOpWqV5bXTUiNjdZDT05S62Ynvr5KkXBa97pE0cuZ7cwGgjxC5Pe\nDteGcbUe/Eikx8FVvqfuYRvKuhAoPdFZ1miR3sWDHOlFjkuJxaGvUbBDlBwcKHx9\nAOxAruhCRmG5/dYAIJ5KlEYJF5xsZz9SmQ7A6Q63AoGAEuaZ2MOsd8YGVgX2cnwI\nIXQsVbtxwWey6dLlx5KxxeQGj55ZBGlW/uAxudxxEx+eOTL523NyOQhYoP5W5sHT\nBlTwEbWYLMbWb4VRN9HYEDM8iVQzPxsxT+bTU2asRrFkZJWg2ABby315AU2RsrZs\nhXq8eCeyGzTl6iyowGHem3ECgYBdfVCKFctnzitPyDGcY5yA7JYt7+KTKQdA9d3p\nSOKziEID9lMB9ffCDIultMi40w5KLa30YgqvTvKw4b5qwrv6FUQuawWqCdF0xyLo\nAGMUzpvTwDPmvVpf50aeqz5cQvqMU4RHUmTLWC2XTEuh3SicVYf6lCpQFaKzbNnS\nvDQvfQKBgQCAqS2DriooFdQpcblAyN2twGBQbGUYw6FUA0N9B8jMzxoq+B3P6jku\nDmhhVv0pf0pgwlmWqNk6IkSEhc3MfWkeC1c8NJFpGrL9PgI/5SDRFH5PJStot0Ng\nAh7rumeduRbfH8pr/u3k3yKP767pEALoVW3/M/myGAHliXn+wqC/3Q==\n-----END RSA PRIVATE KEY-----\n'

export TF_VAR_ssh_authorized_keys='ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEAnb7/OPSHIDalG5nQyOJi+5xkKvtGfRr7YOeVvdPxljrbvdzBdJ5+kNXljFI5cy/DZxpfKy7vBoxhSo2Un87eZkcd/TsqSEb92Pxmf4n+W9RXBW4WWxsncv9xQqyUVsK6qIN5ZwmXoXDutz2xfmNFsFBQYwYaPax8wsPJWxVj/9icO4wI+oiWrdDepLPVWYz+8mubhCO/S51IVl1agNT6s1sZqDVKTLTRNeLsMlkHJoVyQhjgWH1PNYxHgJbQR+JxDnp/H8YX7cUrrqjaYU3vCwCGCkFK2Q3RXGcE5QcyOUqzPFyREkNUjuAAE+kYlBMjXLm+b/b6tn0okIAsKlfHCQ== rsa-key-20170117cotudor-ro-mac@COTUDOR-mac'

export TF_VAR_ssh_private_key=$'-----BEGIN RSA PRIVATE KEY-----\nMIIEoQIBAAKCAQEAnb7/OPSHIDalG5nQyOJi+5xkKvtGfRr7YOeVvdPxljrbvdzB\ndJ5+kNXljFI5cy/DZxpfKy7vBoxhSo2Un87eZkcd/TsqSEb92Pxmf4n+W9RXBW4W\nWxsncv9xQqyUVsK6qIN5ZwmXoXDutz2xfmNFsFBQYwYaPax8wsPJWxVj/9icO4wI\n+oiWrdDepLPVWYz+8mubhCO/S51IVl1agNT6s1sZqDVKTLTRNeLsMlkHJoVyQhjg\nWH1PNYxHgJbQR+JxDnp/H8YX7cUrrqjaYU3vCwCGCkFK2Q3RXGcE5QcyOUqzPFyR\nEkNUjuAAE+kYlBMjXLm+b/b6tn0okIAsKlfHCQIBJQKCAQEAkPSzLWsUYspuNQni\nc2g/SBMrnR5Axf0eWQwovEY38dU4oKFX0vKCJDos4c8EXAJgiEG/PHBRREmlgsdK\nTalWvtmRLeNXSVX+BamqS688wxYmc7FE+cXs5jbWx6WBZHuWxF0jc3CZLJFKECra\ndCPfLGV6TSgzyfhypSKdWp4I3UIkrpCAskg7SwOKY7gdzrU+jFUIRSCwps2Scnvg\n3IYqie+fewqukclWi9WhFBFKhPUU3uMzYvZkyljTTn9258VUp5VfBLUm6LIsHiB6\nVSTVfV2twknI+6U7gg8dDbtNltJm8Kg339En8K8+cHsD/lBXJgTSog0PbINbl4Me\nq0H4FQKBgQDTdE3H7+hYdb8oKmXn938X5vs4hoytJkLI/QfSDxbnvflx3mo1RAB5\nrWLbWmjHnrW8vaZ+8MJUmkElb05tYHobTOlaM9keMwRfEpH8gycW/tKOqHRSAlrX\nk4IUmCdeVi8JaULpCl8g7gND/UVjDYBWMSqfIs3ynTXgwgiyVC3iuwKBgQC++jeW\n3zOFOUKQ0Y0YPVe5kNedwWv3jfNB3L+qVJV7QSgu8RJIAKSIaMxK9s5fYIjEc8Zt\nP8z2Va48Hb2szRGXv+cQV75tgUozMAKU+fc7A8B+Dl6PECTPnT3nRXXhVZEMoAg+\nDMvihFpDbw7HxKTIr1yyjcB3UF36uKUFo1MLCwKBgHgDvXF6U3B6LjlkLAAyhmeD\nGPaRjh0VtzPN4dgWZvI7ZBAyIJrFuxSgrbrEnFWfRI23v1zN1hRXjMI4QUT/Z+X7\nOFXKZnjshfDFWcarTYmXjEMhVsbDEPblBKPnp6Q+wMAm/HZtq50Rd3mdlhWf4Q5T\nQbREL7M2ognxljuzPKNHAoGAUpW3LHwx9GvJwhVtckQKQmgl41qPjaUq7Axujth3\n/fKpl8Ixaz6MVqnbzWPO3SLTW98JsrPOQQJ027nVeyg+9YNq1qJ73FOVtUUxjIfE\n2z/kiYmsWYpwyHtZChAy+ah2E0wfPW1RP1vUAXwiETJwxXxDwtWDqTeCle77QLVU\nV80CgYAXGHPu1mY5hc63CxXirt1rh16ylxzAhbe2o5N5dxHqzp+qhdfnifxKup4d\nT9YDlZiJqJmTkuuqcd3FFvFLmvY9B3Ykwm1vQAdDLZ1mDk0AE4GSJYeZxnK2OX5Q\nSqeQ7RwHm9XM7y+gTGpY9oErDCVuC+WGDde+syHIrN6FYJdRgw==\n-----END RSA PRIVATE KEY-----\n'

echo "------->>> TF_VAR_tenancy_id= ${TF_VAR_tenancy_id}"
echo "------->>> TF_VAR_user_id= ${TF_VAR_user_id}"
echo "------->>> TF_VAR_fingerprint= ${TF_VAR_fingerprint}"
echo "------->>> TF_VAR_private_key= ${TF_VAR_private_key}"
echo "------->>> TF_VAR_ssh_authorized_keys= ${TF_VAR_ssh_authorized_keys}"
echo "------->>> TF_VAR_ssh_private_key= ${TF_VAR_ssh_private_key}"

go test -v $HOME/go/src/terratest/test/tf-oci-instance_test.go -timeout 20m
