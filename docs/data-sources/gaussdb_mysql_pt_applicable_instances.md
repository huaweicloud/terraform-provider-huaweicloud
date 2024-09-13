---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_pt_applicable_instances"
description: |-
  Use this data source to get the GaussDB MySQL instances that a parameter template can be applied to.
---

# huaweicloud_gaussdb_mysql_pt_applicable_instances

Use this data source to get the GaussDB MySQL instances that a parameter template can be applied to.

## Example Usage

```hcl
data "huaweicloud_gaussdb_mysql_pt_applicable_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `configuration_id` - (Required, String) Specifies the parameter template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of instances that the parameter template can be applied to.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.
