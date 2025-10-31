---
subcategory: "RGC"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_template"
description: |-
  Manages an RGC template resource within HuaweiCloud.
---

# huaweicloud_rgc_template

Manages an RGC template resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable template_name {}
variable template_type {}

resource "huaweicloud_rgc_template" "predefined_template" {
  template_name = var.template_name
  template_type = var.template_type
}
```

### Customized template

```hcl
variable template_name {}
variable template_type {}
variable template_description {}
variable template_body {}

resource "huaweicloud_rgc_template" "customized_template" {
  template_name        = var.template_name
  template_type        = var.template_type
  template_description = var.template_description
  template_body        = var.template_body
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `template_name` - (Required, String, NonUpdatable) Specifies the name to create the template.

* `template_type` - (Required, String, NonUpdatable) Specifies the type of the template,
  only **predefined** and **customized** are supported.

* `template_description` - (Optional, String, NonUpdatable) Specifies the description of customized template.

* `template_body` - (Optional, String, NonUpdatable) Specifies the content of customized template,
  it is a zip-type compressed file that has been encoded using base64.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `latest_version_id` - The lastest version ID.

* `create_time` - The creation time.

* `update_time` - The last update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.

* `delete` - Default is 5 minutes.

## Import

The RGC template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rgc_template.test <id>
```
