---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_services"
description: |-
  Use this data source to query the services list within HuaweiCloud.
---

# huaweicloud_identity_services

Use this data source to query the services list within HuaweiCloud.

## Example Usage

### Query services by ID

```hcl
variable "service_id" {}

data "huaweicloud_identity_services" "test" {
  service_id = var.service_id
}
```

### Query all services

```hcl
data "huaweicloud_identity_services" "test" {}
```

## Argument Reference

* `service_id` - (Optional, String) Specifies the service id.

## Attribute Reference

* `services` - Indicates service information list
  The [services](#IdentityServices_Services) structure is documented below.

<a name="IdentityServices_Services"></a>
The `services` block contains:

* `id` - Indicates service ID

* `type` - Indicates the service type.

* `name` - Indicates the service name.

* `link` - Indicates the Resource link of service.

* `enabled` - Indicates the service available.
