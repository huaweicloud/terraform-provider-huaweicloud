---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_application_views"
description: |-
  Use this data source to get the list of COC application views.
---

# huaweicloud_coc_application_views

Use this data source to get the list of COC application views.

## Example Usage

```hcl
data "huaweicloud_coc_application_views" "test" {}
```

## Argument Reference

The following arguments are supported:

* `name_like` - (Optional, String) Specifies the fuzzy query application view name.

* `code_list` - (Optional, List) Specifies the application, component and group code list.

* `is_collection` - (Optional, Bool) Specifies whether to add to collection. The default value is **true**.
  Values can be as follows:
  + **true**: Search for apps, components or groups in my favorites.
  + **false**: Search for apps, components or groups in all apps.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of the application views.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the application view ID.

* `name` - Indicates the application, component or group name.

* `code` - Indicates the application, component or group code.

* `type` - Indicates the type.

* `parent_id` - Indicates the parent ID.

* `component_id` - Indicates the component ID.

* `application_id` - Indicates the application ID.

* `path` - Indicates the path where the node is located, consisting of application, component, group and other IDs.

* `vendor` - Indicates cloud vendor information.

* `related_domain_id` - Indicates the domain ID to which the cross-account resource belongs.
