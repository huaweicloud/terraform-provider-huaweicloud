---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_component_history_configuration"
description: |-
  Use this data source to get the list of component history configuration.
---

# huaweicloud_secmaster_component_history_configuration

Use this data source to get the list of component history configuration.

## Example Usage

```hcl
variable "workspace_id" {}
variable "component_id" {}

data "huaweicloud_secmaster_component_history_configuration" "test" {
  workspace_id = var.workspace_id
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `component_id` - (Required, String) Specifies the component ID.

* `sort_key` - (Optional, String) Specifies the sort field.

* `sort_dir` - (Optional, String) Specifies the sort direction.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The component history configuration list.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `configuration_status` - The node configuration status.  
  The valid values are as follows:
  + **UN_SAVED**: Not saved.
  + **SAVE_AND_UN_DEPLOY**: Saved but not deployed.
  + **DEPLOYING**: Deploying.
  + **MOVE_AND_UN_DEPLOY**: Removed but not applied.
  + **FAIL_DEPLOY**: Deployment failed.
  + **DEPLOYED**: Deployed.

* `list` - The file parameter information list.

  The [list](#list_struct) structure is documented below.

* `node_id` - The node ID.

* `node_name` - The node name.

* `node_status` - The node status.  
  The valid values are as follows:
  + **NORMAL**: Normal.
  + **ANOMALIES**: Abnormal.
  + **FAULTS**: Fault.
  + **LOST_CONTACT**: Lost contact.

* `specification` - The node specification.

<a name="list_struct"></a>
The `list` block supports:

* `configuration_id` - The configuration ID.

* `file_name` - The file name.

* `file_type` - The file type.

* `node_id` - The node ID.

* `param` - The parameter.

* `type` - The type.

* `version` - The version.
