---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_custom_event_report"
description: |-
  Manages a COC custom event report resource within HuaweiCloud.
---

# huaweicloud_coc_custom_event_report

Manages a COC custom event report resource within HuaweiCloud.

~> Deleting custom event report resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "integration_key" {}
variable "alarm_id" {}
variable "alarm_name" {}
variable "application_id" {}
variable "alarm_desc" {}

resource "huaweicloud_coc_custom_event_report" "test" {
  integration_key = var.integration_key
  alarm_id        = var.alarm_id
  alarm_name      = var.alarm_name
  alarm_level     = "Critical"
  time            = 1709118444540
  namespace       = "shanghai"
  application_id  = var.application_id
  alarm_desc      = var.alarm_desc
  alarm_source    = "coc"
  additional      = jsonencode({
    "key": "test"
  })
}
```

## Argument Reference

The following arguments are supported:

* `integration_key` - (Required, String, NonUpdatable) Specifies the integration key.

* `alarm_id` - (Required, String, NonUpdatable) Specifies the alarm ID.

* `alarm_name` - (Required, String, NonUpdatable) Specifies the alarm name.

* `alarm_level` - (Required, String, NonUpdatable) Specifies the alarm level.
  Values can be as follows:
  + **Critical**: Critical.
  + **Major**: Major.
  + **Minor**: Minor.
  + **Info**: Info.

* `time` - (Required, Int, NonUpdatable) Specifies the alarm occurrence time.

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace of the service.

* `application_id` - (Required, String, NonUpdatable) Specifies the application ID.

* `alarm_desc` - (Required, String, NonUpdatable) Specifies the alarm description.

* `alarm_source` - (Required, String, NonUpdatable) Specifies the alarm source.

* `region_id` - (Optional, String, NonUpdatable) Specifies the area where the alarm occurs.

* `resource_name` - (Optional, String, NonUpdatable) Specifies the resource name.

* `resource_id` - (Optional, String, NonUpdatable) Specifies the resource ID.

* `url` - (Optional, String, NonUpdatable) Specifies the original alert URL.

* `alarm_status` - (Optional, String, NonUpdatable) Specifies the alarm status.
  Values can be as follows:
  + **alarm**: Warning.
  + **ok**: Restored.

* `additional` - (Optional, String, NonUpdatable) Specifies the additional alarm information.

  -> The value of `additional` is a json string.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `integration_key`.
