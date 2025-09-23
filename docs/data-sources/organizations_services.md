---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_services"
description: |-
  Use this data source to get the list of cloud services that can be integrated with Organizations.
---

# huaweicloud_organizations_services

Use this data source to get the list of cloud services that can be integrated with Organizations.

## Example Usage

```hcl
data "huaweicloud_organizations_services" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `services` - Indicates the list of cloud services that can be integrated with Organizations.
