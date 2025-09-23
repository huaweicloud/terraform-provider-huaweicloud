---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_databases"
description: |-
  Use this data source to get the databases of a specified GaussDB OpenGauss instance.
---

# huaweicloud_gaussdb_opengauss_databases

Use this data source to get the databases of a specified GaussDB OpenGauss instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_databases" "this" {
  instance_id = var.instance_id
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the instance. If omitted, the provider-level region will
  be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the data source.

* `databases` - Indicates the list of the databases.

  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `name` - Indicates the database name.

* `owner` - Indicates the database user.

* `character_set` - Indicates the database character set.

* `lc_collate` - Indicates the database collation.

* `size` - Indicates the database size.

* `compatibility_type` - Indicates the database compatibility type.
