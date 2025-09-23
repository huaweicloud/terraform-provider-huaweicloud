---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_notification_template"
description: |-
  Manages an LTS notification template resource within HuaweiCloud.  
---

# huaweicloud_lts_notification_template

Manages an LTS notification template resource within HuaweiCloud.  

## Example Usage

```hcl
resource "huaweicloud_lts_notification_template" "test" {
  name        = "keywords_alarm_demo"
  description = "keywords alarm demo"
  source      = "LTS"
  locale      = "en-us"

  templates {
    sub_type = "sms"
    content  = <<EOF
Account:$${domain_name};
Alarm Rules:<a href="$event.annotations.alarm_rule_url">$${event_name}</a>;
Alarm Status:$event.annotations.alarm_status;
Severity:<span style="color: red">$${event_severity}</span>;
Occurred:$${starts_at};
Type:Keywords;
Condition Expression:$event.annotations.condition_expression;
Current Value:$event.annotations.current_value;
Frequency:$event.annotations.frequency;
Log Group/Stream Name:$event.annotations.results[0].resource_id;
Query Time:$event.annotations.results[0].time;
Query URL:<a href="$event.annotations.results[0].url">details</a>;
EOF
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The name of the notification template.  
  The value can contain 1 to 100 characters, only chinese characters, digits, letters, and underscores (\_).  
  and can not start or end with an underscore (\_).
  Changing this parameter will create a new resource.

* `source` - (Required, String) The source of the notification template.  
  Currently, this parameter is mandatory to **LTS**..

* `locale` - (Required, String) Language.  
  Currently, the value can be **zh-cn** or **en-us**.

* `templates` - (Required, List) The list of notification template body.  
  The [templates](#NotificationTemplate_SubTemplate) structure is documented below.

* `description` - (Optional, String) The description of the notification template.  
  The value can contain 1 to 1024 characters, only chinese characters, digits, letters, and underscores (\_),  
  and can not start or end with an underscore (\_).  

<a name="NotificationTemplate_SubTemplate"></a>
The `templates` block supports:

* `sub_type` - (Required, String) The type of the sub-template.  
  Only the following five types are supported: **sms**, **dingding**, **wechat**, **webhook**, **email**.

* `content` - (Required, String) The content of the sub-template..  
  In the sub-template body, only the following variables are supported for the variables following the **$** symbol.  
  The supported variables vary according to the alarm type (keyword alarm and SQL alarm).  

    + Common variables:  
      * Alarm severity: **${event_severity}**.
      * Occurrence time: **${starts_at}**.
      * Alarm source: **$event.metadata.resource_provider**.
      * Resource type: **$event.metadata.resource_type**.
      * Resource ID: **${resources}**.
      * Expression: **$event.annotations.condition_expression**.
      * current value: **$event.annotations.current_value**.
      * Statistical period: **$event.annotations.frequency**.

    + Keywords alarm specific variable:  
      * query time: **$event.annotations.results[0].time**.
      * Run the **$event.annotations.results[0].raw_results** command to query LTSs.

    + SQL alarm specific variable:  
      * LTS group/stream name: **$event.annotations.results[0].resource_id**.
      * Query statement: **$event.annotations.results[0].sql**.
      * Query time: **$event.annotations.results[0].time**.
      * Query URL: **$event.annotations.results[0].url**.
      * Run the **$event.annotations.results[0].raw_results** command to query LTSs.

  -> semicolon(;) after variable is an English symbol and must be added. Otherwise, the template will fail to be replaced.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals `name`.

## Import

The LTS notification template can be imported using the `id` which equals `name`, e.g.

```bash
$ terraform import huaweicloud_lts_notification_template.test <id>
```
