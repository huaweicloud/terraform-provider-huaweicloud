---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregator_resource_detail"
description: |-
  Use this data source to get the detail about a specific resource in a source account.
---

# huaweicloud_rms_resource_aggregator_resource_detail

Use this data source to get the detail about a specific resource in a source account.

## Example Usage

```hcl
variable "aggregator_id" {}
variable "resource_id" {}
variable "source_account_id" {}
variable "region_id" {}

data "huaweicloud_rms_resource_aggregator_resource_detail" "test" {
  aggregator_id      = var.aggregator_id
  resource_id        = var.resource_id
  cloud_service_type = "vpc"
  type               = "securityGroups"
  source_account_id  = var.source_account_id
  region_id          = var.region_id
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, String) Specifies the resource aggregator ID.

* `resource_id` - (Required, String) Specifies the resource ID.

* `service_type` - (Required, String) Specifies the cloud service type.

* `type` - (Required, String) Specifies the resource type.

* `source_account_id` - (Required, String) Specifies the source account ID.

* `region_id` - (Required, String) Specifies the region to which the resource belongs.

* `resource_name` - (Optional, String) Specifies the resource name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `aggregator_domain_id` - Indicates the resource aggregator account.

* `ep_id` - Indicates the enterprise project ID.

* `created` - Indicates the time when the resource was created.

* `updated` - Indicates the time when the resource was updated.

* `tags` - Indicates the resource tag.

* `properties` - Indicates the properties of the resource.
