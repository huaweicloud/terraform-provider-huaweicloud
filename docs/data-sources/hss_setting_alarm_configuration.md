---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_alarm_configuration"
description: |-
  Use this data source to query the alarm configuration.
---

# huaweicloud_hss_setting_alarm_configuration

Use this data source to query the alarm configuration.

## Example Usage

```hcl
data "huaweicloud_hss_setting_alarm_configuration" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarm_type` - The alarm type.

* `display_name` - The display name.

* `topic_urn` - The SMN URN.

* `daily_alarm` - Whether to enable the daily alarm.

* `realtime_alarm` - Whether to enable the real-time alarm function.

* `alarm_level` - The alarm severity.

* `ignored_event_class_list` - The ignored events.
