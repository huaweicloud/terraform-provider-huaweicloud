---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_endpoint_groups"
description: |-
  Use this data source to get the list of endpoint groups.
---

# huaweicloud_ga_endpoint_groups

Use this data source to get the list of endpoint groups.

## Example Usage

```hcl
variable "endpoint_group_name" {}

data "huaweicloud_ga_endpoint_groups" "test" {
  name = var.endpoint_group_name
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_group_id` - (Optional, String) Specifies the ID of the endpoint group.

* `name` - (Optional, String) Specifies the name of the endpoint group.

* `status` - (Optional, String) Specifies the status of the endpoint group.
  The valid values are as follows:
  + **ACTIVE**: The status of the endpoint group is normal operation.
  + **ERROR**: The status of the endpoint group is error.

* `listener_id` - (Optional, String) Specifies the ID of the listener to which the endpoint group belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `endpoint_groups` - The list of the endpoint groups.
  The [endpoint_groups](#ga_endpoint_groups) structure is documented below.

<a name="ga_endpoint_groups"></a>
The `endpoint_groups` block supports:

* `id` - The ID of the endpoint group.

* `name` - The name of the endpoint group.  

* `description` - The description of the endpoint group.

* `status` - The status of the endpoint group.

* `traffic_dial_percentage` - The percentage of traffic distributed to the endpoint group.

* `region_id` - The region where the endpoint group belongs.

* `listener_id` - The ID of the listener to which the endpoint group belongs.

* `created_at` - The creation time of the endpoint group.

* `updated_at` - The latest update time of the endpoint group.

* `frozen_info` - The frozen details of cloud services or resources.
  The [frozen_info](#endpoint_groups_frozen_info) structure is documented below.

<a name="endpoint_groups_frozen_info"></a>
The `frozen_info` block supports:

* `status` - The status of a cloud service or resource.

* `effect` - The status of the resource after being forzen.

* `scene` - The service scenario.
