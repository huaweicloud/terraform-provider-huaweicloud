---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_proxy_flavors"
description: |-
  Use this data source to get the list of GaussDB MySQL proxy flavors.
---

# huaweicloud_gaussdb_mysql_proxy_flavors

Use this data source to get the list of GaussDB MySQL proxy flavors.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_mysql_proxy_flavors" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of GaussDB MySQL Instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavor_groups` - Indicates the list of flavor groups.

  The [flavor_groups](#flavor_groups_struct) structure is documented below.

<a name="flavor_groups_struct"></a>
The `flavor_groups` block supports:

* `type` - Indicates the group type. The value can be **arm** or **x86**.

* `flavors` - Indicates the list of flavors.

  The [flavors](#flavor_groups_flavors_struct) structure is documented below.

<a name="flavor_groups_flavors_struct"></a>
The `flavors` block supports:

* `id` - Indicates the ID of the proxy flavor.

* `db_type` - Indicates the database type.

* `vcpus` - Indicates the number of vCPUs.

* `ram` - Indicates the memory size in GB.

* `spec_code` - Indicates the proxy specification code.

* `az_status` - Indicates the key/value pairs of the availability zone status.
  **key** indicates the AZ ID, and **value** indicates the specification status in the AZ.
