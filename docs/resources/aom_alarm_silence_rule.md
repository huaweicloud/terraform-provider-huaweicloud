---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_silence_rule"
description: ""
---

# huaweicloud_aom_alarm_silence_rule

Manages an AOM alarm silence rule resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_aom_alarm_silence_rule" "test" {
  name        = "test_rule"
  description = "terraform test"
  time_zone   = "Asia/Shanghai"

  silence_time {
    type      = "WEEKLY"
    starts_at = 64800
    ends_at   = 86399
    scope     = [1, 2, 3, 4, 5]
  }

  silence_conditions {
    conditions {
      key     = "event_severity"
      operate = "EXIST"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the rule name.

  Changing this parameter will create a new resource.

* `time_zone` - (Required, String, ForceNew) Specifies the time zone, e.g. **Asia/Shanghai**.

  Changing this parameter will create a new resource.

* `silence_time` - (Required, List) Specifies the silence time of the rule.
  The [silence_time](#silence_time) structure is documented below.

* `silence_conditions` - (Required, List) Specifies the silence conditions of the rule.
  Different silence conditions are parallel. A maximum of 10 silence conditions are allowed.
  The [silence_conditions](#silence_conditions) structure is documented below.

* `description` - (Optional, String) Specifies the description.

<a name="silence_time"></a>
The `silence_time` block supports:

* `type` - (Required, String) Specifies the effective time type of the silence rule.
  The value can be: **FIXED**, **DAILY**, **WEEKLY** and **MONTHLY**.

* `starts_at` - (Required, Int) Specifies the start time of the silence rule.
  When the `type` is **FIXED**, the value is a time stamp, e.g. **1684466549755**,
  which indicates **2023-05-19 11:22:29.755**. When the `type` is **DAILY**, **WEEKLY**
  or **MONTHLY**, the value range is `0` to `86,399`, which indicates **00:00:00** to **23:59:59**.

* `ends_at` - (Optional, Int) Specifies the end time of the silence rule.
  When the `type` is **FIXED**, the value is a time stamp, e.g. **1684466549755**,
  which indicates **2023-05-19 11:22:29.755**. When the `type` is **DAILY**, **WEEKLY**
  or **MONTHLY**, the value range is `0` to `86,399`, which indicates **00:00:00** to **23:59:59**.

* `scope` - (Optional, List) Specifies the silence time of the rule.
  It's required when the type is **WEEKLY** or **MONTHLY**.

<a name="silence_conditions"></a>
  The `silence_conditions` block supports:

* `conditions` - (Required, List) Specifies the serial conditions.
  A maximum of 10 conditions are allowed.
  The [conditions](#conditions) structure is documented below.

<a name="conditions"></a>
The `conditions` block supports:

* `key` - (Required, String) Specifies the key of the match condition.

* `operate` - (Required, String) Specifies the operate of the match condition.
  The value can be: **EQUALS**, **REGEX** and **EXIST**.

* `value` - (Optional, List) Specifies the value list of the match condition.
  A maximum of 5 values are allowed. This should be empty when the value of operate is *EXIST**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the rule name.

* `created_at` - The creation time.

* `updated_at` - The last update time.

## Import

The application operations management can be imported using the `id` (name), e.g.

```bash
$ terraform import huaweicloud_aom_alarm_silence_rule.test test_rule
```
