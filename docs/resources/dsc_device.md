---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_device"
description: |-
  Manages DSC device resource within HuaweiCloud.
---

# huaweicloud_dsc_device

Manages DSC device resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_dsc_device" "test" {
  name        = var.name
  type        = 0
  mode        = "SINGLE"
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  description = "demo description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the device name.

* `type` - (Required, Int) Specifies the device type.
  Valid values are: `0` (application data audit), `1` (application data security gateway),
  `2` (database firewall), `3` (database encryption), `4` (static masking instance),
  `5` (dynamic masking instance).

* `mode` - (Required, String) Specifies the deployment mode.

* `vpc_id` - (Required, String) Specifies the VPC ID.

* `subnet_id` - (Required, String) Specifies the subnet ID.

* `description` - (Optional, String) Specifies the description of the device.

* `manage_url` - (Optional, String) Specifies the management URL.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also device ID).

* `ip` - The device IP address.

* `status` - The device status.

* `version` - The device version.

* `create_time` - The creation time of the device.

* `update_time` - The update time of the device.

* `related_datasource_policy_list` - The related datasource policy information list.
  The [related_datasource_policy_list](#related_datasource_policy_list_struct) structure is documented below.

<a name="related_datasource_policy_list_struct"></a>
The `related_datasource_policy_list` block supports:

* `datasource_id` - The datasource asset ID.

* `datasource_address` - The datasource address.

* `datasource_port` - The datasource port.

* `proxy_port` - The proxy port.

* `ddm_policies` - The dynamic masking policy information list.
  The [ddm_policies](#ddm_policies_struct) structure is documented below.

* `gde_policies` - The encryption policy information list.
  The [gde_policies](#gde_policies_struct) structure is documented below.

* `sdm_policies` - The static masking policy information list.
  The [sdm_policies](#sdm_policies_struct) structure is documented below.

<a name="ddm_policies_struct"></a>
The `ddm_policies` block supports:

* `namespace` - The namespace name.

* `table` - The table name.

* `columns` - The column information list.
  The [columns](#columns_struct) structure is documented below.

<a name="gde_policies_struct"></a>
The `gde_policies` block supports:

* `action` - The action. **1** means encrypt, **2** means decrypt.

* `alg` - The encryption algorithm. Valid values are: **sm4**, **aes-128**, **aes-256**.

* `table` - The table name.

* `columns` - The column information list.
  The [columns](#columns_struct) structure is documented below.

<a name="sdm_policies_struct"></a>
The `sdm_policies` block supports:

* `table` - The table name.

* `namespace` - The namespace name.

* `do_mask` - Whether to mask data.

* `do_move` - Whether to move data.

* `columns` - The column information list.
  The [columns](#columns_struct) structure is documented below.

<a name="columns_struct"></a>
The `columns` block supports:

* `name` - The column name.

* `mask` - The masking algorithm name or ID.

## Import

The DSC device can be imported by `id`. e.g.

```bash
$ terraform import huaweicloud_dsc_device.test <id>
```
