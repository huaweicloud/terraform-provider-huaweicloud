---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_switches_status"
description: |-
  Use this data source to query the configuration switch status.
---

# huaweicloud_hss_setting_switches_status

Use this data source to query the configuration switch status.

## Example Usage

```hcl
variable "code" {}

data "huaweicloud_hss_setting_switches_status" "test" {
  code = var.code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `code` - (Required, String) Specifies the configuration type.
  The valid values are as follows:
  + **malware_sample_high_detect**: Sensitive malware detection mode.
  + **image_pay_per_scan**: Pay-per-use image scan switch.
  + **image_popup**: Pay-per-use image scan pop-up window switch.
  + **image_free_to_pay_popup**: Switch of the pop-up window for transferring from free to paid image scan.
  + **display_unprotected_host**: Show unprotected servers.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `enabled` - Whether the switch enabled.
