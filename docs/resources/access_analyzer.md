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

* `region` - (Optional, String, ForceNew) Specifies the region in which to create resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the analyzer.

* `type` - (Required, String, NonUpdatable) Specifies the type of the analyzer.
  The value can be: **account**, **organization**, **account_unused_access** and **organization_unused_access**.

* `configuration` - (Optional, List, NonUpdatable) Specifies the configuration of the analyzer.
  The [configuration](#configuration) structure is documented below.

* `tags` - (Optional, Map) Specifies the tags of the analyzer.

<a name="configuration"></a>
The `configuration` block supports:

* `unused_access` - (Optional, List, NonUpdatable) Specifies the unused access.
  The [unused_access](#unused_access) structure is documented below.

<a name="unused_access"></a>
The `unused_access` block supports:

* `unused_access_age` - (Optional, Int, NonUpdatable) Specifies the unused access age in days.
  When the `type` is **account_unused_access** or **organization_unused_access**, the default value is 90.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the analyzer.

* `status_reason` - The status reason of the analyzer.
  The [status_reason](#status_reason) structure is documented below.

* `created_at` - The time when the analyzer was created.

* `urn` - The resource analyzer identifier.

* `organization_id` - The organization ID of the analyzer.

* `last_analyzed_resource` - The URN of last analyzed resource.

* `last_resource_analyzed_at` - The time when the resource was last analyzed.

<a name="status_reason"></a>
The `status_reason` block supports:

* `code` - The code of the status reason.

* `details` - The details of the status reason.

## Import

Analyzers can be imported by their `id`, e.g.

```bash
$ terraform import huaweicloud_access_analyzer.test 3b7e65af-e75b-4d78-ac75-2a87924cd2a2
```
