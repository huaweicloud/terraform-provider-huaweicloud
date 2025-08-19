---
subcategory: "rms"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_cloud_services"
description: |-
  Use this data source to get the list of Config cloud services.
---

# huaweicloud_rms_cloud_services

Use this data source to get the list of Config cloud services.

## Example Usage

```hcl
data "huaweicloud_rms_cloud_services" "test" {
  track = "tracked"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `track` - (Optional, String) Specifies whether resources are tracked by default. Value options: **tracked** and **untracked**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_providers` - The list of resource providers.

  The [resource_providers](#resource_providers_struct) structure is documented below.

<a name="resource_providers_struct"></a>
The `resource_providers` block supports:

* `provider` - The name of the cloud service.

* `display_name` - The display name of the cloud service, which can be set through **X-Language** in the request header.

* `category_display_name` - The category display name of the cloud service, which can be set through **X-Language** in
  the request header.

* `resource_types` - The list of resource types.

  The [resource_types](#resource_providers_resource_types_struct) structure is documented below.

<a name="resource_providers_resource_types_struct"></a>
The `resource_types` block supports:

* `name` - The name of the resource type.

* `display_name` - The display name of the resource type, which can be set through **X-Language** in the request header.

* `global` - Whether it is a global resource.

* `regions` - The list of supported regions.

* `console_endpoint_id` - The console endpoint ID.

* `console_list_url` - The console list page URL.

* `console_detail_url` - The console detail page URL.

* `track` - Whether resources are tracked by default.
