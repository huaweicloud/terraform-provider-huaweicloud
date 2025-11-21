---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_associated_resources"
description: |-
  Use this data source to get the list of associated resources of EPS resource within HuaweiCloud.
---

# huaweicloud_enterprise_associated_resources

Use this data source to get the list of associated resources of EPS resource within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_enterprise_project_associated_resources" "test" {
  resource_id   = "9ff4bcd5-88b7-4c28-911f-6d5be757024b"
  project_id    = "0bc8d0ad6980f56e2f30c00bf75b72b2"
  region_id     = "cn-north-4"
  resource_type = "ecs"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Specifies the resource id.  

* `project_id` - (Required, String) Specifies the project id.

* `region_id` - (Required, String) Specifies the region id.

* `resource_type` - (Required, String) Specifies the resource type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID.

* `name` - Resource name.

* `type` - Type of associated resources.

* `associated_resources` - The list of the associated resources.
  The [associated_resource](#associated_resource) structure is documented below.

* `errors` - Indicates the list of errors.
  The [error](#errors_struct) structure is documented below.

<a name="associated_resource"></a>
The `associated_resource` block supports:

* `id` - Resource ID.

* `name` - Resource name.

* `eip` - EIP information.

* `resource_type` - Resource type.

<a name="errors_struct"></a>
The `error` block supports:

* `error_code` - Indicates the error code.

* `error_msg` - Indicates the error message.

* `project_id` - Indicates the project ID.

* `resource_type` - Indicates the resource type.

* `resource_id` - Indicates the resource id.
