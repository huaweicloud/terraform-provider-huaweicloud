---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_structuring_custom_configuration"
description: |-
  Manages an LTS structuring custom configuration resource within HuaweiCloud.
---

# huaweicloud_lts_structuring_custom_configuration

Manages an LTS structuring custom configuration resource within HuaweiCloud.

## Example Usage

### Creating with json structuring method

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  content       = "{'code':38,'user':{'name':'testdemo'}}"
  layers        = 3
  
  demo_fields {
    is_analysis = true
    field_name  = "code"
    content     = "38"
    type        = "long"
  }

  demo_fields {
    is_analysis = true
    field_name  = "user.name"
    content     = "testdemo"
    type        = "string"
  }

  tag_fields {
    is_analysis = true
    field_name  = "hostIP"
    content     = "192.168.2.134"
    type        = "string"
  }
}
```

### Creating with split structuring method

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  content       = "2023-09-09/18:50:51 Error"
  tokenizer     = " "
  
  demo_fields {
    is_analysis = true
    field_name  = "b1"
    content     = "2023-09-09/18:50:51"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "b2"
    content     = "Error"
    type        = "string"
  }
}
```

### Creating with nginx structuring method

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  content       = "39.149.31.187 - - [12/Mar/2020:12:24:02 +0800] \"GET / HTTP/1.1\" 304 "
  log_format    = "log_format  main   '$remote_addr - $remote_user [$time_local] \"$request\" '\n'$status ';"

  demo_fields {
    is_analysis = true
    field_name  = "remote_addr"
    content     = "39.149.31.187"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "remote_user"
    content     = "-"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "request_method"
    content     = "GET"
    type        = "string"
  }
}
```

### Creating with custom regex structuring method

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  content       = "2023-09-09/18:15:41 this log is Error NO 6323"
  regex_rules   = "^(?<a01>[^ ]+)(?:[^ ]* ){1}(?<a02>\\w+)(?:[^ ]* ){1}(?<a03>\\w+)(?:[^ ]* )"

  demo_fields {
    is_analysis = true
    field_name  = "a01"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "a02"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "a03"
    type        = "string"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `log_group_id` - (Required, String, ForceNew) Specifies the log group ID.
  Changing this parameter will create a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the log stream ID.
  Changing this parameter will create a new resource.

* `content` - (Required, String) Specifies a sample log event. LTS will create parsing rules based on the sample log.

* `demo_fields` - (Required, List) Specifies the list of example fields. The maximum length is `200`. The field sequence
  in `demo_fields` must be the same as that in `content`. The listed fields will be used as log extraction fields.
The [demo_fields](#StructCustomConfig_demo_fields) structure is documented below.

* `regex_rules` - (Optional, String) Specifies the regular expression. The maximum length is `5000` characters.
  When this field is specified, regular analysis will be used to parse the logs.

* `layers` - (Optional, Int) Specifies the maximum parsing layers. The maximum value is `3`.
  When this field is specified, the log body will be parsed in JSON format and split into key-value pairs.

* `tokenizer` - (Optional, String) Specifies the delimiter, such as spaces and colons.
  When this field is specified, the log body will be parsed by specifying separators.

* `log_format` - (Optional, String) Specifies the nginx configuration.
  When this field is specified, key-value pairs are extracted from Nginx log events.

-> The fields `regex_rules`, `layers`, `tokenizer` and `log_format` are mutually exclusive, and one of these fields
must be specified. Refer to [document](https://support.huaweicloud.com/intl/en-us/usermanual-lts/lts_0823.html) for more
information.

* `tag_fields` - (Optional, List) Specifies the tag field array. This field is only needed when tag fields are used for
  parsing.
The [tag_fields](#StructCustomConfig_tag_fields) structure is documented below.

<a name="StructCustomConfig_demo_fields"></a>
The `demo_fields` block supports:

* `field_name` - (Optional, String) Specifies the field name. The value ranges from `1` to `50`.

* `type` - (Optional, String) Specifies the field data type. Valid values are **string**, **long** and **float**.

* `is_analysis` - (Optional, Bool) Specifies whether quick analysis is enabled. Defaults to **false**.

* `content` - (Optional, String) Specifies the field content.

<a name="StructCustomConfig_tag_fields"></a>
The `tag_fields` block supports:

* `field_name` - (Required, String) Specifies the field name. The value ranges from `1` to `50`.

* `type` - (Required, String) Specifies the field data type. Valid values are **string**, **long** and **float**.

* `is_analysis` - (Optional, Bool) Specifies whether quick analysis is enabled. Defaults to **false**.

* `content` - (Optional, String) Specifies the field content.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The LTS structuring custom configuration can be imported using `log_group_id` and `log_stream_id`, separated by a slash,
e.g.

```bash
$ terraform import huaweicloud_lts_structuring_custom_configuration.test <log_group_id>/<log_stream_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `demo_fields`, `regex_rules`, `layers`,
`tokenizer`, `log_format`, `tag_fields`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      demo_fields,
      regex_rules,
      layers,
      tokenizer,
      log_format,
      tag_fields,
    ]
  }
}
```
