---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_application_capacity_orders"
description: |-
  Use this data source to get the list of COC application capacity orders.
---

# huaweicloud_coc_application_capacity_orders

Use this data source to get the list of COC application capacity orders.

## Example Usage

```hcl
variable "application_id" {}

data "huaweicloud_coc_application_capacity_orders" "test" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  application_id     = var.application_id
}
```

## Argument Reference

The following arguments are supported:

* `cloud_service_name` - (Required, String) Specifies the cloud service name.
  The value can be **ecs**, **cce**, **rds** and so on.

* `type` - (Required, String) Specifies the resource type name.
  There are many resource types. Choose the resource type based on your business needs. Common resource types are as follows:
  + **cloudservers**: Elastic Cloud Servers.
  + **servers**: Bare Metal Servers.
  + **clusters**: Cloud Container Engines.
  + **instances**: Cloud Databases.

* `application_id` - (Optional, String) Specifies the application ID.

* `component_id` - (Optional, String) Specifies the component ID.

* `group_id` - (Optional, String) Specifies the group ID.

-> Exactly one of `group_id`, `component_id` or `application_id` must be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the application capacity order list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `type` - Indicates the type of capacity.

* `rank_list` - Indicates the list of TOP 5 ranked objects.

  The [rank_list](#data_rank_list_struct) structure is documented below.

<a name="data_rank_list_struct"></a>
The `rank_list` block supports:

* `id` - Indicates the application, component or group ID.

* `name` - Indicates the application, component or group name.

* `value` - Indicates the capacity value of an application, component or group.
