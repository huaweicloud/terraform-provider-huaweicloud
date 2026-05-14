---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_redis_log_download_link"
description: |-
  Use this data source to obtain the download link for a specific Redis log of a DCS instance within HuaweiCloud.
---

# huaweicloud_dcs_redis_log_download_link

Use this data source to obtain the download link for a specific Redis log of a DCS instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "log_id" {}

data "huaweicloud_dcs_redis_log_download_link" "test" {
  instance_id = var.instance_id
  log_id      = var.log_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the log download link. If omitted, the
  provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `log_id` - (Required, String) Specifies the ID of the Redis log.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `link` - The download link for the Redis log.

* `backup_id` - The background task ID associated with the download link creation.
