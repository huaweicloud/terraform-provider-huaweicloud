---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_resources"
description: |-
  Use this data source to get the list of CCI resources within HuaweiCloud.
---

# huaweicloud_cciv2_resources

Use this data source to get the list of CCI resources within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cciv2_resources" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The resources.
  The [resources](#attrblock_resources) structure is documented below.

<a name="attrblock_resources"></a>
The `resources` block supports:

* `categories` - The categories.

* `group` - The group.

* `kind` - The kind.

* `name` - The name.

* `namespaced` - The namespaced.

* `short_names` - The short names.

* `singular_name` - The singular name.

* `storage_version_hash` - The storage version hash.

* `verbs` - The verbs.

* `version` - The version.
