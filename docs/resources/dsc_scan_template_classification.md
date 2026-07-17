---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_scan_template_classification"
description: |-
  Manages a DSC scan template classification resource within HuaweiCloud.
---

# huaweicloud_dsc_scan_template_classification

Manages a DSC scan template classification resource within HuaweiCloud.

## Example Usage

### Create a parent classification

```hcl
variable "template_id" {}
variable "classification_name" {}

resource "huaweicloud_dsc_scan_template_classification" "test" {
  template_id         = var.template_id
  classification_name = var.classification_name
}
```

### Create a classification under the specified parent classification

```hcl
variable "template_id" {}
variable "classification_name" {}
variable "parent_id" {}

resource "huaweicloud_dsc_scan_template_classification" "test" {
  template_id         = var.template_id
  classification_name = var.classification_name
  parent_id           = var.parent_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the scan template classification is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `template_id` - (Required, String, NonUpdatable) Specifies the ID of the scan template.

* `classification_name` - (Required, String) Specifies the name of the classification.  
  Only Chinese characters, letters, digits, underscores (_) and hyphens (-) are allowed, and it must start with
  a Chinese character and letter.

* `parent_id` - (Optional, String, NonUpdatable) Specifies the parent classification ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `depth` - The depth of the classification.

* `root_id` - The ID of the root to which the classification belongs.

* `create_time` - The time when the classification is created, in RFC3339 format.

* `update_time` - The latest update time of the classification, in RFC3339 format.

## Import

The resource can be imported using the `template_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dsc_scan_template_classification.test <template_id>/<id>
```
