---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_recycling_policy"
description: |-
  Manage a GaussDB MySQL recycling policy resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_recycling_policy

Manage a GaussDB MySQL recycling policy resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_gaussdb_mysql_recycling_policy" "test" {
  retention_period_in_days = "5"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the recycling policy resource. If omitted, the
  provider-level region will be used. Changing this creates a new GaussDB MySQL instance resource.

* `retention_period_in_days` - (Required, String) Specifies the retention period, in days. Value ranges: **1** to **7**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID. The value is the project ID of the region.

## Import

The GaussDB MySQL recycling policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_recycling_policy.test <id>
```
