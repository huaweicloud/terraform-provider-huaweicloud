---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_cloud_service_authorizations"
description: |-
  Use this data source to get the list of AOM cloud service authorizations.
---

# huaweicloud_aom_cloud_service_authorizations

Use this data source to get the list of AOM cloud service authorizations.

## Example Usage

```hcl
data "huaweicloud_aom_cloud_service_authorizations" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String)  Specifies the region where the cloud service authorizations are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `authorizations` - The list of cloud service authorizations.
  The [authorizations](#aom_cloud_service_authorizations) structure is documented below.

<a name="aom_cloud_service_authorizations"></a>
The `authorizations` block supports:

* `service` - The authorization service name.

* `role_name` - The role names list.

* `status` - Whether the authorization is enabled.

* `need_optimized` - Whether the authorization needs optimization.
