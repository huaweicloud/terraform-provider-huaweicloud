---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_bigdata_assets"
description: |-
  Use this data source to get the list of DSC bigdata assets within HuaweiCloud.
---

# huaweicloud_dsc_bigdata_assets

Use this data source to get the list of DSC bigdata assets within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_bigdata_assets" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the bigdata asset type. Valid values are **Elasticsearch**,
  **DLI**,**Hive**, **HBase**, **MRS_HIVE**, **ALL**, **LTS**,**HIVE_ONLY** and **JUST_BIG_DATA**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `assets` - The bigdata asset information list.

  The [assets](#assets_struct) structure is documented below.

<a name="assets_struct"></a>
The `assets` block supports:

* `asset_name` - The asset name.

* `authorize_fail_reason` - The authorization failure reason.

* `authorized` - Whether the asset is authorized.

* `create_time` - The creation timestamp of the asset.

* `default` - Whether the asset is default.

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

* `region` - The region information.

* `security_group_id` - The security group ID.

* `source_type` - The asset source type.

* `subnet_id` - The subnet ID.

* `vpc_id` - The VPC ID.
