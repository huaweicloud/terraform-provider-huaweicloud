---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_resource_group_service_resources"
description: |-
  Use this data source to get the list of resources under a specific service category within a CES resource group.
---

# huaweicloud_ces_resource_group_service_resources

Use this data source to get the list of resources under a specific service category within a CES resource group.

## Example Usage

```hcl
variable "group_id" {}
variable "service" {}

data "huaweicloud_ces_resource_group_service_resources" "test" {
  group_id = var.group_id
  service  = var.service
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `group_id` - (Required, String) Specifies the resource group ID.

* `service` - (Required, String) Specifies the service type.

* `dim_name` - (Optional, String) Specifies the dimension name.

* `dim_value` - (Optional, String) Specifies the dimension value. Fuzzy match is not supported.
  You can specify only one dimension for resources with multiple dimensions.

* `status` - (Optional, String) Specifies the health status.
  The valid values are as follows:
  + **health**: resources for which alarm rules have been configured and no alarm was triggered.
  + **unhealthy**: resources for which alarm rules have been configured and alarms were triggered.
  + **no_alarm_rule**: resources for which alarm rules are not configured.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The resources in a resource group.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `status` - The health status.

* `dimensions` - The dimension information about a resource.

  The [dimensions](#resources_dimensions_struct) structure is documented below.

<a name="resources_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The dimension name.

* `value` - The dimension value.
