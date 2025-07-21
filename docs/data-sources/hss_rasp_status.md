---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_rasp_status"
description: |-
  Use this data source to get the application protection status.
---

# huaweicloud_hss_rasp_status

Use this data source to get the application protection status.

## Example Usage

```hcl
variable "host_id"

data "huaweicloud_hss_rasp_status" "test" {
  host_id               = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Required, String) Specifies the host ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `app_type` - (Optional, String) Specifies the application type.
  The value can be **java** (indicates Java application protection).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - The status of the application protection.
