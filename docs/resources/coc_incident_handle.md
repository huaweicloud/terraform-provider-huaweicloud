---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_incident_handle"
description: |-
  Manages a COC incident handle resource within HuaweiCloud.
---

# huaweicloud_coc_incident_handle

Manages a COC incident handle resource within HuaweiCloud.

~> Deleting incident handle resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "incident_num" {}
variable "user_id" {}

resource "huaweicloud_coc_incident_handle" "test" {
  incident_num = var.incident_num
  operator     = var.user_id
  operate_key  = "acceptedIncident1"
}
```

## Argument Reference

The following arguments are supported:

* `incident_num` - (Required, String, NonUpdatable) Specifies the incident number.

* `operator` - (Required, String, NonUpdatable) Specifies the user ID of operator.

* `operate_key` - (Required, String, NonUpdatable) Specifies the operation type.
  The following scenarios can be applied, and the default values are different in different scenarios:
  + **Accepting an incident ticket**: The value is fixed to `acceptedIncident1`.
  + **Submitting an incident ticket solution**: The value is `commitSolution1`.
  + **Verifying the incident handling result**: The value is `confirm`.

* `parameter` - (Optional, Map, NonUpdatable) Specifies the parameter.
  The following scenarios can be applied, and different scenarios contain different fields:
  + Accepting an incident ticket: Leave this parameter empty.
  + Submitting an incident ticket solution:
      - **mtm_type**: Incident type.
        For details, see [mtm_type](https://support.huaweicloud.com/intl/en-us/api-coc/coc_api_04_03_001_006_02.html).
      - **is_service_interrupt**: Whether the service is interrupted. The value can be `true` or `false`.
      - **cause**: Reason.
      - **solution**: Solution.
      - **start_time**: The fault occurrence timestamp. Required when `is_service_interrupt` is `true`.
      - **fault_recovery_time**: The fault recovery timestamp. Required when `is_service_interrupt` is `true`.
  + Verifying the incident handling result:
      - **virtual_confirm_result**: Verification result. The value can be `true` (default) or `false`.
      - **virtual_confirm_comment**: Verification description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `incident_num`.
