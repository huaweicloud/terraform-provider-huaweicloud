---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_providers"
description: |-
  Use this data source to get the supported cloud services and regions within HuaweiCloud Resource Center.
---

# huaweicloud_rc_providers

Use this data source to get the supported cloud services and regions within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_providers" "test" {}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `resource_providers` - The cloud service detail list.

  The [resource_providers](#resource_providers_struct) structure is documented below.

<a name="resource_providers_struct"></a>
The `resource_providers` block supports:

* `provider` - The cloud service name.

* `display_name` - The cloud service display name, language can be set via X-Language in the request header.

* `category_display_name` - The cloud service category display name, language can be set via X-Language in the request
  header.

* `resource_types` - The resource type list.

  The [resource_types](#resource_types_struct) structure is documented below.

<a name="resource_types_struct"></a>
The `resource_types` block supports:

* `name` - The resource type name.

* `display_name` - The resource type display name, language can be set via X-Language in the request.

* `global` - Whether this is a global resource type.

* `regions` - The supported region list.

* `console_endpoint_id` - The console endpoint ID.

* `console_list_url` - The console list page URL.

* `console_detail_url` - The console detail page URL.
