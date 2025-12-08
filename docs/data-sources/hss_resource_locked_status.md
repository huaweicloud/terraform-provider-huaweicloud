---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_resource_locked_status"
description: |-
  Use this data source to get the locked status of HSS resource within HuaweiCloud.
---

# huaweicloud_hss_resource_locked_status

Use this data source to get the locked status of HSS resource within HuaweiCloud.

## Example Usage

```hcl
variable "resource_id" {}

data "huaweicloud_hss_resource_locked_status" "test" {
  resource_id = var.resource_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `resource_id` - (Required, String) Specifies the quota ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `locked_status` - The locked status of the resource.  
  The valid values are as follows:
  + **true**: The resource is locked and cannot be converted from pay-as-you-go to yearly/monthly billing.
  + **false**: The resource is not locked and can be converted from pay-as-you-go to yearly/monthly billing.
