---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice"
description: ""
---

# huaweicloud_cse_microservice

Manages a dedicated microservice resource within HuaweiCloud.

-> 1. Before creating a configuration, make sure the engine has enabled the rules shown in the appendix
   [table](#microservice_default_engine_access_rules).
   <br/> 2. When deleting a microservice, all instances under it will also be deleted together.

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

* `auth_address` - (Required, String, ForceNew) Specifies the address that used to request the access token.  
  Usually is the connection address of service center.  
  Changing this will create a new resource.

* `connect_address` - (Required, String, ForceNew) Specifies the address that used to access engine and manages
  microservice.  
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

* `name` - (Required, String, ForceNew) Specifies the name of the dedicated microservice.
  The name can contain `1` to `128` characters, only letters, digits, underscore (_), hyphens (-) and dots (.) are
  allowed. The name must start and end with a letter or digit. Changing this will create a new microservice.

* `app_name` - (Required, String, ForceNew) Specifies the name of the dedicated microservice application.
  Changing this will create a new microservice.

* `version` - (Required, String, ForceNew) Specifies the version of the dedicated microservice.
  Changing this will create a new microservice.

* `environment` - (Optional, String, ForceNew) Specifies the environment (stage) type.
  The valid values are **development**, **testing**, **acceptance** and **production**.
  If omitted, the microservice will be deployed in an empty environment.
  Changing this will create a new microservice.

* `level` - (Optional, String, ForceNew) Specifies the microservice level.
  The valid values are **FRONT**, **MIDDLE**, and **BACK**. Changing this will create a new microservice.

* `description` - (Optional, String, ForceNew) Specifies the description of the dedicated microservice.
  The description can contain a maximum of `256` characters.
  Changing this will create a new microservice.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The microservice ID.

* `status` - The microservice status. The values supports **UP** and **DOWN**.

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
| Ingress   | 1        | Allow  | ICMP     | All           | Ipv4      | 0.0.0.0/0             |
| Ingress   | 1        | Allow  | TCP      | 30100-30130   | Ipv4      | 0.0.0.0/0             |
| Ingress   | 1        | Allow  | All      | All           | Ipv4      | cse-engine-default-sg |
| Egress    | 100      | Allow  | All      | All           | Ipv6      | ::/0                  |
| Egress    | 100      | Allow  | All      | All           | Ipv4      | 0.0.0.0/0             |
