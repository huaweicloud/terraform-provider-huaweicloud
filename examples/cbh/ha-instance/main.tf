# This example demonstrates how to create a CBH HA instance in Huawei Cloud using Terraform.
# It includes the creation of a VPC, subnet, security group, and the CBH instance itself.
# The example also includes the necessary variables and outputs for the created resources.

#List the availability zones
data "huaweicloud_availability_zones" "default" {}

#create a VPC
resource "huaweicloud_vpc" "default" {
  name = "cbh_vpc"
  cidr = "192.168.0.0/16"
}

#Create a subnet
resource "huaweicloud_vpc_subnet" "default" {
  name       = "cbh_subnet"
  cidr       = "192.168.0.0/20"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.default.id
}

#create a security group
resource "huaweicloud_networking_secgroup" "default" {
  name = "cbh_security_group"
}

resource "huaweicloud_cbh_instance" "cbh_HA_demo" {
    name                     = var.name
    flavor_id                = var.flavor_id
    vpc_id                   = huaweicloud_vpc.default.id
    subnet_id                = huaweicloud_vpc_subnet.default.id
    security_group_id        = huaweicloud_networking_secgroup.default.id
    master_availability_zone = data.huaweicloud_availability_zones.default.names[0]
    slave_availability_zone  = data.huaweicloud_availability_zones.default.names[0]
    password                 = var.password
    charging_mode            = var.charging_mode
    period_unit              = var.period_unit
    period                   = var.period          
}

#output the CBH instance ID
output "cbh_instance_id" {
    value       = huaweicloud_cbh_instance.cbh_HA_demo.id
    description = "The ID of the CBH HA instance."
}
