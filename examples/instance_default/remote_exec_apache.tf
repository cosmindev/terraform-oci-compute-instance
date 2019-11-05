# Copyright (c) 2019 Oracle and/or its affiliates,  All rights reserved.

/*
 * @Author: cosmin.tudor@oracle.com 
 * @Date: 2018-11-08 12:44:08 
 * @Last Modified by: cosmin.tudor@oracle.com
 * @Last Modified time: 2018-11-09 12:02:13
 */
resource "null_resource" "configure_cluster_node_apache" {
  count = "${var.instance_count}"


  provisioner "file" {
    connection {
      user        = "opc"
      agent       = false
      private_key = "${chomp(file(var.ssh_private_key_path))}"
      timeout     = "10m"
      host        = "${module.instance.public_ip[count.index]}"
    }

    source      = "apache_install.sh"
    destination = "/tmp/apache_install.sh"
  }

  provisioner "remote-exec" {
    connection {
      user        = "opc"
      agent       = false
      private_key = "${chomp(file(var.ssh_private_key_path))}"
      timeout     = "10m"
      host        = "${module.instance.public_ip[count.index]}"
    }

    inline = [
      "chmod uga+x /tmp/apache_install.sh",
      "sudo su - root -c \"/tmp/apache_install.sh ${module.instance.display_name[count.index]}\"",
    ]
  }
  depends_on = ["module.instance"]
}