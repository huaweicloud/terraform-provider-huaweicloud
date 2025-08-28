---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_groups"
description: |-
  Use this data source to get the list of COC groups.
---

# huaweicloud_coc_groups

Use this data source to get the list of COC groups.

## Example Usage

```hcl
variable "component_id" {}

data "huaweicloud_coc_groups" "test" {
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `component_id` - (Required, String) Specifies the component ID.

* `id_list` - (Optional, List) Specifies the group ID list.

* `application_id` - (Optional, String) Specifies the application ID.

* `name_like` - (Optional, String) Specifies the fuzzy query the group name.

* `code` - (Optional, String) Specifies the group code.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the application group list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the UUID assigned by the CMDB.

* `name` - Indicates the group name.

* `vendor` - Indicates the manufacturer information.

* `code` - Indicates the group code.

* `domain_id` - Indicates the tenant ID.

* `region_id` - Indicates the region ID.

* `component_id` - Indicates the component ID.

* `application_id` - Indicates the application ID.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `sync_mode` - Indicates the resource synchronization mode.

* `sync_rules` - Indicates the smart association rules.

  The [sync_rules](#data_sync_rules_struct) structure is documented below.

* `relation_configurations` - Indicates the group configuration information.

  The [relation_configurations](#data_relation_configurations_struct) structure is documented below.

<a name="data_sync_rules_struct"></a>
The `sync_rules` block supports:

* `enterprise_project_id` - Indicates the enterprise project ID.

* `rule_tags` - Indicates the associated tag.

<a name="data_relation_configurations_struct"></a>
The `relation_configurations` block supports:

* `type` - Indicates the configuration type.

* `parameters` - Indicates the configuration parameters.
