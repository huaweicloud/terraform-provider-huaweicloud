---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_message_template"
description: ""
---

# huaweicloud_smn_message_template

Manages a SMN message template resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "protocol" {}

resource "huaweicloud_smn_message_template" "test"{
  name     = var.name
  protocol = var.protocol
  content  = "this content contains {content1}, {content2}, {content3}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the message template name.
  A template name starts with a letter or digit, consists of 1 to 64 characters,
  and can contain only letters, digits,  hyphens (-), and underscores (_).

  Changing this parameter will create a new resource.

* `protocol` - (Required, String, ForceNew) Specifies the protocol supported by the template. Value options:
  + **default**: the default protocol
  + **email**: the email protocol
  + **sms**: the SMS protocol
  + **functionstage**: the FunctionGraph transport protocol
  + **dms**: the DMS transport protocol
  + **http**: the http protocol
  + **https**: the https protocol

  Changing this parameter will create a new resource.

* `content` - (Required, String) Specifies the template content, which supports plain text only.
  The template content cannot be left blank or larger than 256 KB.
  The fields within "{}" can be replaced based on the actual situation
  when you use the template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `tag_names` - Indicates the variable list. The variable name will be quoted in braces ({}) in the template.
  When you use a template to send messages, you can replace the variable with any content.

## Import

The SMN message template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_smn_message_template.test <message_template_id>
```
