---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_service"
description: -|
  Manages a VPC endpoint service resource within HuaweiCloud.
---

# huaweicloud_vpcep_service

Manages a VPC endpoint service resource within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_id" {}
variable "vm_port" {}

resource "huaweicloud_vpcep_service" "demo" {
  name        = "demo-service"
  server_type = "VM"
  vpc_id      = var.vpc_id
  port_id     = var.vm_port
  description = "test description"

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC endpoint service. If omitted, the
  provider-level region will be used. Changing this creates a new VPC endpoint service resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC to which the backend resource of the VPC endpoint
  service belongs. Changing this creates a new VPC endpoint service.

* `server_type` - (Required, String, ForceNew) Specifies the backend resource type. The valid values are as follows:
  + **VM**: Indicates the cloud server, which can be used as a server.
  + **LB**: Indicates the shared load balancer, which is applicable to services with high access traffic and services
    that require high reliability and disaster recovery.

  Changing this creates a new VPC endpoint service.

* `port_id` - (Required, String) Specifies the ID for identifying the backend resource of the VPC endpoint service.
  + If the `server_type` is **VM**, the value is the NIC ID of the ECS where the VPC endpoint service is deployed.
  + If the `server_type` is **LB**, the value is the ID of the port bound to the private IP address of the load
    balancer.

* `port_mapping` - (Required, List) Specifies the port mappings opened to the VPC endpoint service. Structure is
  documented below.

* `name` - (Optional, String) Specifies the name of the VPC endpoint service. The value contains a maximum of 16
  characters, including letters, digits, underscores (_), and hyphens (-).

* `approval` - (Optional, Bool) Specifies whether connection approval is required. The default value is false.

* `permissions` - (Optional, List) Specifies the list of accounts to access the VPC endpoint service.
  The record is in the `iam:domain::domain_id` format, while `*` allows all users to access the VPC endpoint service.

* `organization_permissions` - (Optional, List) Specifies the list of organizations to access the VPC endpoint service.
  The record is in the `organizations:orgPath::org_path` format, while `organizations:orgPath::*` allows all users in
  organizations to access the VPC endpoint service.

* `description` - (Optional, String) Specifies the description of the VPC endpoint service.

* `tcp_proxy` - (Optional, String) Specifies whether to transfer client information (such as source IP address,
  source port number and packet ID) to the server.
  The valid values are as follows:
  + **close**: Neither **TCP TOA** nor **Proxy Protocol** information is carried. Default value.
  + **toa_open**: **TCP TOA** information is carried.
  + **proxy_open**: **Proxy Protocol** information is carried.
  + **open**: Both **TCP TOA** and **Proxy Protocol** information are carried.

  -> 1.**TCP TOA**: The client information is placed into the `tcp option` field and sent to the server.
    This type is available only when the backend resource is an OBS resource.
  <br/>2.**Proxy Protocol**: The client information is placed into the `tcp payload` field and sent to the server.

  -> This parameter is available only when the server can parse the `tcp option` and `tcp payload` fields.

* `ip_version` - (Optional, String, ForceNew) Specifies the IP version of the VPC endpoint service.
  The valid values are as follows:
  + **ipv4** (Default value)
  + **ipv6**

  -> 1.Only professional VPC endpoint service supports this parameter.
    <br>2.Currently, professional VPC endpoint service are available in the **cn-east-4**, **me-east-1**,
    **cn-east-5**, and **af-north-1** regions.

* `snat_network_id` - (Optional, String, ForceNew) Specifies the network ID of any subnet within the VPC used to create
  the VPC endpoint service.

  -> This parameter is valid only when the `ip_version` is set to **ipv6**.

* `enable_policy` - (Optional, Bool, ForceNew) Specifies whether the VPC endpoint policy is enabled. Defaults to **false**.
  Changing this creates a new VPC endpoint service resource.

* `tags` - (Optional, Map) The key/value pairs to associate with the VPC endpoint service.

The `port_mapping` block supports:

* `service_port` - (Required, Int) Specifies the port for accessing the VPC endpoint service. This port is provided by
  the backend service to provide services. The value ranges from `1` to `65,535`.

* `terminal_port` - (Required, Int) Specifies the port for accessing the VPC endpoint. This port is provided by the VPC
  endpoint, allowing you to access the VPC endpoint service. The value ranges from `1` to `65,535`.

* `protocol` - (Optional, String) Specifies the protocol used in port mappings. Only **TCP** is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the VPC endpoint service.

* `status` - The status of the VPC endpoint service. The value can be **available** or **failed**.

* `service_name` - The full name of the VPC endpoint service in the format: *region.name.id* or *region.id*.

* `service_type` - The type of the VPC endpoint service.

* `connections` - An array of VPC endpoints connect to the VPC endpoint service. Structure is documented below.
  + `endpoint_id` - The unique ID of the VPC endpoint.
  + `packet_id` - The packet ID of the VPC endpoint.
  + `domain_id` - The user's domain ID.
  + `status` - The connection status of the VPC endpoint.
  + `description` - The description of the VPC endpoint service connection.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

VPC endpoint services can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpcep_service.test_service <id>
```
