---
layout: "huaweicloud"
page_title: "HuaweiCloud: resource_huaweicloud_vpc_peering_connection_accepter_v2"
sidebar_current: "docs-huaweicloud-resource-vpc-peering-accepter-v2"
description: |-
  Manage the accepter's side of a VPC Peering Connection.
---

# huaweicloud_vpc_peering_connection_accepter_v2

Provides a resource to manage the accepter's side of a VPC Peering Connection.

When a cross-tenant (requester's tenant differs from the accepter's tenant) VPC Peering Connection is created, a VPC Peering Connection resource is automatically created in the
accepter's account.
The requester can use the `huaweicloud_vpc_peering_connection_v2` resource to manage its side of the connection
and the accepter can use the `huaweicloud_vpc_peering_connection_accepter_v2` resource to "adopt" its side of the
connection into management.

## Example Usage

 ```hcl
 provider "huaweicloud"  {
    alias = "main"
    user_name   = "${var.username}"
    domain_name = "${var.domain_name}"
    password    = "${var.password}"
    auth_url    = "${var.auth_url}"
    region      = "${var.region}"
    tenant_id   = "${var.tenant_id}"
}

provider "huaweicloud"  {
    alias = "peer"
    user_name = "${var.peer_username}"
    domain_name = "${var.peer_domain_name}"
    password    = "${var.peer_password}"
    auth_url    = "${var.peer_auth_url}"
    region      = "${var.peer_region}"
    tenant_id   = "${var.peer_tenant_id}"
}

resource "huaweicloud_vpc_v1" "vpc_main" {
    provider = "huaweicloud.main"
    name = "${var.vpc_name}"
    cidr = "${var.vpc_cidr}"
}

resource "huaweicloud_vpc_v1" "vpc_peer" {
    provider = "huaweicloud.peer"
    name = "${var.peer_vpc_name}"
    cidr = "${var.peer_vpc_cidr}"
}

# Requester's side of the connection.
resource "huaweicloud_vpc_peering_connection_v2" "peering" {
    provider = "huaweicloud.main"
    name = "${var.peer_name}"
    vpc_id = "${huaweicloud_vpc_v1.vpc_main.id}"
    peer_vpc_id = "${huaweicloud_vpc_v1.vpc_peer.id}"
    peer_tenant_id =  "${var.tenant_id}"
}

# Accepter's side of the connection.
resource "huaweicloud_vpc_peering_connection_accepter_v2" "peer" {
    provider = "huaweicloud.peer"
    vpc_peering_connection_id = "${huaweicloud_vpc_peering_connection_v2.peering.id}"
    accept = true
  
}
 ```

## Argument Reference

The following arguments are supported:

* `vpc_peering_connection_id` (Required) - The VPC Peering Connection ID to manage. Changing this creates a new VPC peering connection accepter.

* `accept` (Optional)- Whether or not to accept the peering request. Defaults to `false`.


## Removing huaweicloud_vpc_peering_connection_accepter_v2 from your configuration
 
huaweicloud allows a cross-tenant VPC Peering Connection to be deleted from either the requester's or accepter's side. However, Terraform only allows the VPC Peering Connection to be deleted from the requester's side by removing the corresponding `huaweicloud_vpc_peering_connection_v2` resource from your configuration. Removing a `huaweicloud_vpc_peering_connection_accepter_v2` resource from your configuration will remove it from your state file and management, but will not destroy the VPC Peering Connection.

## Attributes Reference

All of the argument attributes except accept are also exported as result attributes.

* `name` - 	The VPC peering connection name.

* `id` - The VPC peering connection ID.

* `status` - The VPC peering connection status.

* `vpc_id` - The ID of requester VPC involved in a VPC peering connection.

* `peer_vpc_id` - The VPC ID of the accepter tenant.

* `peer_tenant_id` - The Tenant Id of the accepter tenant.


