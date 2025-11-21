---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_core_account"
description: |-
  Use this data source to get core account in Resource Governance Center.
---

# `huaweicloud_rgc_core_account`

Use this data source to get core account in Resource Governance Center.

## Example Usage

```hcl
variable account_type {}

data "huaweicloud_rgc_core_account" "test" {
  account_type = var.account_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `account_type` - (Required, String) The type of the core account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `account_id` - The ID of the core account.

* `core_resource_mappings` - The core resource mappings in JSON format.
