---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_incident"
description: |-
  Manages a COC incident resource within HuaweiCloud.
---

# huaweicloud_coc_incident

Manages a COC incident resource within HuaweiCloud.

~> Deleting incident resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "incident_title" {}
variable "user_id" {}
variable "assignee" {}
variable "application_id" {}

resource "huaweicloud_coc_incident" "test" {
  incident_level        = "level_50"
  is_service_interrupt  = false
  incident_type         = "inc_type_p_security_issues"
  incident_title        = var.incident_title
  incident_source       = "incident_source_forwarding"
  creator               = var.user_id
  incident_assignee     = var.assignee
  current_cloud_service = [var.application_id]
  incident_description  = "description"
}
```

## Argument Reference

The following arguments are supported:

* `incident_level` - (Required, String, NonUpdatable) Specifies the incident level.
  For details, see [incident_level](https://support.huaweicloud.com/intl/en-us/api-coc/coc_api_04_03_001_006_01.html)

* `is_service_interrupt` - (Required, Bool, NonUpdatable) Specifies whether the service is interrupted.

* `incident_type` - (Required, String, NonUpdatable) Specifies the incident type.
  For details, see [incident_type](https://support.huaweicloud.com/intl/en-us/api-coc/coc_api_04_03_001_006_02.html)

* `incident_title` - (Required, String, NonUpdatable) Specifies the incident title.

* `incident_source` - (Required, String, NonUpdatable) Specifies the incident source.
  For details, see [incident_source](https://support.huaweicloud.com/intl/en-us/api-coc/coc_api_04_03_001_006_03.html)

* `creator` - (Required, String, NonUpdatable) Specifies the user ID of the creator.

* `regions` - (Optional, List, NonUpdatable) Specifies the region.
  This parameter is mandatory if a war room is automatically started.

* `enterprise_project` - (Optional, List, NonUpdatable) Specifies the enterprise project ID.

* `start_time` - (Optional, Int, NonUpdatable) Specifies the fault occurrence timestamp.

* `current_cloud_service` - (Optional, List, NonUpdatable) Specifies the application ID.

* `incident_ownership` - (Optional, String, NonUpdatable) Specifies the incident ownership.

* `incident_description` - (Optional, String, NonUpdatable) Specifies the incident description.

* `incident_assignee` - (Optional, List, NonUpdatable) Specifies the incident assignee.

* `assignee_scene` - (Optional, String, NonUpdatable) Specifies the scheduling scenario.

* `assignee_role` - (Optional, String, NonUpdatable) Specifies the scheduling role.
  
-> Only one of `incident_assignee` and `assignee_scene` can be specified, and `assignee_scene` and `assignee_role` must
be specified st the same time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `warroom_id` - Indicates the war room ID.

* `handle_time` - Indicates the timestamp of the last solution submission.

* `status` - Indicates the incident status.

* `create_time` - Indicates the incident create time.

* `enum_data_list` - Indicates the enumeration list.
  The [enum_data_list](#attrblock--enum_data_list) structure is documented below.

<a name="attrblock--enum_data_list"></a>
The `enum_data_list` block supports:

* `filed_key` - Indicates the filed key.

* `enum_key` - Indicates the enum key.

* `name_zh` - Indicates the Chinese name.

* `name_en` - Indicates the English name.

## Import

The COC incident can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_incident.test <id>
```
