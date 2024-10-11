---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_applications"
description: |-
  Use this data source to query the list of applications within HuaweiCloud.
---

# huaweicloud_servicestagev3_applications

Use this data source to query the list of applications within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_servicestagev3_applications" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the applications are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `applications` - All application details.  
  The [applications](#servicestage_v3_applications) structure is documented below.

<a name="servicestage_v3_applications"></a>
The `applications` block supports:

* `id` - The application ID.

* `name` - The application name.

* `description` - The description of the application.

* `enterprise_project_id` - The ID of the enterprise project to which the application belongs.

* `tags` - The key/value pairs to associate with the application.

* `creator` - The creator name of the application.

* `created_at` - The creation time of the application, in RFC3339 format.

* `updated_at` - The latest update time of the application, in RFC3339 format.
