---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_associated_resource_status"
description: |-
  Use this data source to query the TMS associated resource status within HuaweiCloud.
---

# huaweicloud_tms_associated_resource_status

Use this data source to query the TMS associated resource tag inheritance status within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_tms_associated_resource_status" "test" {}

output "status" {
  value = data.huaweicloud_tms_associated_resource_status.test.status
}
```

## Argument Reference

There is no argument for the data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - The current status of the associated resource tag inheritance. The value can be **enabled** or **disabled**.