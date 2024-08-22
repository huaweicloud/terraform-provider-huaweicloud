---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice"
description: ""
---

# huaweicloud_cse_microservice

Manages a dedicated microservice resource within HuaweiCloud.

-> When deleting a microservice, all instances under it will also be deleted together.

## Example Usage

### Create a microservice in an engine with RBAC authentication disabled

```hcl
variable "engine_conn_addr" {}
variable "service_name" {}
variable "app_name" {}

resource "huaweicloud_cse_microservice" "test" {
  connect_address = var.engine_conn_addr
  name            = var.service_name
  version         = "1.0.0"
  environment     = "development"
  app_name        = var.app_name
}
```

### Create a microservice in an engine with RBAC authentication enabled

```hcl
variable "engine_conn_addr" {}
variable "service_name" {}
variable "app_name" {}
variable "admin_pass" {}

resource "huaweicloud_cse_microservice" "test" {
  connect_address = var.engine_conn_addr
  name            = var.service_name
  version         = "1.0.0"
  environment     = "development"
  app_name        = var.app_name

  admin_user = "root"
  admin_pass = var.admin_pass
}
```

## Argument Reference

The following arguments are supported:

* `connect_address` - (Required, String, ForceNew) Specifies the connection address of service registry center for the
  specified dedicated CSE engine. Changing this will create a new microservice.

-> We are only support IPv4 addresses yet.

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

* `admin_user` - (Optional, String, ForceNew) Specifies the account name. The initial account name is **root**.
  Required if the `auth_type` of engine is **RBAC**. Changing this will create a new microservice.

* `admin_pass` - (Optional, String, ForceNew) Specifies the account password.
  Required if the `auth_type` of engine is **RBAC**. Changing this will create a new microservice.
  The password format must meet the following conditions:
  + Must be `8` to `32` characters long.
  + A password must contain at least one digit, one uppercase letter, one lowercase letter, and one special character
    (-~!@#%^*_=+?$&()|<>{}[]).
  + Cannot be the account name or account name spelled backwards.
  + The password can only start with a letter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The microservice ID.

* `status` - The microservice status. The values supports **UP** and **DOWN**.

## Import

Microservices can be imported using related `connect_address` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_cse_microservice.test https://124.70.26.32:30100/f14960ba495e03f59f85aacaaafbdef3fbff3f0d
```

If you enabled the **RBAC** authorization, you also need to provide the account name and password, e.g.

```bash
$ terraform import huaweicloud_cse_microservice.test 'https://124.70.26.32:30100/f14960ba495e03f59f85aacaaafbdef3fbff3f0d/root/Test!123'
```

The single quotes can help you solve the problem of special characters reporting errors on bash.
