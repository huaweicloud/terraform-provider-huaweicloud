---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_applications"
description: |-
  Use this data source to get the list of COC applications.
---

# huaweicloud_coc_applications

Use this data source to get the list of COC applications.

## Example Usage

```hcl
data "huaweicloud_coc_applications" "test" {}
```

## Argument Reference

The following arguments are supported:

* `id_list` - (Optional, List) Specifies the ID list.

* `parent_id` - (Optional, String) Specifies the parent application ID.

* `name_like` - (Optional, String) Specifies the fuzzy query the application name.

* `code` - (Optional, String) Specifies the application code.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the application list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the application ID.

* `name` - Indicates the application name.

* `code` - Indicates the application code.

* `domain_id` - Indicates the domain ID.

* `parent_id` - Indicates the parent application ID.

* `description` - Indicates the application description.

* `path` - Indicates the application path.

* `is_collection` - Indicates whether the application is a favorite application.

* `update_time` - Indicates the modification time.

* `create_time` - Indicates the creation time.
