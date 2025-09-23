---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_rules"
description: |-
  Use this data source to get a list of CodeArts pipeline rules.
---

# huaweicloud_codearts_pipeline_rules

Use this data source to get a list of CodeArts pipeline rules.

## Example Usage

```hcl
data "huaweicloud_codearts_pipeline_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Optional, String) Specifies the CodeArts project ID.

* `name` - (Optional, String) Specifies the rule name.

* `type` - (Optional, String) Specifies the rule type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - Indicates the rule list.
  The [rules](#attrblock--rules) structure is documented below.

<a name="attrblock--rules"></a>
The `rules` block supports:

* `id` - Indicates the rule ID.

* `name` - Indicates the rule name.

* `type` - Indicates the rule type.

* `version` - Indicates the rule version.

* `operate_time` - Indicates the operate time.

* `operator` - Indicates the operator.
