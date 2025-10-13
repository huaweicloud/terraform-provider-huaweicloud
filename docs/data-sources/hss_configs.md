---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_configs"
description: |-
  Use this data source to get the list of HSS configs within HuaweiCloud.
---

# huaweicloud_hss_configs

Use this data source to get the list of HSS configs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_configs" "test" {
  config_name_list = ["password_min_len"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `config_name_list` - (Required, List) Specifies the list of config names.  
  The valid values are as follows:
  + **password_min_len**: Minimum password length.
  + **password_digit_min_num**: Minimum number of digits in password.
  + **password_upper_letter_min_num**: Minimum number of uppercase letters in password.
  + **password_lower_letter_min_num**: Minimum number of lowercase letters in password.
  + **password_special_char_min_num**: Minimum number of special characters in password.
  + **weakpwd**: Customize weak password strategy.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of config information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `config_name` - The config name.

* `config_value` - The config content.
