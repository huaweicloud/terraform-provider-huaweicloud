---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_trusted_services"
description: |-
  Use this data source to get the list of the trusted services that are integrated with Organizations
---

# huaweicloud_organizations_trusted_services

Use this data source to get the list of the trusted services that are integrated with Organizations

## Example Usage

```hcl
data "huaweicloud_organizations_trusted_services" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `trusted_services` - Indicates the list of service principals for the services that are integrated with Organizations.

  The [trusted_services](#trusted_services_struct) structure is documented below.

<a name="trusted_services_struct"></a>
The `trusted_services` block supports:

* `service_principal` - Indicates the name of a trusted service.

* `enabled_at` - Indicates the date when the trusted service was integrated with Organizations
