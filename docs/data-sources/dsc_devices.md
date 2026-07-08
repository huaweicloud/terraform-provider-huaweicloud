---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_devices"
description: |-
  Use this data source to get the list of DSC devices within HuaweiCloud.
---

# huaweicloud_dsc_devices

Use this data source to get the list of DSC devices within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_devices" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `devices` - The device information list.

  The [devices](#devices_struct) structure is documented below.

<a name="devices_struct"></a>
The `devices` block supports:

* `id` - The device ID.

* `name` - The device name.

* `ip` - The device IP address.

* `type` - The device type.
  The valid values are as follows:
  + **0**: Application data audit.
  + **1**: Application data security gateway.
  + **2**: Database firewall.
  + **3**: Database encryption.
  + **4**: Static masking instance.
  + **5**: Dynamic masking instance.

* `status` - The device status.

* `mode` - The deployment mode.

* `version` - The device version.

* `description` - The device description.

* `manage_url` - The management URL.

* `vpc_id` - The VPC ID.

* `subnet_id` - The subnet ID.

* `create_time` - The creation time.

* `update_time` - The update time.

* `related_datasource_policy_list` - The related datasource policy information list.

  The [related_datasource_policy_list](#devices_related_datasource_policy_list_struct) structure is documented below.

<a name="devices_related_datasource_policy_list_struct"></a>
The `related_datasource_policy_list` block supports:

* `datasource_id` - The datasource asset ID.

* `datasource_address` - The datasource address.

* `datasource_port` - The datasource port.

* `proxy_port` - The proxy port.

* `ddm_policies` - The dynamic masking policy information list.

  The [ddm_policies](#devices_ddm_policies_struct) structure is documented below.

* `gde_policies` - The encryption policy information list.

  The [gde_policies](#devices_gde_policies_struct) structure is documented below.

* `sdm_policies` - The static masking policy information list.

  The [sdm_policies](#devices_sdm_policies_struct) structure is documented below.

<a name="devices_ddm_policies_struct"></a>
The `ddm_policies` block supports:

* `namespace` - The namespace name.

* `table` - The table name.

* `columns` - The column information list.

  The [columns](#devices_columns_struct) structure is documented below.

<a name="devices_gde_policies_struct"></a>
The `gde_policies` block supports:

* `action` - The action. 1 means encrypt, 2 means decrypt.

* `alg` - The encryption algorithm.

* `table` - The table name.

* `columns` - The column information list.

  The [columns](#devices_columns_struct) structure is documented below.

<a name="devices_sdm_policies_struct"></a>
The `sdm_policies` block supports:

* `table` - The table name.

* `namespace` - The namespace name.

* `do_mask` - Whether to mask data.

* `do_move` - Whether to move data.

* `columns` - The column information list.

  The [columns](#devices_columns_struct) structure is documented below.

<a name="devices_columns_struct"></a>
The `columns` block supports:

* `name` - The column name.

* `mask` - The masking algorithm name or ID.
