# Copyright (c) 2019 Oracle and/or its affiliates,  All rights reserved.

// Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.



provider "oci" {
  tenancy_ocid = "${var.tenancy_ocid}"
  user_ocid    = "${var.user_ocid}"
  fingerprint  = "${var.fingerprint}"
  private_key  = "${var.private_key}"
  region       = "${var.region}"
}

module "instance" {
  source = "../../"

  instance_count             = "${var.instance_count}"
  availability_domain        = "${var.availability_domain}"
  compartment_ocid           = "${var.compartment_ocid}"
  instance_display_name      = "${var.instance_display_name}"
  source_ocid                = "${var.source_ocid}"
  subnet_ocids               = ["${oci_core_subnet.test_subnet.id}"]
  ssh_authorized_keys        = "${var.ssh_authorized_keys}"
  block_storage_sizes_in_gbs = "${var.block_storage_sizes_in_gbs}"
}

resource "oci_core_subnet" "test_subnet" {
  cidr_block        = var.cidr
  compartment_id    = var.compartment_ocid
  vcn_id            = var.vcn_id
  display_name      = "temp_subnet"
  dns_label         = "tempsubnet"
  security_list_ids = [oci_core_security_list.sec-list.id]
  route_table_id    = "ocid1.routetable.oc1.uk-london-1.aaaaaaaassp2z6qhsn7lcivqgav67pfehtk3bddbqcaglnbzb4ezkyfs5tvq"
  dhcp_options_id   = "ocid1.dhcpoptions.oc1.uk-london-1.aaaaaaaa7ocyua574mbo3l6tid6jiqh4wlzntjebzyhj7ftzf5tugajlvpua"
}

resource "oci_core_security_list" "sec-list" {
  compartment_id = "${var.compartment_ocid}"
  display_name   = "techflow-seclist"
  vcn_id         = "${var.vcn_id}"

  egress_security_rules {
    destination = "0.0.0.0/0"
    protocol    = "${local.tcp_protocol}"
  }

  ingress_security_rules {
    source   = "0.0.0.0/0"
    protocol = "${local.tcp_protocol}"

    tcp_options {
      min = 80
      max = 80
    }
  }

  ingress_security_rules {
    source   = "0.0.0.0/0"
    protocol = "${local.tcp_protocol}"

    tcp_options {
      min = 443
      max = 443
    }
  }

  ingress_security_rules {
    source   = "0.0.0.0/0"
    protocol = "${local.tcp_protocol}"

    tcp_options {
      min = 22
      max = 22
    }
  }
}

locals {
  tcp_protocol  = "6"
  icmp_protocol = "1"
  udp_protocol  = "17"
  vrrp_protocol = "112"
  all_protocols = "all"
  anywhere      = "0.0.0.0/0"
}

