---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_recycling_policy"
description: |-
Manage an RDS recycling policy resource within HuaweiCloud.
---

# huaweicloud_rds_recycling_policy

Manage an RDS recycling policy resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_rds_recycling_policy" "test" {
  retention_period_in_days = 5
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the recycling policy resource. If omitted, the
  provider-level region will be used. Changing this creates a new rds instance resource.

* `retention_period_in_days` - (Optional, Int) Specifies the period of retaining deleted DB instances. Value ranges
  from `1` day to `7` days. Defaults to `7`.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - Indicates the resource ID. The value is the project ID of the region.

## Import

The RDS recycling policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_recycling_policy.test <id>
```
