---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_client_auth_config_history"
description: |-
  Use this data source to query the client access authentication configuration modification history of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_client_auth_config_history

Use this data source to query the client access authentication configuration modification history of a GaussDB instance
within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_client_auth_config_history" "test" {
  instance_id = var.instance_id
}
```

### Filter by Time Range

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_client_auth_config_history" "test" {
  instance_id = var.instance_id
  start_time  = "2026-01-01 00:00:00"
  end_time    = "2026-12-31 23:59:59"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the client auth config history.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID of the GaussDB instance.  
  This parameter is the unique identifier of the instance created by the user.

* `start_time` - (Optional, String) Specifies the start time of the query interval.  
  The format is **yyyy-mm-dd hh:mm:ss**.

* `end_time` - (Optional, String) Specifies the end time of the query interval.  
  The format is **yyyy-mm-dd hh:mm:ss**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `hba_histories` - The list of client access authentication configuration modification histories.  
  The [hba_histories](#gaussdb_client_auth_config_history_hba_histories) structure is documented below.

<a name="gaussdb_client_auth_config_history_hba_histories"></a>
The `hba_histories` block supports:

* `id` - The ID of the client access authentication modification record.

* `status` - The status of the client access authentication modification.  
  The valid values are as follows:
  + **success**: The modification was successful.
  + **failed**: The modification failed.

* `time` - The modification time. The format is **yyyy-mm-dd hh:mm:ss**.

* `fail_reason` - The reason for the modification failure. This parameter is only returned when the modification fails.

* `before_confs` - The client access authentication configuration before modification.  
  The [before_confs](#gaussdb_client_auth_config_history_hba_confs) structure is documented below.

* `after_confs` - The client access authentication configuration after modification.  
  The [after_confs](#gaussdb_client_auth_config_history_hba_confs) structure is documented below.

<a name="gaussdb_client_auth_config_history_hba_confs"></a>
The `before_confs` and `after_confs` block supports:

* `type` - The client connection type.  
  The valid values are as follows:
  + **host**: Accepts both normal TCP/IP socket connections and SSL-encrypted TCP/IP socket connections.
  + **hostssl**: Only accepts SSL-encrypted TCP/IP socket connections.
  + **hostnossl**: Only accepts normal TCP/IP socket connections.

* `database` - The database that the record matches and allows access to.  
  The value can be **all** (matches all databases) or a specific database name.

* `user` - The database user that the record matches and allows access to.  
  The value can be **all** (matches all users) or a specific username.

* `address` - The IP address range that matches and is allowed to access.  
  Supports IPv4 and IPv6, e.g., `10.10.0.0/24` or `2001:250::/128`.

* `method` - The authentication method used for the connection.  
  The valid values are as follows:
  + **reject**: Unconditionally reject the connection.
  + **md5**: Perform MD5 password authentication.
  + **sha256**: Perform SHA-256 password authentication.
  + **sm3**: Perform SM3 password authentication.
  + **gss**: Use GSSAPI authentication.
