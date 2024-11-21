---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_proxy_flavors"
description: |-
  Use this data source to get the list of RDS MySQL proxy flavors.
---

# huaweicloud_rds_mysql_proxy_flavors

Use this data source to get the list of RDS MySQL proxy flavors.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_mysql_proxy_flavors" "flavor" {
  instance_id = var.instance_id
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of RDS MySQL instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavor_groups` - Indicates the list of flavor groups.
  The [flavor_groups](#flavor_groups_struct) structure is documented below.

<a name="flavor_groups_struct"></a>
The `flavor_groups` block supports:

* `group_type` - Indicates the specification group type. The value can be **ARM** or **X86**.

* `flavors` - Indicates the list of flavors.

  The [flavors](#flavor_groups_flavors_struct) structure is documented below.

<a name="flavor_groups_flavors_struct"></a>
The `flavors` block supports:

* `id` - Indicates the specification ID of the database proxy.

* `code` - Indicates the specification code of the database proxy.

* `vcpus` - Indicates the number of vCPUs.

* `memory` - Indicates the memory size in GB.

* `db_type` - Indicates the database type.

* `az_status` - Indicates the AZ information. **key** indicates the AZ associated with the specification, and **value**
  indicates the specification status in the AZ. Only the specification status in the AZ where the primary instance is
  located is displayed.
