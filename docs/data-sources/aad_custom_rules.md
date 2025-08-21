---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_custom_rules"
description: |-
  Use this data source to get the list of Advanced Anti-DDos custom rules within HuaweiCloud.
---

# huaweicloud_aad_custom_rules

Use this data source to get the list of Advanced Anti-DDos custom rules within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "overseas_type" {}

data "huaweicloud_aad_custom_rules" "test" {
  domain_name   = var.domain_name
  overseas_type = var.overseas_type
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Specifies the domain name.

* `overseas_type` - (Required, Int) Specifies protection zone.  
  The valid values are as follows:
  + **0**: Mainland.
  + **1**: Overseas.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `items` - The custom rules detail.  
  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID.

* `name` - The name.

* `time` - The precise protection rule effective time.  
  The valid values are as follows:
  + **true**: custom effective time.
  + **false**: take effect immediately.

* `start` - The start timestamp (in seconds) when the precise protection rule takes effect.

* `terminal` - The end timestamp (in seconds) when the precise protection rule takes effect.

* `priority` - The priority of executing this rule. The smaller the value, the higher the priority.
  Value range `0` to `1,000`.

* `conditions` - The conditions.  
  The [conditions](#conditions_struct) structure is documented below. .

* `action` - The action.  
  The [action](#action_struct) structure is documented below. .

* `domain_name` - The domain name.

* `overseas_type` - The protection zone.

<a name="conditions_struct"></a>
The `conditions` block supports:

* `category` - The field type.  
  The valid values are as follows:
  + **url**: URL.
  + **ip**: IPv4.
  + **user-agent**: User Agent.
  + **method**: Method.
  + **referer**: Referer.
  + **params**: Params.
  + **cookie**: Cookie.
  + **header**: Header.
  + **request_line**: Request Line.
  + **request**: Request.

* `index` - The sub-field.  
  Please refer to the document link [reference](https://support.huaweicloud.com/api-aad/ListWafCustomRuleV2.html)
  for values.

* `logic_operation` - The condition list matching logic.  
  Please refer to the document link [reference](https://support.huaweicloud.com/api-aad/ListWafCustomRuleV2.html)
  for values.

* `contents` - The condition list logic matching content.  
  Please refer to the document link [reference](https://support.huaweicloud.com/api-aad/ListWafCustomRuleV2.html)
  for values.

<a name="action_struct"></a>
The `action` block supports:

* `category` - The protection action.  
  The valid values are as follows:
  + **block**: Intercept.
  + **pass**: Allow.
  + **log**: Log only.
