---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_client_auth_config_restore"
description: |-
  Use this resource to restore the client access authentication configuration of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_client_auth_config_restore

Use this resource to restore the client access authentication configuration of a GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_client_auth_config_restore" "test" {
  instance_id = var.instance_id
}
```

### Filter by hba_history_id

```hcl
variable "instance_id" {}
variable "hba_history_id" {}

resource "huaweicloud_gaussdb_client_auth_config_restore" "test" {
  instance_id    = var.instance_id
  hba_history_id = var.hba_history_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to restore the client access authentication
  configuration.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `hba_history_id` - (Optional, String, NonUpdatable) Specifies the client access authentication modification history
  record ID. If empty, it means restoring to the default configuration.
