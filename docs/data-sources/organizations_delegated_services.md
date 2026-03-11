---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_delegated_services"
description: |-
  Use this data source to get the service list for which the delegated administrator account within HuaweiCloud.
---

# huaweicloud_organizations_delegated_services

Use this data source to get the service list for which the delegated administrator account within HuaweiCloud.

-> This data source can only be called from the organization management account, or a member account that is assigned
   as a delegated administrator for a cloud service.

## Example Usage

```hcl
variable "account_id" {}

data "huaweicloud_organizations_delegated_services" "test" {
  account_id = var.account_id
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, String) Specifies the unique ID of an account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `delegated_services` - The list of the delegated services.  
  The [delegated_services](#Organizations_Delegated_Services) structure is documented below.

<a name="Organizations_Delegated_Services"></a>
The `delegated_services` block supports:

* `service_principal` - The name of the service principal.

* `delegation_enabled_at` - The date when the account became a delegated administrator for the service.
