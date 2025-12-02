---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_dictionaries"
description: |-
  Use this data source to get the list of dictionaries.
---

# huaweicloud_hss_setting_dictionaries

Use this data source to get the list of dictionaries.

## Example Usage

```hcl
variable "group_code" {}

data "huaweicloud_hss_setting_dictionaries" "test" {
  group_code = var.group_code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `group_code` - (Required, String) Specifies the dictionary group.
  The valid values are as follows:
  + **featureSwitch**: Page feature switch.

* `scene` - (Optional, String) Specifies the application scenario.
  The valid values are as follows:
  + **hws**: Chinese mainland website.
  + **hec-hk**: International website.

* `code` - (Optional, String) Specifies the dictionary item code.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of dictionaries.
  The [data_list](#hosts_data_list) structure is documented below.

<a name="hosts_data_list"></a>
The `data_list` block supports:

* `code` - The dictionary code.

* `value` - The dictionary value (single value).

* `values` - The dictionary values (multiple values).
