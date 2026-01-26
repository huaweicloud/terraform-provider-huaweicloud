---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project_services"
description: |-
  Use this data source to get the services supported by EPS within HuaweiCloud.
---

# huaweicloud_enterprise_project_services

Use this data source to get the services supported by EPS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_enterprise_project_services" "test" {}
```

## Argument Reference

The following arguments are supported:

* `locale` - (Optional, String) Specifies the display language.

* `service` - (Optional, String) Specifies the cloud service name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `services` - Indicates the cloud services.
  The [services](#service_struct) structure is documented below.

<a name="service_struct"></a>
The `services` block supports:

* `service` - Indicates the cloud service name.

* `service_i18n_display_name` - Indicates the display name of the cloud service. You can set the language by setting the
  locale parameter.

* `resource_types` - Indicates the resource type list.
  The [resource_types](#resource_types_struct) structure is documented below.

<a name="resource_types_struct"></a>
The `resource_types` block supports:

* `resource_type` - Indicates the name of the resource type.

* `resource_type_i18n_display_name` - Indicates the display name of the resource type. You can set the language by
  setting the locale parameter.

* `global` - Whether the resource is a global resource.

* `regions` - Indicates regions supported.
