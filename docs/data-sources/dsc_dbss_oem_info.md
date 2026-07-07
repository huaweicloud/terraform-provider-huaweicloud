---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_dbss_oem_info"
description: |-
  Use this data source to get the list of DSC DBSS OEM info within HuaweiCloud.
---

# huaweicloud_dsc_dbss_oem_info

Use this data source to get the list of DSC DBSS OEM info within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_dbss_oem_info" "test" {
  type = "DBSS"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the instance type.
  Valid values are **GDE** (database encryption), **DOM** (database O&M),
  **DDM** (database dynamic masking), and **DBSS** (database audit).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `ins_info_list` - The instance information list.

  The [ins_info_list](#ins_info_list_struct) structure is documented below.

<a name="ins_info_list_struct"></a>
The `ins_info_list` block supports:

* `ins_id` - The instance ID.

* `ins_name` - The instance name.

* `related_datasource_policy_list` - The list of related data source and security policy information.

  The [related_datasource_policy_list](#related_datasource_policy_list_struct) structure is documented below.

* `subnet_id` - The subnet ID to which the instance belongs.

* `vpc_id` - The VPC ID to which the instance belongs.

<a name="related_datasource_policy_list_struct"></a>
The `related_datasource_policy_list` block supports:

* `datasource_address` - The data source address.

* `datasource_id` - The data source asset ID.

* `datasource_port` - The data source port.

* `ddm_policies` - The list of related dynamic masking policy information.

  The [ddm_policies](#ddm_policies_struct) structure is documented below.

* `gde_policies` - The list of related encryption policy information.

  The [gde_policies](#gde_policies_struct) structure is documented below.

* `proxy_port` - The proxy port.

* `sdm_policies` - The list of related static masking policy information.

  The [sdm_policies](#sdm_policies_struct) structure is documented below.

<a name="ddm_policies_struct"></a>
The `ddm_policies` block supports:

* `columns` - The list of column information configured by the policy.

  The [columns](#ddm_columns_struct) structure is documented below.

* `namespace` - The namespace name.

* `table` - The table name.

<a name="gde_policies_struct"></a>
The `gde_policies` block supports:

* `action` - The action type. 1 means encryption, 2 means decryption.

* `alg` - The encryption algorithm. Valid values are **sm4**, **aes-128**, and **aes-256**.

* `columns` - The list of column information to be encrypted.

  The [columns](#gde_columns_struct) structure is documented below.

* `table` - The table name.

<a name="sdm_policies_struct"></a>
The `sdm_policies` block supports:

* `columns` - The list of column information.

  The [columns](#sdm_columns_struct) structure is documented below.

* `do_mask` - Whether to mask the data.

* `do_move` - Whether to move the data.

* `namespace` - The namespace, specific to HBase.

* `table` - The table name.

<a name="ddm_columns_struct"></a>
<a name="gde_columns_struct"></a>
<a name="sdm_columns_struct"></a>
The `columns` block supports:

* `mask` - The masking algorithm name or ID.

* `name` - The column name.
