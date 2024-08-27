---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_stream"
description: |-
  Manage a log stream resource within HuaweiCloud.
---

# huaweicloud_lts_stream

Manage a log stream resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}

resource "huaweicloud_lts_stream" "test" {
  group_id    = var.group_id
  stream_name = "testacc_stream"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the log stream resource. If omitted, the
  provider-level region will be used. Changing this creates a new log stream resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of a created log group. Changing this parameter will create
  a new resource.

* `stream_name` - (Required, String, ForceNew) Specifies the log stream name. Changing this parameter will create a new
  resource.

* `ttl_in_days` - (Optional, Int, ForceNew) Specifies the log expiration time(days), value range: 1-365.
  If not specified, it will inherit the log group setting. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs of the log stream.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The log stream ID.

* `filter_count` - Number of log stream filters.

* `created_at` - The creation time of the log stream.

## Import

The log stream can be imported using the group ID and stream ID separated by a slash, e.g.

```bash
$ terraform import huaweicloud_lts_stream.test <group_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to `ttl_in_days` attribute missing
from the API response. It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_lts_stream" "test" {
  ...

  lifecycle {
    ignore_changes = [
      ttl_in_days,
    ]
  }
}
