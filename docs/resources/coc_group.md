---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_group"
description: |-
  Manages a COC group resource within HuaweiCloud.
---

# huaweicloud_coc_group

Manages a COC group resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "component_id" {}
variable "region_id" {}

resource "huaweicloud_coc_group" "test" {
  name         = var.name
  component_id = var.component_id
  region_id    = var.region_id
  sync_mode    = "AUTO"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the group name.

* `component_id` - (Required, String, NonUpdatable) Specifies the component ID.

* `region_id` - (Required, String, NonUpdatable) Specifies the region ID.

* `sync_mode` - (Required, String) Specifies the resource synchronization mode. The default value is **MANUAL**.
  Values can be as follows:
  + **MANUAL**: Manual association: Under the corresponding group, the user manually associates the corresponding
  resource data to the group for management.
  + **AUTO**: Smart association: Users can group resources with the same tag under an enterprise project into the same
  resource group using enterprise projects and tags.

* `vendor` - (Optional, String, NonUpdatable) Specifies the manufacturer information.
  Values can be as follows:
  + **RMS**: Huawei Cloud Vendor.
  + **ALI**: Alibaba Cloud Vendor.
  + **OTHER**: Other Vendor.

* `application_id` - (Optional, String, NonUpdatable) Specifies the application ID.

* `sync_rules` - (Optional, List) Specifies the smart association rules.
  The [sync_rules](#sync_rules_struct) structure is documented below.

  -> Only the `sync_rules` is only valid when `sync_mode` is **AUTO**.

* `relation_configurations` - (Optional, List) Specifies the group configuration information.
  The [relation_configurations](#relation_configurations_struct) structure is documented below.

* `force_delete` - (Optional, Bool) Specifies whether to force deletion.
  Values can be as follows:
  + **true**: If you force delete a group, the system will immediately clear all associated resources and delete the group.
  + **false**: Non-force deletion only supports deleting groups without associated resources.

<a name="sync_rules_struct"></a>
The `sync_rules` block supports:

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `rule_tags` - (Optional, String) Specifies the associated tag.

<a name="relation_configurations_struct"></a>
The `relation_configurations` block supports:

* `type` - (Optional, String) Specifies the configuration type.

* `parameters` - (Optional, Map) Specifies the configuration parameters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `code` - Indicates the group code.

* `enterprise_project_id` - Indicates the enterprise project ID.

## Import

The COC group can be imported using `component_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_coc_group.test <component_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `force_delete`.
It is generally recommended running `terraform plan` after importing a group.
You can then decide if changes should be applied to the group, or the resource definition should be updated to align
with the group. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_group" "test" {
    ...

  lifecycle {
    ignore_changes = [
      force_delete
    ]
  }
}
```
