---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_services"
description: |-
  Use this data source to get the list of RMS services.
---

# huaweicloud_rms_services

Use this data source to get the list of RMS services.

## Example Usage

```hcl
data "huaweicloud_rms_services" "test" {
  name = "ecs"
}
```

## Argument Reference

The following arguments are supported:

* `track` - (Optional, String) Specifies whether resources are collected by default.
  The value can be **tracked** and **untracked**

* `name` - (Optional, String) Specifies the service name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `services` - The service details list.

  The [services](#services_struct) structure is documented below.

<a name="services_struct"></a>
The `services` block supports:

* `name` - The service name.

* `display_name` - The display name of the cloud service.

* `category_display_name` - The display name of the service category.

* `resource_types` - The resource type list.

  The [resource_types](#services_resource_types_struct) structure is documented below.

<a name="services_resource_types_struct"></a>
The `resource_types` block supports:

* `name` - The resource type name.

* `display_name` - The display name of the resource type.

* `global` - Indicates whether a resource is a global resource.

* `regions` - The list of supported regions.

* `track` - Indicates whether resources are collected by default.
  
  The value can be:
  + **tracked** indicates that resources are collected by default;
  + **untracked** indicates that resources are not collected by default.

* `console_endpoint_id` - The endpoint ID of the console.

* `console_detail_url` - The URL of the resource details page.

* `console_list_url` - The URL of the resource list page.
