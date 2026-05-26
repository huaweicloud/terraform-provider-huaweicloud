---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_studio_instances"
description: |-
  Use this data source to get the list of DataArts Studio instances.
---

# huaweicloud_dataarts_studio_instances

Use this data source to get the list of DataArts Studio instances.

## Example Usage

```hcl
data "huaweicloud_dataarts_studio_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DataArts Studio instances are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of the DataArts Studio instances.
  The [instances](#dataarts_studio_instances_attr) structure is documented below.

<a name="dataarts_studio_instances_attr"></a>
The `instances` block supports:

* `id` - The ID of the DataArts Studio instance.

* `name` - The name of the DataArts Studio instance.

* `version` - The version (spec code) of the DataArts Studio instance.

* `order_id` - The order ID of the DataArts Studio instance.

* `product_id` - The product ID of the DataArts Studio instance.

* `auto_renew` - Whether auto renew is enabled for the DataArts Studio instance.

* `status` - The status of the DataArts Studio instance.

* `vpc_id` - The VPC ID of the DataArts Studio instance.

* `subnet_id` - The subnet ID of the DataArts Studio instance.

* `availability_zone` - The availability zone of the DataArts Studio instance.

* `enterprise_project_id` - The enterprise project ID to which the DataArts Studio instance belongs.

* `effective_time` - The effective time of the DataArts Studio instance, in RFC3339 format.

* `created_by` - The creator of the DataArts Studio instance.

* `created_at` - The creation time of the DataArts Studio instance, in RFC3339 format.

* `workspace_mode` - The workspace mode of the DataArts Studio instance.
