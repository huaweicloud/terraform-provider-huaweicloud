---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_subscription"
description: ""
---

# huaweicloud_smn_subscription

Manages an SMN subscription resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_smn_topic" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "huaweicloud_smn_subscription" "subscription_1" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "O&M"
}

resource "huaweicloud_smn_subscription" "subscription_2" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "13600000000"
  protocol  = "sms"
  remark    = "O&M"
}

resource "huaweicloud_smn_subscription" "subscription_3" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "https://example.com/notification"
  protocol  = "https"
  remark    = "API webhook"
  
  extension {
    header = {
      "X-Custom-Test" = "test"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SMN subscription resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `topic_urn` - (Required, String, ForceNew) Specifies the resource identifier of a topic, which is unique.
  Changing this parameter will create a new resource.

* `protocol` - (Required, String, ForceNew) Specifies the protocol of the message endpoint. Currently, **email**,
  **sms**, **http**, **https**, **functionstage**, **functiongraph**, **callnotify**, **wechat**, **dingding**,
  **feishu** and **welink** are supported. Changing this parameter will create a new resource.

* `endpoint` - (Required, String, ForceNew) Message endpoint. Changing this parameter will create a new resource.
  + **For an HTTP subscription**, the endpoint starts with `http://`.
  + **For an HTTPS subscription**, the endpoint starts with `https://`.
  + **For an email subscription**, the endpoint is an mail address.
  + **For an SMS message subscription**, the endpoint is a phone number,
    the format is \[+\]\[country code\]\[phone number\], e.g. +86185xxxx0000.
  + **For a functionstage subscription**, the endpoint is a function urn.
  + **For a functiongraph subscription**, the endpoint is a workflow ID.
  + **For a callnotify subscription**, the endpoint is a phone number,
    the format is \[+\]\[country code\]\[phone number\], e.g. +86185xxxx0000.
  + **For a dingding subscription**, the endpoint is an IP address of a DingTalk group chatbot.
  + **For a wechat subscription**, the endpoint is an IP address of a WeChat group chatbot.
  + **For a feishu subscription**, the endpoint is an IP address of a Lark group chatbot.
  + **For a welink subscription**, the endpoint is a WeLink group account.

* `remark` - (Optional, String) Remark information. The remarks must be a UTF-8-coded character string
  containing 128 bytes.

* `extension` - (Optional, List, ForceNew) Specifies the extension configurations.
  The [extension](#extension) structure is documented below.
  Changing this parameter will create a new resource.

<a name="extension"></a>
The `extension` block supports:

* `client_id` - (Optional, String, ForceNew) Specifies the client ID. This field is the tenant ID field in
  the WeLink subscription and is obtained by the tenant from WeLink. This field is mandatory when `protocol`
  is set to **welink**. Changing this parameter will create a new resource.

* `client_secret` - (Optional, String, ForceNew) Specifies the client secret. This field is the client secret
  field obtained by the tenant from WeLink. This field is mandatory when `protocol` is set to **welink**.
  Changing this parameter will create a new resource.

* `keyword` - (Optional, String, ForceNew) Specifies the keyword. When `protocol` is set to **feishu**,
  either `keyword` or `sign_secret` must be specified. When you use `keywords` to configure a security policy
  for the Lark or DingTalk chatbot on SMN, the keywords must have one of the keywords configured on the Lark
  or DingTalk client. Changing this parameter will create a new resource.

* `sign_secret` - (Optional, String, ForceNew) Specifies the key including signature. When `protocol` is set
  to **feishu** or **dingding**, this field or `keyword` must be specified. The key configurations must be
  the same as those on the Lark or DingTalk client. For example, if only key is configured on the Lark client,
  enter the key field obtained from the Lark client. If only keyword is configured on the Lark client, skip this field.
  Changing this parameter will create a new resource.

* `header` - (Optional, Map, ForceNew) Specifies the HTTP/HTTPS headers to be added to the requests when the
  message is delivered via HTTP/HTTPS. This field is used when `protocol` is set to **http** or **https**.
  The following requirements apply to the header keys and values:
  + Header keys must:
    - Contain only letters, numbers, and hyphens (`[A-Za-z0-9-]`)
    - Not end with a hyphen
    - Not contain consecutive hyphens
    - Start with "x-" (e.g., "x-abc-cba", "x-abc")
    - Not start with "x-smn"
    - Be case-insensitive (e.g., "X-Custom" and "x-custom" are considered the same)
    - Not be duplicated
  + Maximum of 10 key-value pairs allowed
  + Total length of all keys and values combined must not exceed 1024 characters
  + Values must only contain ASCII characters (no Chinese or other Unicode characters, spaces are allowed)

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the subscription urn.

* `subscription_urn` - Resource identifier of a subscription, which is unique.

* `owner` - Project ID of the topic creator.

* `status` - Subscription status.
  + **0**: indicates that the subscription is not confirmed.
  + **1**: indicates that the subscription is confirmed.
  + **3**: indicates that the subscription is canceled.

* `filter_policies` - The message filter policies of a subscriber.
  The [filter_policies](#smn_subscription_filter_policies_attr) structure is documented below.

<a name="smn_subscription_filter_policies_attr"></a>
The `filter_policies` block supports:

* `name` - The filter policy name.

* `string_equals` - The string array for exact match.

## Import

SMN subscription can be imported using the `id` (subscription urn), e.g.

```bash
$ terraform import huaweicloud_smn_subscription.subscription_1 urn:smn:cn-north-4:0970dd7a1300f5672ff2c003c60ae115:topic_1:a2aa5a1f66df494184f4e108398de1a6
```
