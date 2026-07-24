---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_alarm_infos"
description: |-
  Use this data source to get the list of DSC device alarm infos within HuaweiCloud.
---

# huaweicloud_dsc_alarm_infos

Use this data source to get the list of DSC device alarm infos within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_alarm_infos" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarm_infos` - The list of the alarm infos.

  The [alarm_infos](#alarm_infos_struct) structure is documented below.

<a name="alarm_infos_struct"></a>
The `alarm_infos` block supports:

* `count` - The cumulative number of times the alarm occurred.

* `create_time` - The timestamp when the alarm was generated.

* `description` - The detail description of the alarm.

* `device_ip` - The IP address of the device to which the alarm belongs.

* `id` - The unique ID of the alarm.

* `module` - The module to which the alarm belongs.

* `name` - The name of the alarm.

* `severity` - The severity level of the alarm.
  The valid values are as follows:
  + **1**: CRITICAL.
  + **2**: MAJOR.
  + **3**: MINOR.
  + **4**: WARNING.
  + **5**: INDETERMINATE.
  + **6**: CLEARED.

* `status` - The current handling status of the alarm.

* `type` - The device type code.
