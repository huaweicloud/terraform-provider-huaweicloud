---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_custom_template"
description: |-
  Manages CSS logstash custom template within HuaweiCloud.
---

# huaweicloud_css_logstash_custom_template

Manages CSS logstash custom template within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "name" {}
variable "configuration_name" {}

resource "huaweicloud_css_logstash_custom_template" "test" {
  cluster_id         = var.cluster_id
  name               = var.name
  configuration_name = var.configuration_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the ID of the CSS logstash cluster.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the custom template.
  Changing this creates a new resource.

* `configuration_name` - (Required, String, ForceNew) Specifies the name of the configuration file you want to
  add to the custom template.
  Changing this creates a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the custom template.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `template_id` - The ID of the custom template.

* `conf_content` - The configuration file content of the custom template.

## Import

The CSS logstash custom template can be imported using `name`, e.g.

```bash
$ terraform import huaweicloud_css_logstash_custom_template.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `cluster_id`, `configuration_name`.
It is generally recommended running `terraform plan` after importing the CSS logstash custom template.
You can then decide if changes should be applied to the CSS logstash custom template,
or the CSS logstash custom template definition should be updated to align with the CSS logstash custom template.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_css_logstash_custom_template" "test" {
    ...

  lifecycle {
    ignore_changes = [
      cluster_id, configuration_name,
    ]
  }
}
```
