---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_configuration_dictionaries"
description: |-
  Use this data source to get the list of dictionaries.
---

# huaweicloud_secmaster_configuration_dictionaries

Use this data source to get the list of dictionaries.

## Example Usage

```hcl
data "huaweicloud_secmaster_configuration_dictionaries" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `success_list` - The list of normal dictionaries.  
  The [success_list](#success_list_struct) structure is documented below.

* `failed_list` - The list of abnormal dictionaries.  
  The [failed_list](#failed_list_struct) structure is documented below.

<a name="success_list_struct"></a>
The `success_list` block supports:

* `id` - The UUID.

* `version` - The version number.

* `dict_id` - The dictionary ID.

* `dict_key` - The dictionary key.

* `dict_code` - The dictionary code.

* `dict_val` - The dictionary value.

* `dict_pkey` - The parent key of the dictionary.

* `dict_pcode` - The parent code of the dictionary.

* `create_time` - The creation time.

* `update_time` - The update time.

* `publish_time` - The publish time.

* `scope` - The scope to which the dictionary belongs.

* `description` - The description of the dictionary.

* `extension_field` - The extension fields.

* `project_id` - The project ID.

* `language` - The current language environment.

<a name="failed_list_struct"></a>
The `failed_list` block supports:

* `id` - The UUID.

* `version` - The version number.

* `dict_id` - The dictionary ID.

* `dict_key` - The dictionary key.

* `dict_code` - The dictionary code.

* `dict_val` - The dictionary value.

* `dict_pkey` - The parent key of the dictionary.

* `dict_pcode` - The parent code of the dictionary.

* `create_time` - The creation time.

* `update_time` - The update time.

* `publish_time` - The publish time.

* `scope` - The scope to which the dictionary belongs.

* `description` - The description of the dictionary.

* `extension_field` - The extension fields.

* `project_id` - The project ID.

* `language` - The current language environment.
