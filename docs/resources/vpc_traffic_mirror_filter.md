---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_traffic_mirror_filter"
description: ""
---

# huaweicloud_vpc_traffic_mirror_filter

 Manages a VPC traffic mirror filter resource within HuaweiCloud.

## Example Usage

```hcl
variable "traffic_mirror_filter_name" {}

resource "huaweicloud_vpc_traffic_mirror_filter" "test" {
  name        = var.traffic_mirror_filter_name
  description = "Traffic mirror filter created by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the traffic mirror filter.

* `description` - (Optional, String) Specifies the description of the traffic mirror filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of the traffic mirror filter.

* `updated_at` - The latest update time of the traffic mirror filter.

## Import

The traffic mirror filter can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_traffic_mirror_filter.test <id>
```
