---
subcategory: "Access Analyzer"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_access_analyzer_archive_rule"
description: |-
  Manages an Access Analyzer archive rule resource within HuaweiCloud.
---

# huaweicloud_access_analyzer_archive_rule

Manages an Access Analyzer archive rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "analyzer_id" {}
variable "rule_name" {}

resource "huaweicloud_access_analyzer_archive_rule" "test" {
  analyzer_id = var.analyzer_id
  name        = var.rule_name

  filters {
    key = "resource_type"
    criterion {
        eq = ["iam:agency", "obs:bucket"]
    }
  }

  filters {
    key = "condition.g:SourceVpc"
    criterion {
        exists = true
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `analyzer_id` - (Required, String, NonUpdatable) Specifies the ID of the analyzer to which the archive rule belongs.

* `name` - (Required, String, NonUpdatable) Specifies the name of the archive rule.

* `filters` - (Required, List) Specifies the filters of the archive rule.
  The [filters](#filters) structure is documented below.

<a name="filters"></a>
The `filters` block supports:

* `key` - (Required, String) Specifies the key of the filter. The value can be: **resource**, **resource_type**,
  **resource_owner_account**, **is_public**, **id**, **status**, **principal_type**, **principal_identifier**,
  **change_type**, **existing_finding_id**, **existing_finding_status**, **condition.g:PrincipalUrn**,
  **condition.g:PrincipalId** and **condition.g:PrincipalAccount**.

* `criterion` - (Required, List, NonUpdatable) Specifies the criterion of the filter.
  The [criterion](#criterion) structure is documented below.

* `organization_id` - (Optional, String) Specifies the organization ID of the filter.

<a name="criterion"></a>
The `criterion` block supports:

* `contains` - (Optional, List) Specifies the values of the **contains** operator.

* `eq` - (Optional, List) Specifies the values of the **equals** operator.

* `neq` - (Optional, List) Specifies the values of the **not equals** operator.

* `exists` - (Optional, String) Specifies the values of the **exists** operator.
  The value can be: **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - The resource archive rule identifier.

* `created_at` - The time when the archive rule was created.

* `updated_at` - The time when the archive rule was updated.

## Import

Archive rules can be imported by the `analyzer_id` and `id`, e.g.

```bash
$ terraform import huaweicloud_access_analyzer_archive_rule.test <analyzer_id>/<id>
```
