---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud_cce_pvc

Manages a CCE Persistent Volume Claim resource within HuaweiCloud.

## Example Usage

### Create PVC by importing EVS volume

```hcl
variable "cluster_id" {}
variable "namespace" {}
variable "pvc_name" {}
variable "volume_id" {}

resource "huaweicloud_cce_pvc" "test" {
  namespace   = var.namespace
  name        = var.pvc_name
  volume_type = "bs"
  volume_id   = var.volume_id
}
```

### Create PVC by importing OBS bucket

```hcl
variable "cluster_id" {}
variable "namespace" {}
variable "pvc_name" {}
variable "obs_bucket_id" {}

resource "huaweicloud_cce_pvc" "test" {
  clsuter_id  = var.cluster_id
  namespace   = var.namespace
  name        = var.pvc_name
  volume_type = "obs"
  volume_id   = var.obs_bucket_id
}
```

### Create PVC by importing SFS file system

```hcl
variable "cluster_id" {}
variable "namespace" {}
variable "pvc_name" {}
variable "sfs_file_system_id" {}

resource "huaweicloud_cce_pvc" "test" {
  clsuter_id  = var.cluster_id
  namespace   = var.namespace
  name        = var.pvc_name
  volume_type = "nfs"
  volume_id   = var.sfs_file_system_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the PVC resource.
  If omitted, the provider-level region will be used. Changing this will create a new PVC resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the cluster ID to which the CCE PVC belongs.

* `namespace` - (Required, String, ForceNew) Specifies the namespace to logically divide your containers into different
  group. Changing this will create a new PVC resource.

* `name` - (Required, String, ForceNew) Specifies the unique name of the PVC resource. This parameter can contain a
  maximum of 63 characters, which may consist of lowercase letters, digits and hyphens (-), and must start and end with
  lowercase letters and digits. Changing this will create a new PVC resource.

* `volume_id` - (Required, String, ForceNew) Specifies the ID of the storage bound to the CCE namespace. Changing this
  will create a new PVC resource.

* `volume_type` - (Optional, String, ForceNew) Specifies the type of the storage bound to the CCE namespace.
  The valid values are as follows:
  + **bs**: EVS disk.
  + **obs**: OBS bucket.
  + **nfs**: SFS file system.

  default to **bs**.

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone in which to create the PVC.
  Please following [reference](https://developer.huaweicloud.com/intl/en-us/endpoint?CCE) for the values.
  Changing this creates a new PVC resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The PVC ID in UUID format.

* `creation_timestamp` - The server time when PVC was created.

* `status` - The current phase of the PVC.
  + **Pending**: Not yet bound.
  + **Bound**: Already bound.

* `access_modes` - Access mode of the volume (EVS disk, OBS bucket or SFS file system).
  + **ReadWriteOnce**: The volume can be mounted as read-write by a single node.
  + **ReadOnlyMany**: The volume can be mounted as read-only by many nodes.
  + **ReadWriteMany**: The volume can be mounted as read-write by many nodes.

## Import

PVCs can be imported using their `id`, `namespace` and the `cluster_id` to which the pvc and namespace belongs,
separated by a slash, e.g.

```
$ terraform import huaweicloud_cce_pvc.test <cluster_id>/<namespace>/<id>
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.
* `delete` - Default is 3 minute.
