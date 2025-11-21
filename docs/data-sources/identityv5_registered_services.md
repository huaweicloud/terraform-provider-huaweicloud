---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_registered_services"
description: |-
  Use this data source to get the list of registered services for auth schema in the Identity and Access Management V5 service.
---

# huaweicloud_identityv5_registered_services

Use this data source to get the list of registered services for auth schema in the Identity and Access Management V5 service.

## Example Usage

```hcl
data "huaweicloud_identityv5_registered_services" "services" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `service_codes` - The list of service codes.
