---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_modify_alarm_notification"
description: |-
  Manages a resource to update WAF alarm notification within HuaweiCloud.
---

# huaweicloud_waf_modify_alarm_notification

Manages a resource to update WAF alarm notification within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

```hcl
variable "alert_id" {}
variable "name" {}
variable "topic_urn" {}
variable "notice_class" {}

resource "huaweicloud_waf_modify_alarm_notification" "test" {
  alert_id     = var.alert_id
  name         = var.name
  topic_urn    = var.topic_urn
  notice_class = var.notice_class
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `alert_id` - (Required, String, NonUpdatable) Specifies the ID of the alarm notification.

* `name` - (Required, String, NonUpdatable) Specifies the name of the alarm notification.

* `topic_urn` - (Required, String, NonUpdatable) Specifies the topic URN of the SMN.

* `notice_class` - (Required, String, NonUpdatable) Specifies the type of the alarm notification.
  The valid values are as follows:
  + **threat_alert_notice**: Indicates protection event.
  + **cert_alert_notice**: Indicates certificate expiration. Currently, the type not supports.

* `enabled` - (Optional, Bool, NonUpdatable) Specifies whether to enable the alarm notification.
  The value can be **true** or **false**.

* `sendfreq` - (Optional, Int, NonUpdatable) Specifies the time interval, in minutes.
  The valid values are `5`, `15`, `30`, `60`, `120`, `360`, `720` and `1,440`.

* `locale` - (Optional, String, NonUpdatable) Specifies the language.
  The value can be **zh-cn**(Chineses) or **en-us**(English).

* `times` - (Optional, Int, NonUpdatable) Specifies the threshold of attack.
  An alarm notification is sent when the number of attack reaches the threshold.

* `threat` - (Optional, List, NonUpdatable) Specifies the event type.
  The valid values are as follows:
  **anticrawler**, **robot**, **advanced_bot**, **webshell**, **illegal**, **lfi**, **rfi**, **custom_whiteblackip**,
  **cmdi**, **custom_custom**, **custom_idc_ip**, **xss**, **cc**, **llm_prompt_sensitive**, **antitamper**, **vuln**,
  **leakage**, **llm_prompt_injection**, **third_bot_river**, **antiscan_dir_traversal**, **antiscan_high_freq_scan**,
  **botm**, **sqli**, **custom_geoip** and **llm_response_sensitive**.

  If you want to set all value above, you can set this parameter to **all**.
  If you not specifies this parameter, default is empty.

* `nearly_expired_time` - (Optional, String, NonUpdatable) Specifies the certificate expiration notification
  start time.
  This parameter is mandatory when the `notice_class` is set to **cert_alert_notice**.

* `is_all_enterprise_project` - (Optional, Bool, NonUpdatable) Specifies whether all enterprise projects
  are involved.
  The value can be **true** or **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
