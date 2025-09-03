---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_enterprise_project_collections"
description: |-
  Use this data source to get the list of COC enterprise project collections.
---

# huaweicloud_coc_enterprise_project_collections

Use this data source to get the list of COC enterprise project collections.

## Example Usage

```hcl
data "huaweicloud_coc_enterprise_project_collections" "test" {}
```

## Argument Reference

The following arguments are supported:

* `unique_id` - (Optional, String) Specifies the unique identifier ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the enterprise project collection list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the unique identifier ID.

* `user_id` - Indicates the user ID.

* `ep_id_list` - Indicates the list of enterprise project favorites.
