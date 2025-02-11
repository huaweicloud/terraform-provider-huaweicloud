---
subcategory: "Access Analyzer"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_access_analyzers"
description: |-
  Use this data source to get a list of Access Analyzers within HuaweiCloud.
---

# huaweicloud_access_analyzers

Use this data source to get a list of Access Analyzers within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_access_analyzers" "test" {
  type = "account"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the type of the analyzer.
  The value can be: **account**, **organization**, **account_unused_access** and **organization_unused_access**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data souce ID.

* `analyzers` - The list of analyzers.
  The [analyzers](#analyzers) structure is documented below.

<a name="analyzers"></a>
The `analyzers` block supports:

* `name` - The name of the analyzer.

* `type` - The type of the analyzer.

* `configuration` - The configuration of the analyzer.
  The [configuration](#configuration) structure is documented below.

* `tags` - The tags of the analyzer.

* `status` - The status of the analyzer.

* `status_reason` - The status reason of the analyzer.
  The [status_reason](#status_reason) structure is documented below.

* `created_at` - The time when the analyzer was created.

* `urn` - The resource analyzer identifier.

* `organization_id` - The organization ID of the analyzer.

* `last_analyzed_resource` - The URN of last analyzed resource.

* `last_resource_analyzed_at` - The time when the resource was last analyzed.

<a name="configuration"></a>
The `configuration` block supports:

* `unused_access` - The unused access.
  The [unused_access](#unused_access) structure is documented below.

<a name="unused_access"></a>
The `unused_access` block supports:

* `unused_access_age` - The unused access age in days.

<a name="status_reason"></a>
The `status_reason` block supports:

* `code` - The code of the status reason.

* `details` - The details of the status reason.
