---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_cluster_component_batch_add"
description: |-
  Use this resource to batch add components to the specified MRS cluster within HuaweiCloud.
---

# huaweicloud_mapreduce_cluster_component_batch_add

Use this resource to batch add components to the specified MRS cluster within HuaweiCloud.

~> If you use this resource, please use `lifecycle.ignore_changes` to ignore the changes of `component_list`,
   `analysis_core_nodes[0].assigned_roles`, `streaming_core_nodes[0].assigned_roles`,
   `analysis_task_nodes[0].assigned_roles`, `streaming_task_nodes[0].assigned_roles`, and
   `custom_nodes[0].assigned_roles` in `huaweicloud_mapreduce_cluster`.

-> 1. Only MRS `3.1.2` and later normal versions and MRS `3.1.2-LTS.2` and later LTS versions of **custom clusters**
   support adding components.
   <br>2. This resource is a one-time action resource used to batch add components to the MRS cluster. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information
   from the tfstate file.

## Example Usage

```hcl
variable "components" {
  type = list(object({
    component = string

    node_groups = list(object({
      name           = string
      assigned_roles = list(string)
    }))

    user_password    = optional(string)
    default_password = optional(string)
  }))
}

resource "huaweicloud_mapreduce_cluster_component_batch_add" "test" {
  cluster_id = huaweicloud_mapreduce_cluster.test.id

  dynamic "components_install_mode" {
    for_each = var.components

    content {
      component = components_install_mode.value.component

      dynamic "node_groups" {
        for_each = components_install_mode.value.node_groups

        content {
          name           = node_groups.value.name
          assigned_roles = node_groups.value.assigned_roles
        }
      }

      component_user_password    = components_install_mode.value.user_password
      component_default_password = components_install_mode.value.default_password
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the components to be added are located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the MRS cluster.

* `components_install_mode` - (Required, List, NonUpdatable) Specifies the list of components to be added.  
  The [components_install_mode](#cluster_component_batch_add_components) structure is documented below.

<a name="cluster_component_batch_add_components"></a>
The `components_install_mode` block supports:

* `component` - (Required, String, NonUpdatable) Specifies the name of the component.

* `node_groups` - (Required, List, NonUpdatable) Specifies the node groups where the component roles will be deployed.  
  The [node_groups](#cluster_component_batch_add_node_groups) structure is documented below.

* `component_user_password` - (Optional, String, NonUpdatable) Specifies the password for the component user.  
  This password is used for the ClickHouse component machine user to connect.

* `component_default_password` - (Optional, String, NonUpdatable) Specifies the password for the component
  default user.  
  This password is used for the ClickHouse component machine default user to connect.

<a name="cluster_component_batch_add_node_groups"></a>
The `node_groups` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of the node group.

* `assigned_roles` - (Required, List, NonUpdatable) Specifies the list of roles to be assigned to this node group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following Timeouts configuration options:

* `create` - Default is 30 minutes.
