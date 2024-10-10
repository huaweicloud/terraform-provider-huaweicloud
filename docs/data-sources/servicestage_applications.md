---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestage_applications"
description: |-
  Use this data source to query available applications within HuaweiCloud.
---

# huaweicloud_servicestage_applications

Use this data source to query available applications within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_servicestage_applications" "test" {}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region where the applications are located.  
  If omitted, the provider-level region will be used.

* `ignore_environments` - (Optional, Bool) Specifies whether to ignore environments query.  
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - All queried applications.
  The [applications](#servicestage_applications) structure is documented below.

<a name="servicestage_applications"></a>
The `applications` block contains:

* `id` - The ID of the application.

* `name` - The name of the application.

* `description` - The description of the application.

* `enterprise_project_id` - The enterpeise project ID to which the application belongs.

* `creator` - The creator name.

* `created_at` - The creation time of the application, in RFC3339 format.

* `updated_at` - The latest update time of the application, in RFC3339 format.

* `component_count` - The number of the components associated with the application.

* `environments` - The environment configuration associated with the application.
  The [environments](#environment_configurations_associated_app) structure is documented below.

<a name="environment_configurations_associated_app"></a>
The `environments` block contains:

* `id` - The environment ID.

* `variables` - The variables of the environment.
  The [variables](#environment_variables_associated_app) structure is documented below.

<a name="environment_variables_associated_app"></a>
The `variables` block contains:

* `name` - The variable name.

* `value` - The variable value.
