---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_registered_services"
description: |-
  Use this data source to get the list of registered services within HuaweiCloud.
---

# huaweicloud_identityv5_registered_services

Use this data source to get the list of registered services within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identityv5_registered_services" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `service_codes` - The list of service codes.
