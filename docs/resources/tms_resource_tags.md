---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_resource_tags"
description: ""
---

# huaweicloud_tms_resource_tags

Using this resource to manage tags of other service resources in batches within HuaweiCloud.

~> The `tags` parameters of this resource and each service resource will affect each other, and should be managed in
only one way as much as possible. You can use `lifecycle.ignore_changes` to ignore resource changes.

## Example Usage

```hcl
variable "resources_project_id" {}
variable "resources_configuration" {
  type = list(object({
    type = string
    id   = string
  }))
}

resource "huaweicloud_tms_resource_tags" "test" {
  project_id = var.resources_project_id

  dynamic "resources" {
    for_each = var.resources_configuration

    content {
      resource_type = resources.value["type"]
      resource_id   = resources.value["id"]
    }
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional, String, ForceNew) Specifies project ID to which the managed resources belong.  
  Required if the resources are project level. Changing this will create a new resource.

* `resources` - (Required, List) Specifies the managed resource configuration.  
  The [resources](#tags_resources) structure is documented below.

* `tags` - (Required, Map) Specifies resource tags for batch management.
  + The valid length of the tag key is limited from `1` to `36`, only letters, digits, hyphens (-), underscores (_) and
  Chinese characters are allowed.
  + The valid length of the tag value is limited from `0` to `43`, only letters, digits, periods (.), hyphens (-),
  underscores (_) and Chinese characters are allowed.

<a name="tags_resources"></a>
The `resources` block supports:

* `resource_type` - (Required, String) Specifies the resource type.

* `resource_id` - (Required, String) Specifies the resource ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
