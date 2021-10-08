---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud_cce_pvc

Manages a CCE Persistent Volume Claim resource within HuaweiCloud.

## Example Usage

### Create PVC with EVS

```hcl
variable "cluster_id" {}
variable "namespace" {}
variable "pvc_name" {}

resource "huaweicloud_cce_pvc" "test" {
  cluster_id  = var.cluster_id
  namespace   = var.namespace
  name        = var.pvc_name
  annotations = {
    "everest.io/disk-volume-type" = "SSD"
  }
  storage_class_name = "csi-disk"
  access_modes = ["ReadWriteOnce"]
  storage = "10Gi"
}
```

### Create PVC with OBS

```hcl
variable "cluster_id" {}
variable "namespace" {}
variable "pvc_name" {}

resource "huaweicloud_cce_pvc" "test" {
  cluster_id  = var.cluster_id
  namespace   = var.namespace
  name        = var.pvc_name
  annotations = {
    "everest.io/obs-volume-type" = "STANDARD"
    "csi.storage.k8s.io/fstype" =  "obsfs"
  }
  storage_class_name = "csi-obs"
  access_modes = ["ReadWriteMany"]
  storage = "1Gi"
}
```

### Create PVC with SFS

```hcl
variable "cluster_id" {}
variable "namespace" {}
variable "pvc_name" {}

resource "huaweicloud_cce_pvc" "test" {
  cluster_id  = var.cluster_id
  namespace   = var.namespace
  name        = var.pvc_name
  storage_class_name = "csi-nas"
  access_modes = ["ReadWriteMany"]
  storage = "10Gi"
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

* `annotations` - (Optional, Map, ForceNew) An unstructured key value map for external parameters. Changing this
  will create a new PVC resource.

* `labels` - (Optional, Map, ForceNew) Map of string keys and values for labels. Changing this
  will create a new PVC resource.

* `storage_class_name` - (Required, String, ForceNew) Specifies the type of the storage bound to the CCE pvc.
  The valid values are as follows:
  + **csi-disk**: EVS.
  + **csi-obs**: OBS.
  + **csi-nas**: SFS.
  + **csi-sfsturbo**: SFS-Turbo.

* `access_modes` - (Required, List, ForceNew) Specifies the desired access modes the volume should have.
  The valid values are as follows:
  + **ReadWriteOnce**: The volume can be mounted as read-write by a single node.
  + **ReadOnlyMany**: The volume can be mounted as read-only by many nodes.
  + **ReadWriteMany**: The volume can be mounted as read-write by many nodes.

* `storage` - (Required, String, ForceNew) Specifies the minimum amount of storage resources required.
  Changing this creates a new PVC resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The PVC ID in UUID format.

* `creation_timestamp` - The server time when PVC was created.

* `status` - The current phase of the PVC.
  + **Pending**: Not yet bound.
  + **Bound**: Already bound.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.
* `delete` - Default is 3 minute.
