---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_available_volumes"
description: |-
  Use this data source to get available volume list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_available_volumes

Use this data source to get available volume list of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_app_available_volumes" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `volume_types` - The list of available volume types.  
  The [volume_types](#workspace_app_volume_types) structure is documented below.

<a name="workspace_app_volume_types"></a>
The `volume_types` block supports:

* `resource_spec_code` - The resource specification code.

* `volume_type` - The volume type.

* `volume_product_type` - The volume product type.

* `resource_type` - The resource type.

* `cloud_service_type` - The cloud service type code.

* `name` - The volume type name in different languages.

* `volume_type_extra_specs` - The extra specifications of volume type.  
  The [volume_type_extra_specs](#workspace_app_volume_type_extra_specs) structure is documented below.

<a name="workspace_app_volume_type_extra_specs"></a>
The `volume_type_extra_specs` block supports:

* `availability_zone` - The availability zone for this volume type.

* `sold_out_availability_zone` - The sold out availability zone for this volume type.
