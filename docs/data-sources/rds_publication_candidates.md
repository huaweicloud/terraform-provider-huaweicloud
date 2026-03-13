---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_publication_candidates"
description: |-
  Use this data source to get the available publications for an RDS SQL Server instance.
---

# huaweicloud_rds_publication_candidates

Use this data source to get the available publications for an RDS SQL Server instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_publication_candidates" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `publication_instance_id` - (Optional, String) Specifies the publisher instance ID.

* `publication_instance_name` - (Optional, String) Specifies the publisher instance name (fuzzy match).

* `publication_name` - (Optional, String) Specifies the publication name (fuzzy match).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_publications` - Indicates the list of available publications.

  The [instance_publications](#instance_publications_struct) structure is documented below.

<a name="instance_publications_struct"></a>
The `instance_publications` block supports:

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `publication_id` - Indicates the publication ID.

* `publication_name` - Indicates the publication name.
