---
subcategory: "Access Analyzer"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_access_analyzer_archive_rules"
description: |-
  Use this data source to get a list of Access Analyzers archive rules within HuaweiCloud.
---

# huaweicloud_access_analyzer_archive_rules

Use this data source to get a list of Access Analyzers archive rules within HuaweiCloud.

## Example Usage

```hcl
variable "analyzer_id" {}

data "huaweicloud_access_analyzer_archive_rules" "test" {
  analyzer_id = var.analyzer_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `analyzer_id` - (Required, String) Specifies the ID of the analyzer to which the archive rule belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data souce ID.

* `archive_rules` - The list of archive rules.
  The [archive_rules](#archive_rules) structure is documented below.

<a name="archive_rules"></a>
The `archive_rules` block supports:

* `id` - The archive rule ID.

* `name` - The name of the archive rule.

* `filters` - The filters of the archive rule.
  The [filters](#filters) structure is documented below.

* `urn` - The resource archive rule identifier.

* `created_at` - The time when the archive rule was created.

* `updated_at` - The time when the archive rule was updated.

<a name="filters"></a>
The `filters` block supports:

* `key` - The key of the filter.

* `criterion` - The criterion of the filter.
  The [criterion](#criterion) structure is documented below.

<a name="criterion"></a>
The `criterion` block supports:

* `contains` - The values of the **contains** operator.

* `eq` - The values of the **equals** operator.

* `neq` - The values of the **not equals** operator.

* `exists` - The values of the **exists** operator.
