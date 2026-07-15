---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_authorize_databases"
description: |-
  Use this data source to get the list of authorized databases for database instances within HuaweiCloud.
---

# huaweicloud_dsc_authorize_databases

Use this data source to get the list of authorized databases for database instances within HuaweiCloud.

## Example Usage

```hcl
variable "instance_type" {}

data "huaweicloud_dsc_authorize_databases" "test" {
  instance_type = var.instance_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_type` - (Required, String) Specifies the instance type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `databases` - The authorized database information list.
  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `id` - The database ID.

* `asset_name` - The asset name.

* `auth_type` - The authorization type.

* `authorize_fail_reason` - The authorization failure reason.

* `authorized` - Whether the database is authorized.

* `create_time` - The creation time.

* `db_address` - The database address.

* `db_authorized` - The database authorization status.

* `db_name` - The database name.

* `db_port` - The database port.

* `db_type` - The database type.

* `db_user` - The database username.

* `db_version` - The database version.

* `default` - Whether it is the default database.

* `ins_id` - The instance ID.

* `ins_name` - The instance name.

* `ins_type` - The instance type.

* `region` - The region where the instance is located.

* `service_name` - The service name.

* `sid` - The session ID.

* `subnet_ids` - The subnet ID list.

* `vpc_ids` - The VPC ID list.
