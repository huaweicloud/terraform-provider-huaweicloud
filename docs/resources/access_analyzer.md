---
subcategory: "Access Analyzer"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_access_analyzer"
description: |-
  Manages an Access Analyzer resource within HuaweiCloud.
---

# huaweicloud_access_analyzer

Manages an Access Analyzer resource within HuaweiCloud.

## Example Usage

```hcl
variable "analyzer_name" {}

resource "huaweicloud_access_analyzer" "test" {
  name = var.analyzer_name
  type = "account"
  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, NonUpdatable) Specifies the name of the analyzer.

* `type` - (Required, String, NonUpdatable) Specifies the type of the analyzer.
  The value can be: **account**.

* `tags` - (Optional, Map) Specifies the tags of the analyzer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the analyzer.

* `created_at` - The time when the analyzer was created.

* `urn` - The resource analyzer identifier.

* `last_analyzed_resource` - The URN of last analyzed resource.

* `last_resource_analyzed_at` - The time when the resource was last analyzed.

## Import

Analyzers can be imported by their `id`, e.g.

```bash
$ terraform import huaweicloud_access_analyzer.test 3b7e65af-e75b-4d78-ac75-2a87924cd2a2
```
