---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice_instance"
description: ""
---

# huaweicloud_cse_microservice_instance

Manages a dedicated microservice instance resource within HuaweiCloud.

-> Before creating a configuration, make sure the engine has enabled the rules shown in the appendix
   [table](#microservice_instance_default_engine_access_rules).

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

* `auth_address` - (Required, String, ForceNew) Specifies the address that used to request the access token.  
  Usually is the connection address of service center.  
  Changing this will create a new resource.

* `connect_address` - (Required, String, ForceNew) Specifies the address that used to access engine and manages
  microservice instance.  
  Usually is the connection address of service center.  
  Changing this will create a new resource.

-> We are only support IPv4 addresses yet (for `auth_address` and `connect_address`).

* `admin_user` - (Optional, String, ForceNew) Specifies the account name for **RBAC** login.
  Changing this will create a new resource.

* `admin_pass` - (Optional, String, ForceNew) Specifies the account password for **RBAC** login.
  The password format must meet the following conditions:
  + Must be `8` to `32` characters long.
  + A password must contain at least one digit, one uppercase letter, one lowercase letter, and one special character
    (-~!@#%^*_=+?$&()|<>{}[]).
  + Cannot be the account name or account name spelled backwards.
  + The password can only start with a letter.

  Changing this will create a new resource.

-> Both `admin_user` and `admin_pass` are required if **RBAC** is enabled for the microservice engine.

* `microservice_id` - (Required, String, ForceNew) Specifies the ID of the dedicated microservice to which the instance
  belongs. Changing this will create a new microservice instance.

* `host_name` - (Required, String, ForceNew) Specifies the host name, such as `localhost`.
  Changing this will create a new microservice instance.

* `endpoints` - (Required, List, ForceNew) Specifies the access addresses information.
  Changing this will create a new microservice instance.

* `version` - (Optional, String, ForceNew) Specifies the version of the dedicated microservice instance.
  Changing this will create a new microservice instance.

* `properties` - (Optional, Map, ForceNew) Specifies the extended attributes.
  Changing this will create a new microservice instance.

  -> The internal key-value pair cannot be configured or overwritten, such as **engineID** and **engineName**.

* `health_check` - (Optional, List, ForceNew) Specifies the health check configuration.
  The [object](#microservice_instance_health_check) structure is documented below.
  Changing this will create a new microservice instance.

* `data_center` - (Optional, List, ForceNew) Specifies the data center configuration.
  The [object](#microservice_instance_data_center) structure is documented below.
  Changing this will create a new microservice instance.

* `admin_user` - (Optional, String, ForceNew) Specifies the account name. The initial account name is **root**.
  Required if the `auth_type` of engine is **RBAC**. Changing this will create a new microservice instance.

* `admin_pass` - (Optional, String, ForceNew) Specifies the account password.
  Required if the `auth_type` of engine is **RBAC**. Changing this will create a new microservice instance.
  The password format must meet the following conditions:
  + Must be `8` to `32` characters long.
  + A password must contain at least one digit, one uppercase letter, one lowercase letter, and one special character
    (-~!@#%^*_=+?$&()|<>{}[]).
  + Cannot be the account name or account name spelled backwards.
  + The password can only start with a letter.

<a name="microservice_instance_health_check"></a>
The `health_check` block supports:

* `mode` - (Required, String, ForceNew) Specifies the heartbeat mode. The valid values are **push** and **pull**.
  Changing this will create a new microservice instance.

* `interval` - (Required, Int, ForceNew) Specifies the heartbeat interval. The unit is **s** (second).
  Changing this will create a new microservice instance.

* `max_retries` - (Required, Int, ForceNew) Specifies the maximum retries.
  Changing this will create a new microservice instance.

* `port` - (Optional, Int, ForceNew) Specifies the port number.
  Changing this will create a new microservice instance.

<a name="microservice_instance_data_center"></a>
The `data_center` block supports:

* `name` - (Required, String, ForceNew) Specifies the data center name.
  Changing this will create a new microservice instance.

* `region` - (Required, String, ForceNew) Specifies the custom region name of the data center.
  Changing this will create a new microservice instance.

* `availability_zone` - (Required, String, ForceNew) Specifies the custom availability zone name of the data center.
  Changing this will create a new microservice instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The microservice instance ID.

* `status` - The microservice instance status. The values supports **UP**, **DOWN**, **STARTING** and **OUTOFSERVICE**.

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
| Ingress   | 1        | Allow  | ICMP     | All           | Ipv4      | 0.0.0.0/0             |
| Ingress   | 1        | Allow  | TCP      | 30100-30130   | Ipv4      | 0.0.0.0/0             |
| Ingress   | 1        | Allow  | All      | All           | Ipv4      | cse-engine-default-sg |
| Egress    | 100      | Allow  | All      | All           | Ipv6      | ::/0                  |
| Egress    | 100      | Allow  | All      | All           | Ipv4      | 0.0.0.0/0             |
