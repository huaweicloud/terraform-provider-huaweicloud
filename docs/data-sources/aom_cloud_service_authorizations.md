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

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `authorizations` - Indicates the authorizations list.
  The [authorizations](#attrblock--authorizations) structure is documented below.

<a name="attrblock--authorizations"></a>
The `authorizations` block supports:

* `service` - Indicates the authorization service.

* `role_name` - Indicates the role names list.

* `status` - Indicates the authorization status.
