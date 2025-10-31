---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_delegated_services"
description: |-
  Use this data source to get the list of delegated services.
---

# huaweicloud_organizations_delegated_services

Use this data source to list the services for which the specified account is a delegated administrator.
It can be called only from the organization's management account or from a member account that is a delegated
administrator for a cloud service.

## Example Usage

```hcl
variable "account_id" {}

data "huaweicloud_organizations_delegated_services" "test" {
  account_id = var.account_id
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, String) Unique ID of an account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `delegated_services` - List of the services for which the specified account is a delegated administrator.
  The [delegated_services](#Organizations_Delegated_Services) structure is documented below.

<a name="Organizations_Delegated_Services"></a>
The `delegated_services` block supports:

* `service_principal` - Name of the service principal.

* `delegation_enabled_at` - Date when the account became a delegated administrator for the service.
