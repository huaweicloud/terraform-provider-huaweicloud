---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_audit_log_delete"
description: |-
  Manages a DDS audit log delete resource within HuaweiCloud.
---

# huaweicloud_dds_audit_log_delete

Manages a DDS audit log delete resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "file_names" {}

resource "huaweicloud_dds_audit_log_delete" "test" {
  instance_id = var.instance_id
  file_names  = var.file_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID.
  Changing this creates a new resource.

* `file_names` - (Required, List, ForceNew) Specifies the audit log file names.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 50 minutes.
