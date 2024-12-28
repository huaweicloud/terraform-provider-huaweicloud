---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice_engine_configurations"
description: |-
  Use this data source to query managed microservice engine configurations within HuaweiCloud.
---

# huaweicloud_cse_microservice_engine_configurations

Use this data source to query managed microservice engine configurations within HuaweiCloud.

## Example Usage

```hcl
variable "microservice_engine_id" {} // Enable the EIP access

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  fileter_engines = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == var.microservice_engine_id]
}

data "huaweicloud_cse_microservice_engine_configurations" "test" {
  auth_address    = local.fileter_engines[0].service_registry_addresses[0].public
  connect_address = local.fileter_engines[0].config_center_addresses[0].public
}
```

## Argument Reference

The following arguments are supported:

* `auth_address` - (Required, String) Specifies the address that used to request the access token.

* `connect_address` - (Required, String) Specifies the address that used to access engine and manages
  configuration.

-> We are only support IPv4 addresses yet (for `auth_address` and `connect_address`).

* `admin_user` - (Optional, String) Specifies the account name for **RBAC** login.

* `admin_pass` - (Optional, String) Specifies the account password for **RBAC** login.
  The password format must meet the following conditions:
  + Must be `8` to `32` characters long.
  + A password must contain at least one digit, one uppercase letter, one lowercase letter, and one special character
    (-~!@#%^*_=+?$&()|<>{}[]).
  + Cannot be the account name or account name spelled backwards.
  + The password can only start with a letter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - All queried configurations of the dedicated microservice engine.  
  The [configurations](#cse_microservice_engine_configurations) structure is documented below.

<a name="cse_microservice_engine_configurations"></a>
The `configurations` block supports:

* `id` - The ID of the microservice engine configuration.

* `key` - The configuration key (item name).

* `value_type` - The type of the configuration value.
  + **ini**
  + **json**
  + **text**
  + **yaml**
  + **properties**
  + **xml**

* `value` - The configuration value.

* `status` - The configuration status.
  + **enabled**
  + **disabled**

* `tags` - The key/value pairs to associate with the configuration.

* `created_at` - The creation time of the configuration, in RFC3339 format.

* `updated_at` - The latest update time of the configuration, in RFC3339 format.

* `create_revision` - The create version of the configuration.

* `update_revision` - The update version of the configuration.
