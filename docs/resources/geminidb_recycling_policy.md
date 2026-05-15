---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_recycling_policy"
description: |-
  Manages a GeminiDB recycling policy resource within HuaweiCloud.
---

# huaweicloud_geminidb_recycling_policy

Manages a GeminiDB recycling policy resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_geminidb_recycling_policy" "test" {
  retention_period_in_days = 7
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Chenge this parameter will create a new resource.

* `retention_period_in_days` - (Required, Int) Specifies the period of retaining deleted GeminiDB instances.
  The value ranges from 1 to 7 days. The default value is 7 days.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the project ID of the region.

## Import

The GeminiDB recycling policy can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_geminidb_recycling_policy.test <id>
```
