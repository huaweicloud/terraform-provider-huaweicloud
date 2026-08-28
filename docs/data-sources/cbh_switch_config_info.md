---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_switch_config_info"
description: |-
  Use this data source to query the CBH switch config info within HuaweiCloud.
---

# huaweicloud_cbh_switch_config_info

Use this data source to query the CBH switch config info within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_cbh_switch_config_info" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the switch config info.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `switch_info` - The switch information of the CBH service.
  The [switch_info](#cbh_switch_config_info_switch_info) structure is documented below.

* `version_info` - The version information of the CBH service.
  The [version_info](#cbh_switch_config_info_version_info) structure is documented below.

<a name="cbh_switch_config_info_switch_info"></a>
The `switch_info` block supports:

* `is_support_unibuy` - Whether unibuy is supported.
* `is_support_float_ipv6` - Whether floating IPv6 is supported.
* `is_support_admin_login` - Whether admin login is supported.
* `is_support_update_ha` - Whether HA update is supported.
* `is_support_tms` - Whether TMS is supported.
* `is_support_eps` - Whether EPS is supported.
* `is_support_iam_login` - Whether IAM login is supported.
* `is_support_ipv6` - Whether IPv6 is supported.
* `is_support_ha` - Whether HA is supported.
* `is_support_reset` - Whether reset admin password and admin login mode is supported.
* `is_support_upgrade_instance` - Whether instance upgrade is supported.
* `is_support_change_security_group` - Whether security group change is supported.
* `is_support_manually_ip` - Whether manual IP is supported.
* `is_support_capacity_expantion` - Whether capacity expansion is supported.
* `is_support_ha_expantion` - Whether HA expansion is supported.
* `is_support_agency_authorize` - Whether agency authorization is supported.
* `is_support_change_vpc` - Whether VPC change is supported.
* `is_support_cluster` - Whether cluster is supported.
* `is_support_ondemand` - Whether on-demand billing is supported.
* `is_support_period` - Whether period billing is supported.

<a name="cbh_switch_config_info_version_info"></a>
The `version_info` block supports:

* `require_eip` - The instance version that supports EIP.
* `iam_login` - The instance version that supports IAM login.
* `admin_login` - The instance version that supports admin login.
* `float_ipv6` - The instance version that supports floating IPv6.
