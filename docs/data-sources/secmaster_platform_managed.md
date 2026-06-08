---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_platform_managed"
description: |-
  Use this data source to get the platform managed information.
---

# huaweicloud_secmaster_platform_managed

Use this data source to get the platform managed information.

## Example Usage

```hcl
data "huaweicloud_secmaster_platform_managed" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `create_time` - The creation time.

* `dw_region` - The region.

* `platform_managed_domain_id` - The platform tenant ID.

* `publish_status` - The publish status.

* `tenant_managed_domain_id` - The tenant managed domain ID.

* `update_time` - The update time.

* `whitelist` - Whether the tenant is in the whitelist.
