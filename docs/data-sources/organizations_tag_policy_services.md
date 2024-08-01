---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_tag_policy_services"
description: |-
  Use this data source to get the services that support enforcement with tag policies.
---

# huaweicloud_organizations_tag_policy_services

Use this data source to get the services that support enforcement with tag policies.

## Example Usage

```hcl
data "huaweicloud_organizations_tag_policy_services" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `services` - Indicates the services that support enforcement with tag policies.

  The [services](#services_struct) structure is documented below.

<a name="services_struct"></a>
The `services` block supports:

* `service_name` - Indicates the service name of the service.

* `resource_types` - Indicates the resource types.

* `support_all` - Indicates whether resource_type support all services (wildcard *).
