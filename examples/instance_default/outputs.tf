// Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.

output "instance_id" {
  description = "ocid of created instances. "
  value       = ["${module.instance.instance_id}"]
}

output "private_ip" {
  description = "Private IPs of created instances. "
  value       = ["${module.instance.private_ip}"]
}

output "public_ip" {
  description = "Public IPs of created instances. "
  value       = "${module.instance.public_ip}"
}

output "instance_url" {
  description = "http url for Apache app running on the instance "
  value       = "http://${module.instance.public_ip[0]}"
}

output "public_instance_ip" {
  description = "http url for Apache app running on the instance "
  value       = "${module.instance.public_ip[0]}"
}

output "instance_username" {
  description = "Usernames to login to Windows instance. "
  value       = ["${module.instance.instance_username}"]
}

output "instance_password" {
  description = "Passwords to login to Windows instance. "
  sensitive   = true
  value       = ["${module.instance.instance_password}"]
}

output "display_name" {
  description = "Display name. "
  sensitive   = true
  value       = ["${module.instance.display_name}"]
}
