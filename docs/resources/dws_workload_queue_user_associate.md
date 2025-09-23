---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_queue_user_associate"
description: |-
  Use this resource to bind the users to the workload queue within HuaweiCloud.
---

# huaweicloud_dws_workload_queue_user_associate

Use this resource to bind the users to the workload queue within HuaweiCloud.

-> A user can only be associated with one workload queue.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "workload_queue_name" {}
variable "associated_user_names" {
  type = list(string)
}

resource "huaweicloud_dws_workload_queue_user_associate" "test" {
  cluster_id = var.dws_cluster_id
  queue_name = var.workload_queue_name
  user_names = var.associated_user_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the DWS cluster ID.
  Changing this creates a new resource.

* `queue_name` - (Required, String, ForceNew) Specifies the workload queue name to associate with the users.
  Changing this creates a new resource.

* `user_names` - (Required, List) Specifies the user names bound to the workload queue.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consists of `cluster_id` and `queue_name`, separated by a slash.

## Import

The resource can be imported using `cluster_id` and `queue_name` (also `id`), separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dws_workload_queue_user_associate.test <cluster_id>/<queue_name>
```
