---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_alarm_config"
description: |-
  Manages a CFW alarm configuration resource within HuaweiCloud.
---

# huaweicloud_cfw_alarm_config

Manages a CFW alarm configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "alarm_type" {}
variable "alarm_time_period" {}
variable "severity" {}
variable "frequency_count" {}
variable "frequency_time" {}
variable "topic_urn" {}

resource "huaweicloud_cfw_alarm_config" "test" {
  fw_instance_id    = var.fw_instance_id
  alarm_type        = var.alarm_type
  alarm_time_period = var.alarm_time_period
  frequency_count   = var.frequency_count
  frequency_time    = var.frequency_time
  severity          = var.severity
  topic_urn         = var.topic_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall ID.

* `alarm_type` - (Required, Int, NonUpdatable) Specifies the alarm type.
  The valid values are as follows.
  + **0**: attack;
  + **1**: traffic threshold crossing;
  + **2**: EIP unprotected;
  + **3**: threat intelligence;

* `alarm_time_period` - (Required, Int) Specifies the alarm period.
  The valid values are as follows:
  + **0**: 8:00 to 22:00;
  + **1**: all day;

* `frequency_count` - (Required, Int) Specifies the alarm triggering frequency.
  + If `alarm_type` is **0** or **3**, the value of `frequency_count` must be between **1** and **2000**.
  + If `alarm_type` is **1** or **2**, the value of `frequency_count` should be **1**.

* `frequency_time` - (Required, Int) Specifies the alarm frequency time range.
  + If `alarm_type` is **0** or **3**, the value of `frequency_time` must be between **1** and **60**.
  + If `alarm_type` is **1** or **2**, the value of `frequency_time` should be **1**.

* `severity` - (Required, String) Specifies the alarm severity.
  + If `alarm_type` is **0** or **3**, the value of `severity` can be a combination of **CRITICAL**, **HIGH**,
  **MEDIUM**, and **LOW**, separated by commas.
  + If `alarm_type` is **1**, the value of `severity` can be **0** (70%), **1** (80%), or **2** (90%).
  + If `alarm_type` is **2**, the value of `severity` must be **3** (EIP).

* `topic_urn` - (Required, String) Specifies the alarm URN.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `language` - The language.

* `username` - The username.

## Import

The alarm configuration can be imported using `fw_instance_id`, `alarm_type`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cfw_alarm_config.test <fw_instance_id>/<alarm_type>
```
