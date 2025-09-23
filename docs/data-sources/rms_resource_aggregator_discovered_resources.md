---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregator_discovered_resources"
description: |-
  Use this data source to get the list of RMS resource aggregator discovered resources.
---

# huaweicloud_rms_resource_aggregator_discovered_resources

Use this data source to get the list of RMS resource aggregator discovered resources.

## Example Usage

```hcl
variable "aggregator_id" {}

data "huaweicloud_rms_resource_aggregator_discovered_resources" "test" {
  aggregator_id = var.aggregator_id
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, String) Specifies the aggregator ID.

* `service_type` - (Optional, String) Specifies the service type. For example, **vpc**.

* `resource_type` - (Optional, String) Specifies the resource type. For example, **vpcs**.

* `filter` - (Optional, List) Specifies the filter. The [filter](#filter) structure is documented below.

<a name="filter"></a>
The `filter` block supports:

* `account_id` - (Optional, String) Specifies the ID of account to which the resource belongs.

* `region_id` - (Optional, String) Specifies the ID of region to which the resource belongs.

* `resource_id` - (Optional, String) Specifies resource ID.

* `resource_name` - (Optional, String) Specifies resource name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The service details list.

  The [resources](#resources) structure is documented below.

<a name="resources"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `service` - The service name.

* `type` - The resource type.

* `region_id` - The region to which the resource belongs.

* `source_account_id` - The ID of the account to which the resource belongs.
