---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice"
description: |-
  Manages a microservice resource under a specified microservice engine within HuaweiCloud.
---

# huaweicloud_cse_microservice

Manages a microservice resource under a specified microservice engine within HuaweiCloud.

-> 1. Before creating a microservice, make sure that the engine is bound to the EIP and that the rules shown in the
   appendix [table](#microservice_default_engine_access_rules) are enabled in the corresponding security group.
   <br/>2. When deleting a microservice, all instances under it will also be deleted together.

## Example Usage

### Create a microservice in an engine with RBAC authentication disabled

```hcl
variable "microservice_engine_id" {} // Enable the EIP access
variable "service_name" {}
variable "app_name" {}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  fileter_engines = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == var.microservice_engine_id]
}

resource "huaweicloud_cse_microservice" "test" {
  auth_address    = local.fileter_engines[0].service_registry_addresses[0].public
  connect_address = local.fileter_engines[0].service_registry_addresses[0].public

  name        = var.service_name
  version     = "1.0.0"
  environment = "development"
  app_name    = var.app_name
}
```

### Create a microservice in an engine with RBAC authentication enabled

```hcl
variable "microservice_engine_id" {} // Enable the EIP access
variable "microservice_engine_admin_password" {}
variable "service_name" {}
variable "app_name" {}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  fileter_engines = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == var.microservice_engine_id]
}

resource "huaweicloud_cse_microservice" "test" {
  auth_address    = local.fileter_engines[0].service_registry_addresses[0].public
  connect_address = local.fileter_engines[0].service_registry_addresses[0].public
  admin_user      = "root"
  admin_pass      = var.microservice_engine_admin_password

  name        = var.service_name
  version     = "1.0.0"
  environment = "development"
  app_name    = var.app_name
}
```

## Argument Reference

The following arguments are supported:

* `auth_address` - (Required, String, NonUpdatable) Specifies the address that used to request the access token.  
  Usually is the connection address of service center.

* `connect_address` - (Required, String, NonUpdatable) Specifies the address that used to access engine and manages
  microservice.  
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

* `name` - (Required, String, NonUpdatable) Specifies the name of the microservice.  
  The valid value is limited from `1` to `128` characters, only letters, digits, underscore (_), hyphens (-) and
  dots (.) are allowed. The name must start and end with a letter or digit.

* `app_name` - (Required, String, NonUpdatable) Specifies the name of the microservice application.

* `version` - (Required, String, NonUpdatable) Specifies the version of the microservice.

* `environment` - (Optional, String, NonUpdatable) Specifies the environment (stage) type of the microservice.  
  The valid values are as follows:
  + **development**
  + **testing**
  + **acceptance**
  + **production**

  If omitted, the microservice will be deployed in an empty environment.

* `level` - (Optional, String, NonUpdatable) Specifies the level of the microservice.  
  The valid values are as follows:
  + **FRONT**
  + **MIDDLE**
  + **BACK**

* `description` - (Optional, String, NonUpdatable) Specifies the description of the microservice.  
  The description can contain a maximum of `256` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the microservice.

* `status` - The status of the microservice.

## Import

Microservices can be imported using related `auth_address`, `connect_address` and their `id`, separated by the
slashes (/), e.g.

```bash
$ terraform import huaweicloud_cse_microservice.test <auth_address>/<connect_address>/<id>
```

If you enabled the **RBAC** authorization, you also need to provide the account name (`admin_user`) and password
(`admin_pass`) of the microservice engine. All fields separated by the slashes (/), e.g.

```bash
$ terraform import huaweicloud_cse_microservice.test <auth_address>/<connect_address>/<id>/<admin_user>/<admin_pass>
```

The single quotes (') or backslashes (\\) can help you solve the problem of special characters reporting errors on bash.

```bash
$ terraform import huaweicloud_cse_microservice.test https://124.70.26.32:30100/https://124.70.26.32:30100/f14960ba495e03f59f85aacaaafbdef3fbff3f0d/root/Test\!123
```

```bash
$ terraform import huaweicloud_cse_microservice.test 'https://124.70.26.32:30100/https://124.70.26.32:30100/f14960ba495e03f59f85aacaaafbdef3fbff3f0d/root/Test!123'
```

## Appendix

<a name="microservice_default_engine_access_rules"></a>
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
