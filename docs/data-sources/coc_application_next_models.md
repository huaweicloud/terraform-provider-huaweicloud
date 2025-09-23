---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_application_next_models"
description: |-
  Use this data source to get the list of COC application next models.
---

# huaweicloud_coc_application_next_models

Use this data source to get the list of COC application next models.

## Example Usage

### query current application next models

```hcl
variable "application_id" {}

data "huaweicloud_coc_application_next_models" "test" {
  application_id = var.application_id
}
```

### query current component next models

```hcl
variable "component_id" {}

data "huaweicloud_coc_application_next_models" "test" {
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional, String) Specifies the application ID.

* `component_id` - (Optional, String) Specifies the component ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sub_applications` - Indicates the list of the sub applications.

  The [sub_applications](#sub_applications_struct) structure is documented below.

* `components` - Indicates the list of the components.

  The [components](#components_struct) structure is documented below.

* `groups` - Indicates the list of the groups.

  The [groups](#groups_struct) structure is documented below.

<a name="sub_applications_struct"></a>
The `sub_applications` block supports:

* `id` - Indicates the application ID.

* `name` - Indicates the application name.

* `code` - Indicates the application code.

* `description` - Indicates the application description.

* `domain_id` - Indicates the tenant ID.

* `parent_id` - Indicates the parent ID.

* `path` - Indicates the node path.

* `create_time` - Indicates the creation time.

* `update_time` - Indicates the modification time.

<a name="components_struct"></a>
The `components` block supports:

* `id` - Indicates the component ID.

* `name` - Indicates the component name.

* `code` - Indicates the component code.

* `domain_id` - Indicates the account ID.

* `application_id` - Indicates the application ID.

* `path` - Indicates the component node path.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - Indicates the group ID.

* `name` - Indicates the group name.

* `code` - Indicates the group code.

* `domain_id` - Indicates the account ID.

* `region_id` - Indicates the region ID.

* `application_id` - Indicates the application ID.

* `component_id` - Indicates the component ID.

* `sync_mode` - Indicates the resource association method.

* `vendor` - Indicates the manufacturer.

* `sync_rules` - Indicates the intelligent association rules.

  The [sync_rules](#groups_sync_rules_struct) structure is documented below.

* `relation_configurations` - Indicates the group configuration information.

<a name="groups_sync_rules_struct"></a>
The `sync_rules` block supports:

* `enterprise_project_id` - Indicates the enterprise project ID.

* `rule_tags` - Indicates the associated tag.
