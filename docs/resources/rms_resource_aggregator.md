---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregator"
description: ""
---

# huaweicloud_rms_resource_aggregator

Manages a RMS aggregator resource within HuaweiCloud.

## Example Usage

### Account Based Aggregation

```hcl
variable "name" {}
variable "source_account_list" {
  type = list(string)
}

resource "huaweicloud_rms_resource_aggregator" "account" {
  name        = var.name
  type        = "ACCOUNT"
  account_ids = var.source_account_list
}
```

### Organization Based Aggregation

```hcl
variable "name" {}

resource "huaweicloud_rms_resource_aggregator" "organization" {
  name = var.name
  type = "ORGANIZATION"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the resource aggregator name.
  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the resource aggregator type, which can be **ACCOUNT** or **ORGANIZATION**.
  Changing this parameter will create a new resource.

* `account_ids` - (Optional, List) Specifies the source account list being aggregated.
  This parameter is only valid in **ACCOUNT** type.

* `tags` - (Optional, Map)  Specifies the key/value pairs to associate with the aggregator.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the resource aggregator identifier.

## Import

The aggregator can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rms_resource_aggregator.example 5dbcb2e0804f46cfabea2a6a1a68b0ae
```
