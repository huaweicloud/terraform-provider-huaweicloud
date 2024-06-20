---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ransomware_protection_policies"
description: |-
  Use this data source to get the list of HSS ransomware protection policies within HuaweiCloud.
---

# huaweicloud_hss_ransomware_protection_policies

Use this data source to get the list of HSS ransomware protection policies within HuaweiCloud.

## Example Usage

```hcl
variable policy_id {}

data "huaweicloud_hss_ransomware_protection_policies" "test" {
  policy_id = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS ransomware protection policies.
  If omitted, the provider-level region will be used.

* `policy_id` - (Optional, String) Specifies the ID of the ransomware protection policy to be queried.

* `name` - (Optional, String) Specifies the name of the ransomware protection policy to be queried.
  This field will undergo a fuzzy matching query, the query result is for all ransomware protection policies whose names
  contain this value.

* `operating_system` - (Optional, String) Specifies the operating system supported by the ransomware protection policies
  to be queried. The value can be **Windows** or **Linux**.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the ransomware
  protection policies belong. For enterprise users, if omitted, will query the ransomware protection policies under all
  enterprise projects.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `policies` - All ransomware protection policies that match the filter parameters.  
  The [policies](#hss_ransomware_protection_policies) structure is documented below.

<a name="hss_ransomware_protection_policies"></a>
The `policies` block supports:

* `id` - The ID of the ransomware protection policy.

* `name` - The name of ransomware protection policy.

* `protection_mode` - The protection mode of the ransomware protection policy.

* `bait_protection_status` - The bait protection enabled status of the ransomware protection policy.
  The value can be **opened** or **closed**.

* `deploy_mode` - The dynamic bait protection enabled status of the ransomware protection policy.
  The value can be **opened** or **closed**.

* `protection_directory` - The protection directory of the ransomware protection policy, multiple directories separated
  by semicolons.

* `protection_type` - The protection file type of the ransomware protection policy, multiple file types separated by
  commas.

* `exclude_directory` - The exclude directories of the ransomware protection policy, multiple directories separated by
  semicolons.

* `runtime_detection_status` - The runtime detection enabled status of the ransomware protection policy.
  The value can be **opened** or **closed**.

* `count_associated_server` - The number of hosts associated with the ransomware protection policy.

* `operating_system` - The operating system supported by the ransomware protection policy.

* `process_whitelist` - The process whitelist of the ransomware protection policy.
  The [process_whitelist](#hss_ransomware_process_whitelist) structure is documented below.

* `default_policy` - Is it the default policy.
  The value can be `1` or `0`. `1` represents the default policy, `0` represents the non-default policy.

<a name="hss_ransomware_process_whitelist"></a>
The `process_whitelist` block supports:

* `path` - The process path.

* `hash` - The process hash.
