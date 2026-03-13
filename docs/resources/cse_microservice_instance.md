---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice_instance"
description: |-
  Manages a microservice instance resource under a specified microservice engine within HuaweiCloud.
---

# huaweicloud_cse_microservice_instance

Manages a microservice instance resource under a specified microservice engine within HuaweiCloud.

-> Before creating a microservice, make sure that the engine is bound to the EIP and that the rules shown in the
   appendix [table](#microservice_instance_default_engine_access_rules) are enabled in the corresponding security group.

## Example Usage

### Create a microservice instance under a microservice with RBAC authentication of engine disabled

```hcl
variable "microservice_engine_id" {} // Enable the EIP access
variable "microservice_id" {}
variable "region_name" {}
variable "az_name" {}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  fileter_engines = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == var.microservice_engine_id]
}

resource "huaweicloud_cse_microservice_instance" "test" {
  auth_address    = local.fileter_engines[0].service_registry_addresses[0].public
  connect_address = local.fileter_engines[0].service_registry_addresses[0].public

  microservice_id = var.microservice_id
  host_name       = "localhost"
  endpoints       = ["grpc://127.0.1.132:9980", "rest://127.0.0.111:8081"]
  version         = "1.0.0"

  properties = {
    "_TAGS"  = "A, B"
    "attr1"  = "a"
    "nodeIP" = "127.0.0.1"
  }

  health_check {
    mode        = "push"
    interval    = 30
    max_retries = 3
  }

  data_center {
    name              = "dc"
    region            = var.region_name
    availability_zone = var.az_name
  }
}
```

### Create a microservice instance under a microservice with RBAC authentication of engine enabled

```hcl
variable "microservice_engine_id" {} // Enable the EIP access
variable "microservice_engine_admin_password" {}
variable "microservice_id" {}
variable "region_name" {}
variable "az_name" {}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  fileter_engines = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == var.microservice_engine_id]
}

resource "huaweicloud_cse_microservice_instance" "test" {
  auth_address    = local.fileter_engines[0].service_registry_addresses[0].public
  connect_address = local.fileter_engines[0].service_registry_addresses[0].public
  admin_user      = "root"
  admin_pass      = var.microservice_engine_admin_password

  microservice_id = var.microservice_id
  host_name       = "localhost"
  endpoints       = ["grpc://127.0.1.132:9980", "rest://127.0.0.111:8081"]
  version         = "1.0.0"

  properties = {
    "_TAGS"  = "A, B"
    "attr1"  = "a"
    "nodeIP" = "127.0.0.1"
  }

  health_check {
    mode        = "push"
    interval    = 30
    max_retries = 3
  }

  data_center {
    name              = "dc"
    region            = var.region_name
    availability_zone = var.az_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `auth_address` - (Required, String, NonUpdatable) Specifies the address that used to request the access token.  
  Usually is the connection address of service center.

* `connect_address` - (Required, String, NonUpdatable) Specifies the address that used to access engine and manages
  microservice instance.  
  Usually is the connection address of service center.

-> We are only support IPv4 addresses yet (for `auth_address` and `connect_address`).

* `admin_user` - (Optional, String, NonUpdatable) Specifies the user name that used to pass the **RBAC** control.

* `admin_pass` - (Optional, String, NonUpdatable) Specifies the user password that used to pass the **RBAC** control.  
  The password format must meet the following conditions:
  + Must be `8` to `32` characters long.
  + A password must contain at least one digit, one uppercase letter, one lowercase letter, and one special character
    (-~!@#%^*_=+?$&()|<>{}[]).
  + Cannot be the account name or account name spelled backwards.
  + The password can only start with a letter.

-> Both `admin_user` and `admin_pass` are required if **RBAC** is enabled for the microservice engine.

~> Please make sure that all the above parameter values ​​are correct; otherwise, **Terraform** will assume the resource
   does not exist and remove it from the local `.tfstate` file, after which it can only be managed by importing it.

* `microservice_id` - (Required, String, NonUpdatable) Specifies the ID of the microservice to which the microservice
  instance belongs.

* `host_name` - (Required, String, NonUpdatable) Specifies the host name of the microservice instance.

* `endpoints` - (Required, List, NonUpdatable) Specifies the list of access addresses of the microservice instance.

* `version` - (Optional, String, NonUpdatable) Specifies the version of the microservice instance.

* `properties` - (Optional, Map, NonUpdatable) Specifies the extended attributes of the microservice instance,
  in key/value format.

  -> The internal key-value pair cannot be configured or overwritten, such as **engineID** and **engineName**.

* `health_check` - (Optional, List, NonUpdatable) Specifies the health check configuration of the microservice
  instance.  
  The [health_check](#cse_microservice_instance_health_check) structure is documented below.

* `data_center` - (Optional, List, NonUpdatable) Specifies the data center configuration of the microservice instance.  
  The [data_center](#cse_microservice_instance_data_center) structure is documented below.

<a name="cse_microservice_instance_health_check"></a>
The `health_check` block supports:

* `mode` - (Required, String, NonUpdatable) Specifies the heartbeat mode of the health check.  
  The valid values are as follows:
  + **push**
  + **pull**

* `interval` - (Required, Int, NonUpdatable) Specifies the heartbeat interval of the health check, in seconds.

* `max_retries` - (Required, Int, NonUpdatable) Specifies the maximum retry number of the health check.

* `port` - (Optional, Int, NonUpdatable) Specifies the port of the health check.

<a name="cse_microservice_instance_data_center"></a>
The `data_center` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of the data center.

* `region` - (Required, String, NonUpdatable) Specifies the custom region name of the data center.

* `availability_zone` - (Required, String, NonUpdatable) Specifies the custom availability zone name of the data center.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the microservice instance.

* `status` - The status of the microservice instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 3 minutes.

## Import

Microservice instances can be imported using related `auth_address`, `connect_address`, `microservice_id` and their `id`,
separated by the slashes (/), e.g.

```bash
$ terraform import huaweicloud_cse_microservice_instance.test <auth_address>/<connect_address>/<microservice_id>/<id>
```

If you enabled the **RBAC** authorization in the microservice engine, it's necessary to provide the account
name (`admin_user`) and password (`admin_pass`) of the microservice engine.
All fields separated by the slashes (/), e.g.

```bash
$ terraform import huaweicloud_cse_microservice_instance.test <auth_address>/<connect_address>/<microservice_id>/<id>/<admin_user>/<admin_pass>
```

The single quotes (') or backslashes (\\) can help you solve the problem of special characters reporting errors on bash.

```bash
$ terraform import huaweicloud_cse_microservice_instance.test https://124.70.26.32:30100/https://124.70.26.32:30100/f14960ba495e03f59f85aacaaafbdef3fbff3f0d/336e7428dd9411eca913fa163e7364b7/root/Test\!123
```

```bash
$ terraform import huaweicloud_cse_microservice_instance.test 'https://124.70.26.32:30100/https://124.70.26.32:30100/f14960ba495e03f59f85aacaaafbdef3fbff3f0d/336e7428dd9411eca913fa163e7364b7/root/Test!123'
```

## Appendix

<a name="microservice_instance_default_engine_access_rules"></a>
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
