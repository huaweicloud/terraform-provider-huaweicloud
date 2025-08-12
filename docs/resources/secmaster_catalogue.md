---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_catalogue"
description: |-
  Manages a SecMaster catalogue resource within HuaweiCloud.
---

# huaweicloud_secmaster_catalogue

Manages a SecMaster catalogue resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "layout_id" {}
variable "catalogue_address" {}
variable "second_catalogue_code" {}

resource "huaweicloud_secmaster_catalogue" "test" {
  workspace_id          = var.workspace_id
  layout_id             = var.layout_id
  catalogue_address     = var.catalogue_address
  second_catalogue_code = var.second_catalogue_code
  parent_catalogue      = "first-level-dir"
  parent_alias_en       = "first-level-dir-en"
  parent_alias_zh       = "一级目录"
  second_catalogue      = "second-level-dir"
  second_alias_en       = "second-level-dir-en"
  second_alias_zh       = "二级目录"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the catalogue belongs.

* `parent_catalogue` - (Optional, String) Specifies the name of the first-level directory.

* `parent_alias_en` - (Optional, String) Specifies the English alias of the first-level directory.

* `parent_alias_zh` - (Optional, String) Specifies the Chinese alias of the first-level directory.

* `second_catalogue` - (Optional, String) Specifies the name of the second-level directory.

* `second_alias_en` - (Optional, String) Specifies the English alias of the second-level directory.

* `second_alias_zh` - (Optional, String) Specifies the Chinese alias of the second-level directory.

* `second_catalogue_code` - (Optional, String) Specifies the code of the second-level directory.
  This field is not returned in the API response body.

* `layout_id` - (Optional, String) Specifies the ID of the layout.

* `layout_name` - (Optional, String) Specifies the name of the layout.

* `catalogue_address` - (Optional, String) Specifies the address of the directory.

* `publisher_name` - (Optional, String) Specifies the name of the publisher.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (the catalogue ID).

* `catalogue_status` - The flag indicating whether it is a built-in directory.

* `is_card_area` - The flag indicating whether to display the card area.

* `is_display` - The flag indicating whether to display the directory.

* `is_landing_page` - The flag indicating whether it is a landing page.

* `is_navigation` - The flag indicating whether to display the breadcrumb navigation.

## Import

The catalogue can be imported using the workspace ID and the catalogue ID, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_catalogue.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `second_catalogue_code`.
It is generally recommended running `terraform plan` after importing a catalogue.
You can then decide if changes should be applied to the catalogue, or the resource definition should be updated to
align with the catalogue. Also you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_catalogue" "test" {
  ...

  lifecycle {
    ignore_changes = [
      second_catalogue_code,
    ]
  }
}
```
