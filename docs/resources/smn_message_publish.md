---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_message_publish"
description: |-
  Manages a SMN message publishment resource within HuaweiCloud.
---

# huaweicloud_smn_message_publish

Manages a SMN message publishment resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

### Basic Example

```hcl
variable "topic_urn" {}
variable "subject" {}

resource "huaweicloud_smn_message_publish" "test" {
  topic_urn = var.topic_urn
  subject   = var.subject
  message   = "test"
}
```

### Message Structure Example

```hcl
variable "topic_urn" {}
variable "subject" {}

resource "huaweicloud_smn_message_publish" "test" {
  topic_urn = var.topic_urn
  subject   = var.subject

  message_structure = jsonencode({
    default       = "send default msg"
    sms           = "send msg by sms protocol"
    email         = "send msg by email protocol"
    http          = "send msg by http protocol"
    functiongraph = "send msg by functiongraph protocol"
    https         = "send msg by https protocol"
  })
}
```

### Message Template Example

```hcl
variable "topic_urn" {}
variable "subject" {}

resource "huaweicloud_smn_message_publish" "test" {
  topic_urn             = var.topic_urn
  subject               = var.subject
  message_template_name = "test"

  tags = {
    key = "value"
  }
}
```

### Add Message Attributes type Example

```hcl
variable "topic_urn" {}
variable "subject" {}

resource "huaweicloud_smn_message_publish" "test" {
  topic_urn = var.topic_urn
  subject   = var.subject
  message   = "test"

  # STRING type format
  message_attributes {
    name  = "test"
    type  = "STRING"
    value = "aaa"
  }

  # STRING_ARRAY type format
  message_attributes {
    name   = "aaa"
    type   = "STRING_ARRAY"
    values = ["aaa", "aaaa"]
  }

  # PROTOCOL type format
  message_attributes {
    name   = "smn_protocol"
    type   = "PROTOCOL"
    values = ["https", "http", "email", "sms"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `topic_urn` - (Required, String, NonUpdatable) Specifies the resource identifier of a topic.

* `subject` - (Optional, String, NonUpdatable) Specifies the message title.

* `message` - (Optional, String, NonUpdatable) Specifies the message content.

* `message_attributes` - (Optional, List, NonUpdatable) Specifies the message filter policies of a subscriber.
  The [message_attributes](#smn_message_attributes) structure is documented below.

* `message_structure` - (Optional, String, NonUpdatable) Specifies the message structure.

* `message_template_name` - (Optional, String, NonUpdatable) Specifies the message template name.

* `time_to_live` - (Optional, String, NonUpdatable) Specifies the maximum retention time of the message within the SMN system.
  After this retention time, the system will no longer send the message. The unit is second, and the default value
  of the variable is **3600** second. The value is a positive integer and less than or equal to **3600*24**.

* `tags` - (Optional, Map, NonUpdatable) Specifies a dictionary consisting of tag and parameters to replace the tag.
  The value corresponding to the label in the message template. Message publishing using message template mode must
  carry this parameter. The key in the dictionary is the parameter name in the message template, which should not
  exceed **21** characters. The value in the dictionary is the value after replacing the parameters in the message
  template, which does not exceed 1KB.

<a name="smn_message_attributes"></a>
The `message_attributes` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the property name.

* `type` - (Required, String, NonUpdatable) Specifies the property type.
  The value can be **STRING**, **STRING_ARRAY** or **PROTOCOL**.

* `value` - (Optional, String, NonUpdatable) Specifies the property value.
  This parameter is valid only when the `type` set to **STRING**. The attribute value can only contain Chinese
  and English, numbers, and underscores, and the length is **1** to **32** characters.

* `values` - (Optional, List, NonUpdatable) Specifies the property values.
  This parameter is valid when the `type` set to **STRING_ARRAY** or **PROTOCOL**.
  + When the `type` is **STRING_ARRAY**, the `values` is a string array, the array length is
    **1** to **10**, the element content in the array cannot be repeated, each string in the array can only contain
    Chinese and English, numbers, and underscores, and length is **1** to **32** characters.
  + When the `type` is **PROTOCOL**, the `values` is a string array of supported protocol types.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
