---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_components"
description: |-
  Use this data source to get the list of CAE components within HuaweiCloud.
---

# huaweicloud_cae_components

Use this data source to get the list of CAE components within HuaweiCloud.

## Example Usage

### Query all components under the default enterprise project or EPS service is not enable

```hcl
variable "environment_id" {}
variable "application_id" {}

data "huaweicloud_cae_components" "test" {
  environment_id = var.environment_id
  application_id = var.application_id
}
```

### Query all components under the specified enterprise project

```hcl
variable "environment_id" {}
variable "application_id" {}
variable "enterprise_project_id" {}

data "huaweicloud_cae_components" "test" {
  environment_id        = var.environment_id
  application_id        = var.application_id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the components are located.  
  If omitted, the provider-level region will be used.

* `environment_id` - (Required, String) Specifies the ID of the environment to which the components belong.

* `application_id` - (Required, String) Specifies the ID of the application to which the components belong.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the components
  belong.  
  If the `application_id` belongs to the non-default enterprise project, this parameter is required and is only valid
  for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `components` - All queried components.  
  The [components](#cae_components) structure is documented below.

<a name="cae_components"></a>
The `components` block supports:

* `id` - The ID of the component.

* `name` - The name of the component.

* `annotations` - The parameters of key/value pairs related to the component.

* `spec` - The configuration information of the component.  
  The [spec](#cae_components_spec) structure is documented below.

* `created_at` - The creation time of the component, in RFC3339 format.

* `updated_at` - The latest update time of the component, in RFC3339 format.

<a name="cae_components_spec"></a>
The `spec` block supports:

* `runtime` - The component runtime.

* `environment_id` - The ID of the environment to which the component belongs.

* `replica` - The instance number of the component.

* `available_replica` - The available instance number of the component.

* `source` - The code source configuration information corresponding to the component.

* `build` - The build information of the code source corresponding to the component.

* `resource_limit` - The instance specification corresponding to the component.

* `image_url` - The image URL that component used.

* `status` - The status of the component.
