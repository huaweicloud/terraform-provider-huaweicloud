---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_catalogues"
description: |-
  Use this data source to get the list of SecMaster catalogues.
---

# huaweicloud_secmaster_catalogues

Use this data source to get the list of SecMaster catalogues.
This data source provides detailed information about all catalogues in the specified workspace.

## Example Usage

```hcl
variable "workspace_id" {}
variable "catalogue_type" {}
variable "catalogue_code" {}

# Get all catalogues
data "huaweicloud_secmaster_catalogues" "all" {
  workspace_id = var.workspace_id
}

# Get specific catalogues by type and code
data "huaweicloud_secmaster_catalogues" "filtered" {
  workspace_id   = var.workspace_id
  catalogue_type = var.catalogue_type
  catalogue_code = var.catalogue_code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `catalogue_type` - (Optional, String) Specifies the type of catalogues to filter.

* `catalogue_code` - (Optional, String) Specifies the code of catalogues to filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of the catalogues.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The unique identifier of the catalogue.

* `parent_catalogue` - The name of the parent (first-level) catalogue.

* `second_catalogue` - The name of the second-level catalogue.

* `catalogue_status` - Whether the catalogue is built-in.

* `catalogue_address` - The address of the catalogue.

* `layout_id` - The ID of the layout associated with the catalogue.

* `layout_name` - The name of the layout associated with the catalogue.

* `publisher_name` - The name of the publisher.

* `is_card_area` - Whether to display the card area.

* `is_display` - Whether to display the catalogue.

* `is_landing_page` - Whether it is a landing page.

* `is_navigation` - Whether to display the breadcrumb navigation.

* `parent_alisa_en` - The English alias of the parent catalogue.

* `parent_alisa_zh` - The Chinese alias of the parent catalogue.

* `second_alias_en` - The English alias of the second-level catalogue.

* `second_alias_zh` - The Chinese alias of the second-level catalogue.
