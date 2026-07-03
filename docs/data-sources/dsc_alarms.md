---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_alarms"
description: |-
  Use this data source to get the list of DSC alarms within HuaweiCloud.
---

# huaweicloud_dsc_alarms

Use this data source to get the list of DSC alarms within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_alarms" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `alarm_level` - (Optional, String) Specifies the alarm level.

* `alarm_name` - (Optional, String) Specifies the alarm name.

* `alarm_status` - (Optional, String) Specifies the alarm status.

* `end_time` - (Optional, Int) Specifies the end time of the query (timestamp).

* `start_time` - (Optional, Int) Specifies the start time of the query (timestamp).

* `responsible_person` - (Optional, String) Specifies the responsible person for the alarm.

* `source_name` - (Optional, String) Specifies the alarm source name.

* `source_type` - (Optional, String) Specifies the alarm source type.

* `verification_status` - (Optional, String) Specifies the verification status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `alarms` - The alarm information list.

  The [alarms](#alarms_struct) structure is documented below.

<a name="alarms_struct"></a>
The `alarms` block supports:

* `id` - The alarm ID.

* `alarm_name` - The alarm name.

* `alarm_level` - The alarm level.

* `alarm_status` - The alarm status.

* `alarm_type` - The alarm type.

* `close_reason` - The reason for closing the alarm.

* `create_time` - The alarm creation time.

* `description` - The alarm description.

* `disposal_suggestion` - The disposal suggestion.

* `domain_id` - The domain ID.

* `occur_time` - The alarm occurrence time.

* `project_id` - The project ID.

* `source_module` - The alarm source module.

* `source_sub_type` - The alarm source sub type.

* `source_type` - The alarm source type.

* `verification_status` - The verification status.

* `affected_asset` - The affected assets.

* `responsible_person` - The responsible person information.

  The [responsible_person](#alarms_responsible_person_struct) structure is documented below.

* `source_instance` - The alarm source instance information.

  The [source_instance](#alarms_source_instance_struct) structure is documented below.

<a name="alarms_responsible_person_struct"></a>
The `responsible_person` block supports:

* `user_id` - The user ID.

* `user_name` - The user name.

<a name="alarms_source_instance_struct"></a>
The `source_instance` block supports:

* `instance_id` - The instance ID.

* `instance_name` - The instance name.
