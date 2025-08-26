---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_components"
description: |-
  Use this data source to get the list of COC components.
---

# huaweicloud_coc_components

Use this data source to get the list of COC components.

## Example Usage

```hcl
data "huaweicloud_coc_components" "test" {}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional, String) Specifies the application ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the component query information list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the UUID assigned by the CMDB.

* `name` - Indicates the component name.

* `code` - Indicates the component code.

* `application_id` - Indicates the application ID.

* `ep_id` - Indicates the enterprise project ID.
