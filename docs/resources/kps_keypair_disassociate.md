---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_keypair_disassociate"
description: |-
  Manages a KPS keypair disassociate resource within HuaweiCloud.
---

# huaweicloud_kps_keypair_disassociate

Manages a KPS keypair disassociate resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

-> Please refer to the API document link for more precautions  
[reference](https://support.huaweicloud.com/intl/en-us/usermanual-dew/dew_01_0071.html).

## Example Usage

### Disassociate SSH keypair with ECS opened

```hcl
variable "server_id" {}
variable "port" {}
variable "type" {}
variable "key" {}

resource "huaweicloud_kps_keypair_disassociate" "test" {
  server {
    id   = var.server_id
    port = var.port
    
    auth {
      type = var.type
      key  = var.key
    }
  }
}
```

### Disassociate SSH keypair with ECS closed

```hcl
variable "server_id" {}

resource "huaweicloud_kps_keypair_disassociate" "test" {
  server {
    id = var.server_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `server` - (Required, List, NonUpdatable) Specifies the ECS information that requires disassociating keypair.
  The [server](#kps_server) structure is documented below.

<a name="kps_server"></a>
The `server` block supports:

* `id` - (Required, String, NonUpdatable) Specifies ID of the ECS which need to disassociate the SSH keypair.

* `auth` - (Optional, List, NonUpdatable) Specifies the authentication information. This parameter is required for replacement
  and not required for reset.
  The [auth](#server_auth) structure is documented below.

* `port` - (Optional, Int, NonUpdatable) Specifies the SSH listening port. The default value is `22`.

-> When the ECS is shut down, the disassociate `port` is fixed at `22` and cannot be configured.
At the same time, `auth` can not be set. When the ECS is turned on, the disassociate `port` can be
configured and defaults to `22`, and `auth` is required, otherwise the operation will fail.

<a name="server_auth"></a>
The `auth` block supports:

* `type` - (Optional, String, NonUpdatable) Specifies the value of an authentication type.
  The valid values are **password** and **keypair**.

* `key` - (Optional, String, NonUpdatable) Specifies the value of the key depending on the `type`.
  When the `type` is set to **password**, it represents the password.  
  When the `type` is set to **keypair**, it represents the private key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is also the task ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Defaults to 30 minutes.
