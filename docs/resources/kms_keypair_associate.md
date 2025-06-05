---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_keypair_associate"
description: |-
  Manages a KMS keypair associate resource within HuaweiCloud.
---

# huaweicloud_kms_keypair_associate

Manages a KMS keypair associate resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

-> There are several prerequisites for using this resource:  
  <br/>1. The image ECS used must be the puclic image provided by HuaweiCloud.  
  <br/>2. The operation of keypair is achieved by modifying the server's `/root/.ssh/authorized_keys` file.
  Please ensure that the file has not been modified before using the keypair, otherwise the operation will fail.  
  <br/>3. The SSH port of the ECS security group needs to be opened for network segment `100.125.0.0/16` in advance.  
  <br/>4. There can be at most `10` ECS associating keypair simultaneously.

-> Please refer to the API document link for more precautions  
  [reference](https://support.huaweicloud.com/intl/en-us/usermanual-dew/dew_01_0071.html).

## Example Usage

### Associate/Replace SSH keypair with ECS opened

```hcl
variable "keypair_name" {}
variable "server_id" {}
variable "port" {}
variable "type" {}
variable "key" {}

resource "huaweicloud_kms_keypair_associate" "test" {
  keypair_name  = var.keypair_name

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

### Associate/Reset SSH keypair with ECS closed

```hcl
variable "keypair_name" {}
variable "server_id" {}

resource "huaweicloud_kms_keypair_associate" "test" {
  keypair_name  = var.keypair_name

  server {
    id = var.server_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `keypair_name` - (Required, String, NonUpdatable) Specifies the name of SSH keypair.

* `server` - (Required, List, NonUpdatable) Specifies the ECS information that requires associating keypair.
  These public images which are not supported to associate keypair are as follows:  
  **CoreOS**, **OpenEuler**, **FreeBSD（Other）**, **Kylin V10 64bit**, **UnionTech OS Server 20**,  
  **Euler 64bit** and **CentOS Stream 8 64bit**.

  The [server](#kms_server) structure is documented below.

<a name="kms_server"></a>
The `server` block supports:

* `id` - (Required, String, NonUpdatable) Specifies ID of the ECS which need to associate (replace or reset) the SSH keypair.

* `auth` - (Optional, List) Specifies the authentication information.
  The [auth](#server_auth) structure is documented below.

* `port` - (Optional, Int, NonUpdatable) Specifies the SSH listening port. The default value is `22`.

-> When the ECS is shut down, the operation (associate, disassociate, reset) `port` is fixed at `22` and cannot be configured.
At the same time, `auth` can not be set. When the ECS is turned on, the operation (associate,replace) `port` can be configured
and defaults to `22`, and `auth` is required, otherwise the operation will fail.

* `disable_password` - (Optional, Bool, NonUpdatable) Specifies whether the password is disabled.
  The valid values are as follows:  
  + **true**: Indicates disable SSH login for virtual machines.  
  + **false**: Indicates enable SSH login for virtual machines. Defaults to **false**.

<a name="server_auth"></a>
The `auth` block supports:

* `type` - (Optional, String, NonUpdatable) Specifies the value of the authentication type.
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
