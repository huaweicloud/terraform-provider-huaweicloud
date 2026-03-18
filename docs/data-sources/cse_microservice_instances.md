---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice_instances"
description: |-
  Use this data source to get the list of the microservice instances under specified microservice within HuaweiCloud.
---

# huaweicloud_cse_microservice_instances

Use this data source to get the list of the microservice instances under specified microservice within HuaweiCloud.

-> Before creating a microservice, make sure that the engine is bound to the EIP and that the rules shown in the
   appendix [table](#microservice_instances_default_engine_access_rules) are enabled in the corresponding security group.

## Example Usage

```hcl
variable "microservice_engine_id" {} // Enable the EIP access
variable "admin_user" {}
variable "admin_password" {}
variable "microservice_id" {}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  filter_engines = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == var.microservice_engine_id]
}

data "huaweicloud_cse_microservice_instances" "test" {
  auth_address    = local.filter_engines[0].service_registry_addresses[0].public
  connect_address = local.filter_engines[0].service_registry_addresses[0].public
  admin_user      = var.admin_user
  admin_pass      = var.admin_password

  microservice_id = var.microservice_id
}
```

## Argument Reference

The following arguments are supported:

* `auth_address` - (Required, String) Specifies the address that used to request the access token.

* `connect_address` - (Required, String) Specifies the address that used to access engine and query microservice
  instances.

-> We are only support IPv4 addresses yet (for `auth_address` and `connect_address`).

* `admin_user` - (Optional, String) Specifies the user name that used to pass the **RBAC** control.

* `admin_pass` - (Optional, String) Specifies the user password that used to pass the **RBAC** control.  
  The password format must meet the following conditions:
  + Must be `8` to `32` characters long.
  + A password must contain at least one digit, one uppercase letter, one lowercase letter, and one special character
    (-~!@#%^*_=+?$&()|<>{}[]).
  + Cannot be the account name or account name spelled backwards.
  + The password can only start with a letter.

-> Both `admin_user` and `admin_pass` are required if **RBAC** is enabled for the microservice engine.

~> Please make sure that all the above parameter values ​​are correct; otherwise, **Terraform** will report a connection
   error.

* `microservice_id` - (Required, String) Specifies the ID of the microservice to which the microservice instances
  belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of the microservice instances.
  The [instances](#microservice_instances) structure is documented below.

<a name="microservice_instances"></a>
The `instances` block supports:

* `id` - The ID of the microservice instance.

* `host_name` - The host name of the microservice instance.

* `endpoints` - The list of the access addresses of the microservice instance.

* `version` - The version of the microservice instance.

* `properties` - The extended attributes of the microservice instance, in key/value format.

* `health_check` - The health check configuration of the microservice instance.
  The [health_check](#microservice_instances_health_check) structure is documented below.

* `data_center` - The data center configuration of the microservice instance.
  The [data_center](#microservice_instances_data_center) structure is documented below.

* `status` - The current status of the microservice instance.
  + **UP**
  + **DOWN**
  + **STARTING**
  + **OUTOFSERVICE**

* `created_at` - The creation time of the microservice instance, in RFC3339 format.

* `updated_at` - The latest update time of the microservice instance, in RFC3339 format.

<a name="microservice_instances_health_check"></a>
The `health_check` block supports:

* `mode` - The heartbeat mode of the health check.

* `interval` - The heartbeat interval of the health check, in seconds.

* `max_retries` - The maximum retry number of the health check.

* `port` - The port of the health check.

<a name="microservice_instances_data_center"></a>
The `data_center` block supports:

* `name` - The name of the data center.

* `region` - The custom region name of the data center.

* `availability_zone` - The custom availability zone of the data center.

## Appendix

<a name="microservice_instances_default_engine_access_rules"></a>
Security group rules required to access the engine:
(Remote is not the minimum range and can be adjusted according to business needs)

| Direction | Priority | Action | Protocol | Ports         | Ethertype | Remote                |
| --------- | -------- | ------ | -------- | ------------- | --------- | --------------------- |
| Ingress   | 1        | Allow  | ICMP     | All           | Ipv6      | ::/0                  |
| Ingress   | 1        | Allow  | TCP      | 30100-30130   | Ipv6      | ::/0                  |
| Ingress   | 1        | Allow  | All      | All           | Ipv6      | cse-engine-default-sg |
| Ingress   | 1        | Allow  | ICMP     | All           | Ipv4      | A CIDR containing the public IP address of the machine from which you executed this terraform script, such as **0.0.0.0/0** |
| Ingress   | 1        | Allow  | TCP      | 30100-30130   | Ipv4      | A CIDR containing the public IP address of the machine from which you executed this terraform script, such as **0.0.0.0/0** |
| Ingress   | 1        | Allow  | All      | All           | Ipv4      | cse-engine-default-sg |
| Egress    | 100      | Allow  | All      | All           | Ipv6      | ::/0                  |
| Egress    | 100      | Allow  | All      | All           | Ipv4      | 0.0.0.0/0             |
