---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_pay_per_scan_switch_status"
description: |-
  Manages an HSS update antivirus pay-per-scan switch status resource within HuaweiCloud.
---

# huaweicloud_hss_antivirus_pay_per_scan_switch_status

Manages an HSS update antivirus pay-per-scan switch status resource within HuaweiCloud.

-> This resource is a one-time action resource using to update HSS antivirus pay-per-scan switch status. Deleting this
  resource will not clear the corresponding request record, but will only remove the resource information from the tf
  state file.

## Example Usage

```hcl
resource "huaweicloud_hss_antivirus_pay_per_scan_switch_status" "test" {
  enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `enabled` - (Required, Bool, NonUpdatable) Specifies whether to enable virus scanning and billing per use function.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset  under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
