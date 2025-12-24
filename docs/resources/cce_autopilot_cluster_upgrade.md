---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_cluster_upgrade"
description: |-
  Use this resource to upgrade a CCE Autopilot cluster within HuaweiCloud.
---

# huaweicloud_cce_autopilot_cluster_upgrade

Use this resource to upgrade a CCE Autopilot cluster within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_cce_autopilot_cluster_upgrade" "test" {
  cluster_id     = var.cluster_id
  target_version = "v1.31"

  strategy {
    type = "inPlaceRollingUpdate"

    in_place_rolling_update {
      user_defined_step = 20
    }
  }
}

```

~> Deleting cluster upgrade resource is not supported, it will only be removed from the state.  
  When the cluster is upgraded, the `version` will be changed. You can ignore the change as below.

```hcl
resource "huaweicloud_cce_autopilot_cluster" "test" {
    ...

  lifecycle {
    ignore_changes = [
      ignore_changes = [ version ]
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the CCE Autopilot cluster upgrade resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `target_version` - (Required, String, NonUpdatable) Specifies the target version, e.g. **v1.29.13-r10**.

* `current_version` - (Optional, String, NonUpdatable) Specifies the current version, e.g. **v1.28.15-r30**.

* `strategy` - (Required, List, NonUpdatable) Specifies the upgrade strategy.
  The [strategy](#strategy) structure is documented below.

* `addons` - (Optional, List, NonUpdatable) Specifies the add-on configuration list
  The [addons](#addons) structure is documented below.

* `node_order` - (Optional, Map, NonUpdatable) Specifies the upgrade sequence of nodes in the node pools.
  The key is the node pool ID, **DefaultPool** indicates the default pool.
  The value is a json string which indicates the priority of nodes in this pool. Please check the example.

* `nodepool_order` - (Optional, Map, NonUpdatable) Specifies the upgrade sequence of node pools, in key-value pairs.
  The key is the node pool ID, **DefaultPool** indicates the default pool.
  The value is the priority of the node pool. **0** indicating the lowest priority.
  A larger value indicates a higher priority.

* `is_snapshot` - (Optional, Bool, NonUpdatable) Specifies whether the cluster is snapshotted.

* `is_postcheck` - (Optional, Bool, NonUpdatable) Specifies whether to run postcheck.

<a name="addons"></a>
The `addons` block supports:

* `addon_template_name` - (Required, String, NonUpdatable) Specifies the add-on name.

* `operation` - (Required, String, NonUpdatable) Specifies the execution action.
  For current upgrades, the value can be **patch**.

* `version` - (Required, String, NonUpdatable) Specifies the target add-on version.
  The target add-on version must match the target cluster version.

* `values` - (Optional, List, NonUpdatable) Specifies the add-on template installation parameters.
  These parameters vary depending on the add-on. The [values](#values) is documented below.

<a name="values"></a>
The `values` block supports:

* `basic_json` - (Optional, String, NonUpdatable) Specifies the json string vary depending on the add-on.

* `custom_json` - (Optional, String, NonUpdatable) Specifies the json string vary depending on the add-on.

* `flavor_json` - (Optional, String, NonUpdatable) Specifies the json string vary depending on the add-on.

~> Arguments which can be passed to the `basic_json`, `custom_json` and `flavor_json` add-on parameters depends on
  the add-on type and version. For more detailed description of add-ons
  see [add-ons description](https://github.com/huaweicloud/terraform-provider-huaweicloud/blob/master/examples/cce/basic/cce-addon-templates.md)

<a name="strategy"></a>
The `strategy` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the upgrade strategy type.
  The value can be **inPlaceRollingUpdate**.

* `in_place_rolling_update` - (Optional, List, NonUpdatable) Specifies the in-place upgrade settings.
  It's mandatory when the `type` is set to **inPlaceRollingUpdate**.
  The [in_place_rolling_update](#in_place_rolling_update) structure is documented below.

<a name="in_place_rolling_update"></a>
The `in_place_rolling_update` block supports:

* `user_defined_step` - (Optional, Int, NonUpdatable) Specifies the node upgrade step.
  The value ranges from **1** to **40**. The recommended value is **20**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
