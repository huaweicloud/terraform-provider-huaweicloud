---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregator_source_statuses"
description: |-
  Use this data source to get the list of RMS aggregator source statuses.
---

# huaweicloud_rms_resource_aggregator_source_statuses

Use this data source to get the list of RMS aggregator source statuses.

## Example Usage

```hcl
variable "aggregator_id" {}

data "huaweicloud_rms_resource_aggregator_source_statuses" "test" {
  aggregator_id = var.aggregator_id
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, String) Specifies the resource aggregator ID.

* `status` - (Optional, String) Specifies the status of the aggregated source account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `aggregated_source_statuses` - The list of source statuses

  The [aggregated_source_statuses](#aggregated_source_statuses_struct) structure is documented below.

<a name="aggregated_source_statuses_struct"></a>
The `aggregated_source_statuses` block supports:

* `source_name` - The source name.

* `source_id` - The source ID.
  The value can be an account ID of organizatin ID.

* `source_type` - The source account type.
  The value can be **ACCOUNT** or **ORGANIZATION**.

* `last_error_code` - The error code returned when the last resource aggregation for the source fails.

* `last_error_message` - The error message returned when the last resource aggregation for the source fails.

* `last_update_status` - The latest status of the source.

* `last_update_time` - The last update time of the source.
