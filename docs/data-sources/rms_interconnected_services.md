---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_interconnected_services"
description: |-
  Use this data source to get the list of all cloud services and resources interconnected with Config.
---

# huaweicloud_rms_interconnected_services

Use this data source to get the list of all cloud services and resources interconnected with Config.

## Example Usage

```hcl
data "huaweicloud_rms_interconnected_services" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_providers` - Indicates the list of cloud service details.

  The [resource_providers](#resource_providers_struct) structure is documented below.

<a name="resource_providers_struct"></a>
The `resource_providers` block supports:

* `provider` - Indicates the service name.

* `display_name` - Indicates the display name of the cloud service.

* `category_display_name` - Indicates the display name of the service category.

* `resource_types` - Indicates the resource type list.
  The [resource_types](#resource_providers_resource_types_struct) structure is documented below.

<a name="resource_providers_resource_types_struct"></a>
The `resource_types` block supports:

* `name` - Indicates the resource type name.

* `display_name` - Indicates the display name of the resource type.

* `global` - Indicates whether a resource is a global resource.

* `regions` - Indicates the list of supported regions.

* `track` - Indicates whether resources are collected by default.
  The value can be:
  + **tracked**: indicates that resources are collected by default;
  + **untracked**: indicates that resources are not collected by default.

* `console_endpoint_id` - Indicates the endpoint ID of the console.

* `console_detail_url` - Indicates the URL of the resource details page.

* `console_list_url` - Indicates the URL of the resource list page.
