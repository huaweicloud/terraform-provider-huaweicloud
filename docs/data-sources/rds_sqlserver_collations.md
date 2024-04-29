---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sqlserver_collations"
description: ""
---

# huaweicloud_rds_sqlserver_collations

Use this data source to get the list of RDS SQLServer collations.

## Example Usage

```hcl
data "huaweicloud_rds_sqlserver_collations" "collations" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `char_sets` - Indicates the character set information list.
