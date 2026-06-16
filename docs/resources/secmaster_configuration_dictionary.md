---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_configuration_dictionary"
description: |-
  Manages a configuration dictionary resource within HuaweiCloud.
---

# huaweicloud_secmaster_configuration_dictionary

Manages a configuration dictionary resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_secmaster_configuration_dictionary" "test" {
  dict_id      = "3027"
  dict_key     = "alert_comments"
  dict_code    = "Open"
  dict_val     = "Open"
  language     = "zh"
  version      = "1.0.0"
  dict_pkey    = "alert_comments_pkey"
  scope        = "ALERT"
  description  = "alert_comments_action"
  extend_field = {"test": "extend"}
  is_built_in  = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `dict_id` - (Required, String, NonUpdatable) Specifies the dictionary ID.

* `dict_key` - (Required, String, NonUpdatable) Specifies the dictionary key.

* `dict_code` - (Required, String) Specifies the dictionary code.

* `dict_val` - (Required, String) Specifies the dictionary value.

* `language` - (Required, String, NonUpdatable) Specifies the language environment.

* `version` - (Optional, String, NonUpdatable) Specifies the version number.

* `dict_pkey` - (Optional, String) Specifies the parent key of the dictionary.

* `dict_pcode` - (Optional, String) Specifies the parent code of the dictionary.

* `scope` - (Optional, String, NonUpdatable) Specifies the domain to which the dictionary belongs.

* `description` - (Optional, String) Specifies the description of the dictionary.

* `extend_field` - (Optional, Map) Specifies the extension field of the dictionary.

* `is_built_in` - (Optional, Bool, NonUpdatable) Specifies whether to create a built-in dictionary.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the UUID of the dictionary.

* `create_time` - The creation time of the dictionary.

* `update_time` - The latest update time of the dictionary.

* `publish_time` - The publish time of the dictionary.

## Import

The configuration dictionary can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_secmaster_configuration_dictionary.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include `is_built_in`. It is generally recommended running `terraform plan`
after importing a resource. You can then decide if changes should be applied to the resource,
or the resource definition should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_configuration_dictionary" "test" {
  ...

  lifecycle {
    ignore_changes = [
      is_built_in,
    ]
  }
}
```
