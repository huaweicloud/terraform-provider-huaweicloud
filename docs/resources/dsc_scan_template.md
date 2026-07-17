---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_scan_template"
description: |-
  Manages a DSC scan template resource within HuaweiCloud.
---

# huaweicloud_dsc_scan_template

Manages a DSC scan template resource within HuaweiCloud.

## Example Usage

### Create a new template with built-in rules

```hcl
resource "huaweicloud_dsc_scan_template" "test" {
  action             = "ADD"
  name               = "test_template"
  description        = "test_template_description"
  add_built_in_rules = true
}
```

### Copy an existing template

```hcl
variable "origin_template_id" {}

resource "huaweicloud_dsc_scan_template" "test" {
  name               = "copied_template"
  action             = "COPY"
  description        = "test_template_description"
  origin_template_id = var.origin_template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the template name.

* `description` - (Optional, String) Specifies the template description.

* `action` - (Optional, String, NonUpdatable) Specifies the operation type, such as creating or copying a template.

* `add_built_in_rules` - (Optional, Bool, NonUpdatable) Specifies whether to add built-in rules when creating the template.

* `origin_template_id` - (Optional, String, NonUpdatable) Specifies the origin template ID, used when copying a template.

* `is_default` - (Optional, Bool) Specifies whether the template is the default template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the template ID in UUID format.

* `category` - The template category, such as 'BUILT_IN' or 'CUSTOM'.

* `create_time` - The creation time.

* `update_time` - The update time.

## Import

The scan template can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_dsc_scan_template.test <template_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include `action`, `add_built_in_rules` and `origin_template_id`. It is generally
recommended running `terraform plan` after importing a resource. You can then decide if changes should be applied to the
resource, or the resource definition should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dsc_scan_template" "test" {
  ...

  lifecycle {
    ignore_changes = [
      action,
      add_built_in_rules,
      origin_template_id,
    ]
  }
}
```
