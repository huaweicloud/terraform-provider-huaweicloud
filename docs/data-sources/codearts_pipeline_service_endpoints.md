---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_service_endpoints"
description: |-
  Use this data source to get a list of CodeArts pipeline service endpoints.
---

# huaweicloud_codearts_pipeline_service_endpoints

Use this data source to get a list of CodeArts pipeline service endpoints.

## Example Usage

```hcl
variable "codearts_project_id" {}

data "huaweicloud_codearts_pipeline_service_endpoints" "test" {
  project_id = var.codearts_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `module_id` - (Optional, String) Specifies the module ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `endpoints` - Indicates the endpoint list.
  The [endpoints](#attrblock--endpoints) structure is documented below.

<a name="attrblock--endpoints"></a>
The `endpoints` block supports:

* `id` - Indicates the endpoint ID.

* `created_by` - Indicates the permission information.
  The [created_by](#attrblock--endpoints--created_by) structure is documented below.

* `module_id` - Indicates the module ID.

* `name` - Indicates the endpoint name.

* `url` - Indicates the URL.

<a name="attrblock--endpoints--created_by"></a>
The `created_by` block supports:

* `user_id` - Indicates the user ID.

* `user_name` - Indicates the user name.
