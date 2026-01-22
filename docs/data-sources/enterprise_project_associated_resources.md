---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project_associated_resources"
description: |-
  Use this data source to get the list of associated resources of EPS resource within HuaweiCloud.
---

# huaweicloud_enterprise_project_associated_resources

Use this data source to get the list of associated resources of EPS resource within HuaweiCloud.

## Example Usage

```hcl
variable "resource_id" {}
variable "resource_type" {}

data "huaweicloud_enterprise_project_associated_resources" "test" {
  resource_id   = var.resource_id
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Specifies the resource id.  

* `resource_type` - (Required, String) Specifies the resource type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the value of `resource_id`).

* `name` - The resource name.

* `type` - The type of associated resource.

* `associated_resources` - The list of the associated resources.
  The [associated_resource](#associated_resource_struct) structure is documented below.

* `errors` - The list of errors.
  The [errors](#errors_struct) structure is documented below.

<a name="associated_resource_struct"></a>
The `associated_resource` block supports:

* `id` - The resource ID.

* `name` - The resource name.

* `eip` - The EIP information.

* `resource_type` - The resource type.

<a name="errors_struct"></a>
The `errors` block supports:

* `project_id` - Indicates the project ID.

* `resource_type` - Indicates the resource type.

* `resource_id` - Indicates the resource id.

* `error_code` - Indicates the error code.

* `error_msg` - Indicates the error message.
