---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_components"
description: |-
  Use this data source to get the list of SecMaster soc components within HuaweiCloud.
---

# huaweicloud_secmaster_soc_components

Use this data source to get the list of SecMaster soc components within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_soc_components" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The soc components list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The ID.

* `name` - The name.

* `dev_language` - The dev language.

* `dev_language_version` - The dev language version.

* `alliance_id` - The alliance ID.

* `alliance_name` - The alliance name.

* `description` - The description.

* `logo` - The market head icon.

* `label` - The label.

* `create_time` - The creation time.

* `update_time` - The update time.

* `creator_name` - The creator name.

* `operate_history` - The plugin operation history.

  The [operate_history](#operate_history_struct) structure is documented below.

* `component_versions` - The plugin version information, compatible with previous Java versions with plugin granularity.

  The [component_versions](#component_versions_struct) structure is documented below.

* `component_type` - The plugin type.  
  The valid values are as follows:
  + **subscribe**: Subscription.
  + **custom**: Custom development.
  + **system**: System built-in.

<a name="operate_history_struct"></a>
The `operate_history` block supports:

* `operate_name` - The operation name.

* `operate_time` - The operating time.

<a name="component_versions_struct"></a>
The `component_versions` block supports:

* `id` - The ID.

* `version_num` - The version.

* `version_desc` - The version description.

* `status` - The status.

* `package_name` - The package name.

* `component_attachment_id` - The component attachment ID.

* `sub_version_id` - The sub-version ID.

* `connection_config` - The operational connection configuration list.

  The [connection_config](#connection_config_struct) structure is documented below.

<a name="connection_config_struct"></a>
The `connection_config` block supports:

* `default_value` - The default value.

* `description` - The description.

* `key` - The key.

* `name` - The connection config name.

* `readonly` - The readonly.

* `required` - Is it required.

* `type` - The type.

* `value` - The value.
