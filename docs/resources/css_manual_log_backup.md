---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_manual_log_backup"
description: |-
  Manages CSS manual log backup resource within HuaweiCloud.
---

# huaweicloud_css_manual_log_backup

Manages CSS manual log backup resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_css_manual_log_backup" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the ID of the CSS cluster.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `job_id` - The ID of the log backup job.

* `type` - The type of the log backup job.

* `status` - The status of the log backup job.

* `log_path` - The storage path of backed up logs in the OBS bucket.

* `created_at` - The creation time.

* `finished_at` - The end time.

* `failed_msg` - The error information.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
