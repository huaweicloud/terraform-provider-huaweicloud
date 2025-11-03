---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_cluster_authorize_access_key"
description: |-
  Manages an access key authorization resource for a CPCS cluster within HuaweiCloud.
---

# huaweicloud_cpcs_cluster_authorize_access_key

Manages an access key authorization resource for a CPCS cluster within HuaweiCloud.

-> Currently, this resource is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "cluster_id" {}
variable "app_id" {}
variable "access_key_id" {}

resource "huaweicloud_cpcs_cluster_authorize_access_key" "test" {
  cluster_id    = var.cluster_id
  app_id        = var.app_id
  access_key_id = var.access_key_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the CPCS cluster.

* `app_id` - (Required, String, NonUpdatable) Specifies the ID of the application.
  It must be the value of the `app_id` that is already associated with the cluster.

* `access_key_id` - (Required, String, NonUpdatable) Specifies the ID of the access key to be authorized.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (Randomly generated UUID).

* `status` - The status of the access key.

* `app_name` - The name of the application.

* `access_key` - The access key (sensitive, will be marked as sensitive in the state).

* `key_name` - The name of the access key.

* `create_time` - The creation time of the access key, UNIX timestamp in milliseconds.

## Import

The CPCS cluster access key authorization resource can be imported using the `cluster_id` and `access_key_id`,
separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cpcs_cluster_authorize_access_key.test <cluster_id>/<access_key_id>
```
