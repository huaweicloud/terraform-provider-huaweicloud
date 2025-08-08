---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_catalogues_search"
description: |-
  Use this data source to query the SecMaster catalogues within HuaweiCloud.
---

# huaweicloud_secmaster_catalogues_search

Use this data source to query the SecMaster catalogues within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_catalogues_search" "example" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the catalogues.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to query catalogues.

* `parent_catalogue` - (Optional, String) Specifies the parent catalogue name.

* `second_catalogue` - (Optional, String) Specifies the second-level catalogue name.

* `catalogue_status` - (Optional, Bool) Specifies the status of the catalogue. Defaults to **false**.

* `layout_name` - (Optional, String) Specifies the layout name.

* `publisher_name` - (Optional, String) Specifies the publisher name.

* `analysis_version` - (Optional, String) Specifies the analysis version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of catalogues that match the query criteria.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The ID of the catalogue.

* `parent_catalogue` - The name of the parent catalogue.

* `second_catalogue` - The name of the second-level catalogue.

* `catalogue_status` - The status of the catalogue (builtin/custom).

* `catalogue_address` - The address of the catalogue.

* `layout_id` - The ID of the layout.

* `layout_name` - The name of the layout.

* `publisher_name` - The name of the publisher.

* `is_card_area` - Whether to display the card area.

* `is_display` - Whether to display the catalogue.

* `is_landing_page` - Whether it is a landing page.

* `is_navigation` - Whether to display navigation.

* `parent_alias_en` - The English alias of the parent catalogue.

* `parent_alias_zh` - The Chinese alias of the parent catalogue.

* `second_alias_en` - The English alias of the second-level catalogue.

* `second_alias_zh` - The Chinese alias of the second-level catalogue.
