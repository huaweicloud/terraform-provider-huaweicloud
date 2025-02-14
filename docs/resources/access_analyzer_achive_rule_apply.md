---
subcategory: "Access Analyzer"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_access_analyzer_archive_rule_apply"
description: |-
  Use this resource to apply an Access Analyzer archive rule within HuaweiCloud.
---

# huaweicloud_access_analyzer_archive_rule_apply

Use this resource to apply an Access Analyzer archive rule within HuaweiCloud.

## Example Usage

```hcl
variable "analyzer_id" {}
variable "archive_rule_id" {}

resource "huaweicloud_access_analyzer_archive_rule_apply" "test" {
  analyzer_id     = var.analyzer_id
  archive_rule_id = var.archive_rule_id
}
```

~> This resource is only a one-time action resource for archive rule apply resource. Deleting this resource will
  not change the status of the currently archive rule apply resource, but will only remove the resource information
  from the tfstate file.

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `analyzer_id` - (Required, String, NonUpdatable) Specifies the ID of the analyzer to which the archive rule belongs.

* `archive_rule_id` - (Required, String, NonUpdatable) Specifies ID of the archive rule to apply.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the `archive_rule_id`.
