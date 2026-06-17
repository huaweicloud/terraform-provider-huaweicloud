---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_layout"
description: |-
  Manages a SecMaster layout resource within HuaweiCloud.
---

# huaweicloud_secmaster_layout

Manages a SecMaster layout resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "name" {}

resource "huaweicloud_secmaster_layout" "test" {
  workspace_id = var.workspace_id
  name         = var.name
  used_by      = "DATACLASS"
  layout_type  = "List"
  binding_code = "Alert"
  boa_version  = "v3"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `name` - (Required, String) Specifies the layout name.

* `used_by` - (Required, String, NonUpdatable) Specifies the business type that uses the layout.
  The value can be **DATACLASS**, **AOP_WORKFLOW**, **SECURITY_REPORT**, or **DASHBOARD**.

* `description` - (Optional, String) Specifies the description.

* `cloud_pack_id` - (Optional, String) Specifies the cloud pack ID.

* `cloud_pack_name` - (Optional, String) Specifies the cloud pack name.

* `cloud_pack_version` - (Optional, String) Specifies the cloud pack version.

* `layout_json` - (Optional, String) Specifies the layout information in JSON format.

* `region_id` - (Optional, String, NonUpdatable) Specifies the region ID.

* `domain_id` - (Optional, String, NonUpdatable) Specifies the domain ID.

* `thumbnail` - (Optional, String, NonUpdatable) Specifies the template thumbnail.

* `layout_type` - (Optional, String, NonUpdatable) Specifies the layout type.  
  + When `used_by` is **SECURITY_REPORT**, or **DASHBOARD**, this field is not returned.

* `binding_id` - (Optional, String, NonUpdatable) Specifies the data class ID or workflow ID.  
  + When `used_by` is **SECURITY_REPORT**, or **DASHBOARD**, this field is not returned.

* `binding_code` - (Optional, String, NonUpdatable) Specifies the data class business code.  
  + When `used_by` is **SECURITY_REPORT**, or **DASHBOARD**, this field is not returned.

* `fields_sum` - (Optional, Int) Specifies the total number of fields.  
  + When `used_by` is **SECURITY_REPORT**, or **DASHBOARD**, this field is not returned.

* `wizards_sum` - (Optional, Int) Specifies the total number of pages.  
  + When `used_by` is **SECURITY_REPORT**, or **DASHBOARD**, this field is not returned.

* `sections_sum` - (Optional, Int) Specifies the total number of system sections.

* `tabs_sum` - (Optional, Int) Specifies the total number of custom tabs.

* `boa_version` - (Optional, String) Specifies the BOA version.

* `is_delete` - (Optional, Bool) Specifies whether to directly delete the layout.  
  The valid values are as follows:
  + **true**: Indicates direct deletion.
  + **false**: If there is a reference relationship, deletion fails.

  Defaults to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (the layout ID).

* `create_time` - The creation time.

* `creator_name` - The creator name.

* `creator_id` - The creator ID.

* `parent_id` - The parent layout ID.

* `en_description` - The English description.

* `en_name` - The English name.

* `project_id` - The project ID.

* `update_time` - The update time.

* `layout_cfg` - The layout configuration used to bind icons on the frontend.

* `binding_name` - The data class name or workflow name.

* `modules_sum` - The total number of system modules.

* `version` - The SecMaster version.

## Import

The layout can be imported using the `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_layout.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `is_delete`, `fields_sum`, `wizards_sum`.
It is generally recommended running `terraform plan` after importing a layout.
You can then decide if changes should be applied to the layout, or the resource definition should be updated to
align with the layout. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_layout" "test" {
  ...

  lifecycle {
    ignore_changes = [
      is_delete,
      fields_sum,
      wizards_sum,
    ]
  }
}
```
