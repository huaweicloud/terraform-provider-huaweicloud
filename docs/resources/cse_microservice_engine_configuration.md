---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice_engine_configuration"
description: |-
  Manages a key/value pair under a dedicated microservice engine (2.0+) resource within HuaweiCloud.
---

# huaweicloud_cse_microservice_engine_configuration

Manages a key/value pair under a dedicated microservice engine (2.0+) resource within HuaweiCloud.

-> Before creating a configuration, make sure the engine has enabled the rules shown in the appendix
   [table](#configuration_default_engine_access_rules).

## Example Usage

### Create an engine configuration and the engine RBAC authentication is disabled

```hcl
variable "microservice_engine_id" {} // Enable the EIP access

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  fileter_engines = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == var.microservice_engine_id]
}

resource "huaweicloud_cse_microservice_engine_configuration" "test" {
  auth_address    = local.fileter_engines[0].service_registry_addresses[0].public
  connect_address = local.fileter_engines[0].config_center_addresses[0].public

  key        = "demo"
  value_type = "json"
  value      = jsonencode({
    "foo": "bar"
  })
  status     = "enabled"

  tags = {
    owner = "terraform"
  }
}
```

### Create an engine configuration and the engine RBAC authentication is enabled

```hcl
variable "microservice_engine_id" {} // Enable the EIP access
variable "microservice_engine_admin_password" {}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  fileter_engines = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == var.microservice_engine_id]
}

resource "huaweicloud_cse_microservice_engine_configuration" "test" {
  auth_address    = local.fileter_engines[0].service_registry_addresses[0].public
  connect_address = local.fileter_engines[0].config_center_addresses[0].public
  admin_user      = "root"
  admin_pass      = var.microservice_engine_admin_password

  key        = "demo"
  value_type = "json"
  value      = jsonencode({
    "foo": "bar"
  })
  status     = "enabled"

  tags = {
    owner = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `auth_address` - (Required, String, ForceNew) Specifies the address that used to request the access token.  
  Changing this will create a new resource.

* `connect_address` - (Required, String, ForceNew) Specifies the address that used to access engine and manages
  configuration.  
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

* `key` - (Required, String, ForceNew) Specifies the configuration key (item name).  
  The valid length is limited from `1` to `2,048` characters, only letters, digits, hyphens (-), underscores (_),
  colons (:) and periods (.) are allowed.  
  Changing this will create a new resource.

* `value_type` - (Required, String, ForceNew) Specifies the type of the configuration value.
  The valid values are as follows:
  + **ini**
  + **json**
  + **text**
  + **yaml**
  + **properties**
  + **xml**

  Changing this will create a new resource.

* `value` - (Required, String) Specifies the configuration value.

* `status` - (Optional, String) Specifies the configuration status.  
  The valid values are as follows:
  + **enabled**
  + **disabled**

* `tags` - (Optional, Map, ForceNew) Specifies the key/value pairs to associate with the configuration that used to
  filter resource.  
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `created_at` - The The creation time of the configuration, in RFC3339 format.

* `updated_at` - The latest update time of the configuration, in RFC3339 format.

* `create_revision` - The create version of the configuration.

* `update_revision` - The update version of the configuration.

## Import

If the related engine is disable the `RBAC`, configurations (key/value pairs) can be imported using their
`auth_address`, `connect_address` and `key`, e.g.

```bash
$ terraform import huaweicloud_cse_microservice_engine_configuration.test <auth_address>/<connect_address>/<key>
```

If you enabled the **RBAC** authorization in the microservice engine, it's necessary to provide the account
name (`admin_user`) and password (`admin_pass`) of the microservice engine.
All fields separated by the slashes (/), e.g.

```bash
$ terraform import huaweicloud_cse_microservice_engine_configuration.test <auth_address>/<connect_address>/<key>/<admin_user>/<admin_pass>
```

The single quotes (') or backslashes (\\) can help you solve the problem of special characters reporting errors on bash.

```bash
$ terraform import huaweicloud_cse_microservice_engine_configuration.test https://124.70.26.32:30100/https://124.70.26.32:30110/demo/root/Test\!123
```

```bash
$ terraform import huaweicloud_cse_microservice_engine_configuration.test 'https://124.70.26.32:30100/https://124.70.26.32:30110/demo/root/Test!123'
```

Note that the imported state may not be identical to your resource definition, due to security reason.
The missing attribute is `admin_pass`. It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the resource, or the definition should be updated to align with the
resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cse_microservice_engine_configuration" "test" {
  ...

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}
```

## Appendix

<a name="configuration_default_engine_access_rules"></a>
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
