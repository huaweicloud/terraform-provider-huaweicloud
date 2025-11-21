---
subcategory: "RGC"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_account_compliance_status"
description: |-
  Use this data source to get the compliance status for an managed account in Resource Governance Center.
---

# huaweicloud_rgc_account_compliance_status

Use this data source to get the compliance status for an managed account in Resource Governance Center.

## Example Usage

```hcl
variable "managed_account_id" {}

data "huaweicloud_rgc_account_compliance_status" "test" {
  managed_account_id  = var.managed_account_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `control_id` - (Optional, String) The ID of the enabled control policy.

* `managed_account_id` - (Required, String) The ID of the managed account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `compliance_status` - The compliance status.
