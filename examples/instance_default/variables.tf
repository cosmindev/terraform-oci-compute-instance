# Copyright (c) 2019 Oracle and/or its affiliates,  All rights reserved.

variable "tenancy_ocid" {}
variable "user_ocid" {}
variable "fingerprint" {}
variable "private_key_path" {}
variable "region" {}
variable "availability_domain" {}
variable "instance_count" {}
variable "compartment_ocid" {}
variable "instance_display_name" {}

variable "vcn_id" {}
variable "cidr" {}
variable "source_ocid" {}
variable "ssh_authorized_keys" {}

variable "ssh_private_key_path" {}

variable "block_storage_sizes_in_gbs" {
  type = "list"
}