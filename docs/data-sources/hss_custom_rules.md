---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_custom_rules"
description: |-
  Use this data source to get the custom rules list of HSS within HuaweiCloud.
---

# huaweicloud_hss_custom_rules

Use this data source to get the custom rules list of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_custom_rules" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, String) Specifies the rule ID.

* `rule_name` - (Optional, String) Specifies the rule name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The custom rules list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `rule_id` - The rule ID.

* `host_num` - The number of protected hosts.

* `rule_name` - The rule name.

* `rule_status` - The rule status.  
  The valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

* `rule_type` - The rule type.  
  The valid values are as follows:
  + **black_hash**: Black hash.

* `auto_block` - Whether to automatically block alerts.  
  The valid values are as follows:
  + `0`: Do not automatically block alerts.
  + `1`: Automatically block alerts.

* `hash_type` - The hash type.  
  Valid values are: **sha256**, **md5**, and **sha1**.

* `is_all_host` - Whether to select all hosts.

* `create_time` - The creation time in milliseconds.

* `update_time` - The update time in milliseconds.
