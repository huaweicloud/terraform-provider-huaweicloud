---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_log_setting"
description: ""
---

# huaweicloud_css_log_setting

Manages CSS log setting resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "agency" {}
variable "base_path" {}
variable "bucket" {}

resource "huaweicloud_css_log_setting" "test" {
  cluster_id = var.cluster_id
  agency     = var.agency
  base_path  = var.base_path
  bucket     = var.bucket
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies ID of the cluster whose log function you want to enable.
  Changing this creates a new resource.

* `agency` - (Required, String) Specifies the agency name. You can create an agency to allow CSS to
  call other cloud services.

* `base_path` - (Required, String) Specifies the storage path of backed up logs in the OBS bucket.

* `bucket` - (Required, String) Specifies the name of the OBS bucket for storing logs.

* `period` - (Optional, String) Specifies the backup start time. Format: GMT.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `updated_at` - The update time.

* `auto_enabled` - Whether to enable automatic backup.

* `log_switch` - Whether to enable the log function.

## Import

The CSS log setting can be imported using `cluster_id`, e.g.

```bash
$ terraform import huaweicloud_css_log_setting.test <cluster_id>
```
