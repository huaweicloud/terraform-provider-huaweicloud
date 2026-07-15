---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_bigdata_instance_databases"
description: |-
  Use this data source to get the database list of a bigdata instance within HuaweiCloud.
---

# huaweicloud_dsc_bigdata_instance_databases

Use this data source to get the database list of a bigdata instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dsc_bigdata_instance_databases" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datasources` - The bigdata asset detail list.

  The [datasources](#datasources_struct) structure is documented below.

<a name="datasources_struct"></a>
The `datasources` block supports:

* `asset_name` - The asset name.

* `authorize_fail_reason` - The authorization failure reason.

* `authorized` - Whether the asset is authorized.

* `create_time` - The asset creation time.

* `default` - Whether the asset is the default asset.

* `ds_address` - The data source address.

* `ds_authorized` - The data source authorization status.

* `ds_name` - The data source name.

* `ds_port` - The data source port.

* `ds_type` - The data source type.

* `ds_user` - The data source username.

* `ds_version` - The data source version.

* `id` - The asset ID.

* `ins_id` - The instance ID.

* `ins_name` - The instance name.

* `ins_type` - The instance type.

* `region` - The region of the asset.

* `security_group_id` - The security group ID.

* `source_type` - The asset source type.

* `subnet_id` - The subnet ID.

* `vpc_id` - The VPC ID.
