---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_stream_index_configuration"
description: |-
  Manages an LTS log stream index configuration resource within HuaweiCloud.
---

# huaweicloud_lts_stream_index_configuration

Manages an LTS log stream index configuration resource within HuaweiCloud.

-> 1. Only one resource can be created in one log stream.
   <br/>2. Deleting this resource will not initialize the currently configured index, but will only remove
   the resource information from the tfstate file.

## Example Usage

```hcl
variable "group_id" {}
variable "stream_id" {}
variable "full_text_index_tokenizer" {}
variable "index_fields" {
  type = list(object({
    field_name      = string
    field_type      = string
    tokenizer       = optional(string)
    case_sensitive  = optional(bool)
    include_chinese = optional(bool)
    quick_analysis  = optional(bool)
    ascii           = optional(list(string))

    lts_sub_fields_info_list = optional(list(object({
      field_name     = string
      field_type     = string
      quick_analysis = optional(bool)
    })))
  }))

  default = []
}

resource "huaweicloud_lts_stream_index_configuration" "test" {
  group_id  = var.group_id
  stream_id = var.stream_id

  full_text_index {
    tokenizer       = var.full_text_index_tokenizer
    enable          = true
    include_chinese = true
  }

  dynamic "fields" {
    for_each = var.index_fields

    content {
      field_name      = fields.value["field_name"]
      field_type      = fields.value["field_type"]
      tokenizer       = fields.value["tokenizer"]
      case_sensitive  = fields.value["case_sensitive"]
      include_chinese = fields.value["include_chinese"]
      quick_analysis  = fields.value["quick_analysis"]
      ascii           = fields.value["ascii"]

      dynamic "lts_sub_fields_info_list" {
        for_each = fields.value["lts_sub_fields_info_list"] != null ? fields.value["lts_sub_fields_info_list"] : []

        content {
          field_name     = lts_sub_fields_info_list.value["field_name"]
          field_type     = lts_sub_fields_info_list.value["field_type"]
          quick_analysis = lts_sub_fields_info_list.value["quick_analysis"]
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the log group to which the index configuration belongs.

* `stream_id` - (Required, String, NonUpdatable) Specifies the ID of the log stream to which the index configuration belongs.

* `full_text_index` - (Optional, List) Specifies the full-text index configuration.  
  The [full_text_index](#stream_index_config_full_text_index) structure is documented below.

* `fields` - (Optional, List) Specifies the list of the index fields.  
  The [fields](#stream_index_config_fields) structure is documented below.

<a name="stream_index_config_full_text_index"></a>
The `full_text_index` block supports:

* `tokenizer` - (Optional, String) Specifies the custom delimiter.

* `enable` - (Optional, Bool) Specifies whether to enable the full-text index.  
  Defaults to **true**.

* `case_sensitive` - (Optional, Bool) Specifies whether letters are case sensitive.  
  Defaults to **false**.

* `include_chinese` - (Optional, Bool) Specifies whether to include Chinese.  
  Defaults to **true**.

* `ascii` - (Optional, List) Specifies the list of the ASCII delimiters.  
  For more ASCII delimiters, please refer to the [document](https://support.huaweicloud.com/intl/en-us/usermanual-lts/lts_05_0008.html#lts_05_0008__section15661144724914).

<a name="stream_index_config_fields"></a>
The `fields` block supports:

* `field_name` - (Required, String) Specifies the name of the field.  
  The field name only letters, digits, hyphens (-), underscores (_) and dots (.) are allowed.  
  The name cannot start with a dot and end with a double underscores (__) or a dot.

* `field_type` - (Required, String) Specifies the type of the field.  
  The valid values are as follows:
  + **string**
  + **json**
  + **long**
  + **float**

* `tokenizer` - (Optional, String) Specifies the custom delimiter.  
  The parameter is available only when the `fields.field_type` parameter set to **string** or **json**.

* `case_sensitive` - (Optional, Bool) Specifies whether letters are case sensitive.  
  Defaults to **false**.  
  The parameter is available only when the `fields.field_type` parameter set to **string** or **json**.

* `include_chinese` - (Optional, Bool) Specifies whether to include Chinese.  
  Defaults to **false**.  
  The parameter is available only when the `fields.field_type` parameter set to **string** or **json**.

* `quick_analysis` - (Optional, Bool) Specifies whether to enable quick analysis.  
  Defaults to **false**.

* `ascii` - (Optional, List) Specifies the list of the ASCII delimiters.  
  The parameter is available only when the `fields.field_type` parameter set to **string** or **json**.  
  For more ASCII delimiters, please refer to the [document](https://support.huaweicloud.com/intl/en-us/usermanual-lts/lts_05_0008.html#lts_05_0008__section15661144724914).
  
* `field_analysis_alias` - (Optional, String) Specifies the alias name of the field.  
  Currently, only available in `cn-north-9`, `ap-southeast-1`, `ap-southeast-3` and `cn-east-3` regions.

* `lts_sub_fields_info_list` - (Optional, List) Specifies the list of of the JSON fields.  
  The [lts_sub_fields_info_list](#stream_index_config_fields_lts_sub_fields_info_list) structure is documented below.  
  The parameter is available only when the `fields.field_type` parameter set to **json**.

<a name="stream_index_config_fields_lts_sub_fields_info_list"></a>
The `lts_sub_fields_info_list` block supports:

* `field_name` - (Required, String) Specifies the name of the field.  
  The field name only letters, digits, hyphens (-), underscores (_) and dots (.) are allowed.  
  The name cannot start with a dot and end with a double underscores (__) or a dot.

* `field_type` - (Required, String) Specifies the type of the field.  
  The valid values are as follows:
  + **string**
  + **long**
  + **float**

* `quick_analysis` - (Optional, Bool) Specifies whether to enable quick analysis.  
  Defaults to **false**.

* `field_analysis_alias` - (Optional, String) Specifies the alias name of the field.  
  Currently, only available in `cn-north-9`, `ap-southeast-1`, `ap-southeast-3` and `cn-east-3` regions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The resource can be imported using the `group_id` and `stream_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_lts_stream_index_configuration.test <group_id>/<stream_id>
```
