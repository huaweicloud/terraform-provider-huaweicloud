---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_topic_attributes"
description: |-
  Manages the SMN topic attributes within HuaweiCloud.
---

# huaweicloud_smn_topic_attributes

Manages the SMN topic attributes within HuaweiCloud.

-> Deleting this resource will not reset the topic attributes, but will only remove the resource information from the
   tfstate file.

## Example Usage

### Attributes with access policy

```hcl
variable "smn_topic_urn" {}

resource "huaweicloud_smn_topic_attributes" "access_policy" {
  topic_urn = var.smn_topic_urn
  name      = "access_policy"
  value     = jsonencode({
    "Version": "2016-09-07",
    "Id": "__default_policy_ID",
    "Statement": [
      {
        "Sid": "__org_path_pub_0",
        "Effect": "Allow",
        "Principal": {
          "OrgPath": [
            "o-xxx/r-xxx/ou-xxx"
          ]
        },
        "Action": [
          "SMN:Publish",
          "SMN:QueryTopicDetail"
        ],
        "Resource": var.smn_topic_urn
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the topic is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `topic_urn` - (Required, String, NonUpdatable) Specifies the topic URN. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the topic attribute name.  
  The valid values are as follows:
  + **access_policy**

* `value` - (Required, String) Specifies the topic attribute value, in JSON format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (format is **{topic_urn}/{name}**).

## Import

The SMN topic attributes can be imported using the `topic_urn` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_smn_topic_attributes.test {topic_urn}/{name}
```
