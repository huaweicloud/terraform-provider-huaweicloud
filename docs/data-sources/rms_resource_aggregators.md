---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregators"
description: |-
  Use this data source to get the list of RMS resource aggregator.
---

# huaweicloud_rms_resource_aggregators

Use this data source to get the list of RMS resource aggregator.

## Example Usage

```hcl
variable "aggregator_name" {}

data "huaweicloud_rms_resource_aggregator" "test" {
  name = var.aggregator_name
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Optional, String) Specifies the resource aggregator ID.

* `name` - (Optional, String) Specifies the resource aggregator name.

* `type` - (Optional, String) Specifies the resource aggregator type, which can be ACCOUNT or ORGANIZATION.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `aggregators` - The resource aggregators.

  The [aggregators](#aggregators_struct) structure is documented below.

<a name="aggregators_struct"></a>
The `aggregators` block supports:

* `name` - The resource aggregator name.

* `id` - The resource aggregator ID.

* `type` - The resource aggregator type.

* `urn` - The resource aggregator identifier.

* `account_ids` - The source account list being aggregated.

* `created_at` - The time when the resource aggregator was created.

* `updated_at` - The time when the resource aggregator was updated.
