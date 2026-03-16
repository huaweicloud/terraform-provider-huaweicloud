---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_cluster_default_tags_switch"
description: |-
  Manages an MRS cluster default tags switch resource within HuaweiCloud.
---

# huaweicloud_mapreduce_cluster_default_tags_switch

Manages an MRS cluster default tags switch resource within HuaweiCloud.

~> 1. Using this resource to enable the default tags switch will add default tags to the cluster and each node, and occupy
  `2` tag quotas.
  <br>2. Destroying this resource will remove the default tags that have been added to the cluster and each node.
  <br>3. Using this resource will affect the `tags` parameter in the `huaweicloud_mapreduce_cluster` resource. You can use
   `lifecycle.ignore_changes` or manually synchronize the changes of this resource.

## Example Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_mapreduce_cluster_default_tags_switch" "test" {
  cluster_id = var.cluster_id
  action     = "create"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the cluster default tags switch is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster.

* `action` - (Required, String, NonUpdatable) Specifies the type of the action.
  The valid values are as follows:
  + **create**: Enable the default tags switch, and add default tags to the cluster and each node.
  + **delete**: Disable the default tags switch, and removes the default tags that have been added to the cluster and
    each node.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `default_tags_enable` - Whether the default tags switch is enabled.

## Timeouts

This resource provides the following Timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The resource can be imported using the `cluster_id`, e.g.

```bash
terraform import huaweicloud_mapreduce_cluster_default_tags_switch.test <cluster_id>
```
