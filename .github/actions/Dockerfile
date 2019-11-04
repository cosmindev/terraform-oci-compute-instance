# Copyright (c) 2019 Oracle and/or its affiliates,  All rights reserved.

FROM golang:1.13

LABEL "name"="test terraform code"
LABEL "version"="1.0"

COPY tf-oci-instance_test.go /tf-oci-instance_test.go
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]