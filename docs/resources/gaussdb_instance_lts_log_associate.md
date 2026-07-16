---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_lts_log_associate"
description: |-
  Manages a GaussDB instance LTS log associate resource within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_lts_log_associate

Manages a GaussDB instance LTS log associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "lts_group_id" {}
variable "lts_stream_id" {}

resource "huaweicloud_gaussdb_instance_lts_log_associate" "test" {
  instance_id   = var.instance_id
  log_type      = "audit_log"
  lts_group_id  = var.lts_group_id
  lts_stream_id = var.lts_stream_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the LTS log associate resource.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID of the GaussDB instance.  
  This parameter is the unique identifier of the instance created by the user.

* `log_type` - (Required, String, NonUpdatable) Specifies the log type.  
  The valid value is **audit_log**.

* `lts_group_id` - (Required, String) Specifies the LTS log group ID.  
  You can obtain it through the LTS API "Query All Log Groups Under an Account".

* `lts_stream_id` - (Required, String) Specifies the LTS log stream ID.  
  You can obtain it through the LTS API "Query All Log Streams Under a Specified Log Group".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, composed of `instance_id` and `log_type`, separated by a slash (`/`).

* `enabled` - Whether the LTS log association is enabled.

## Import

The GaussDB instance LTS log associate resource can be imported using `instance_id` and `log_type` separated by a
slash (`/`), e.g.

```bash
$ terraform import huaweicloud_gaussdb_instance_lts_log_associate.test <instance_id>/<log_type>
```
