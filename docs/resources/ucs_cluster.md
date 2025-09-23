---
subcategory: "Ubiquitous Cloud Native Service (UCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ucs_cluster"
description: ""
---

# huaweicloud_ucs_cluster

Manages a UCS cluster resource within HuaweiCloud.

## Example Usage

### Registering a K8S cluster

```hcl
variable "cluster_name" {}
variable "kubeconfig" {}

resource "huaweicloud_ucs_cluster" "test" {
  category     = "attachedcluster"
  cluster_type = "privatek8s"
  cluster_name = var.cluster_name
  country      = "CN"
  city         = "110000"

  annotations = {
    "kubeconfig" = var.kubeconfig
  }
}
```

### Registering a CCE cluster

```hcl
variable "fleet_id" {}
variable "cluster_id" {}
variable "cluster_region" {}
variable "cluster_project_id" {}

resource "huaweicloud_ucs_cluster" "test" {
  category           = "self"
  cluster_type       = "cce"
  fleet_id           = var.fleet_id
  cluster_id         = var.cluster_id
  cluster_region     = var.cluster_region
  cluster_project_id = var.cluster_project_id
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Required, String, ForceNew) Specifies the category of the cloud.

  Changing this parameter will create a new resource.

* `cluster_type` - (Required, String, ForceNew) Specifies the cluster type.

  Changing this parameter will create a new resource.

* `fleet_id` - (Optional, String) Specifies ID of the fleet to add the cluster into.
  If left empty, that means registering a cluster **discrete** cluster.

* `cluster_id` - (Optional, String, ForceNew) Specifies the cluster id.
  This can only be used when registering a cluster imported from CCE.

  Changing this parameter will create a new resource.

* `cluster_region` - (Optional, String, ForceNew) Specifies the cluster region.
   This can only be used when registering a cluster imported from CCE.

  Changing this parameter will create a new resource.

* `cluster_project_id` - (Optional, String, ForceNew) Specifies the cluster project ID.
   This can only be used when registering a cluster imported from CCE.

  Changing this parameter will create a new resource.

* `cluster_name` - (Optional, String, ForceNew) Specifies the name of the cluster to register.

  Changing this parameter will create a new resource.

* `cluster_labels` - (Optional, Map, ForceNew) Specifies the labels of the cluster to register.

  Changing this parameter will create a new resource.

* `annotations` - (Optional, Map, ForceNew) Specifies the annotations of the cluster to register.

  Changing this parameter will create a new resource.

* `service_provider` - (Optional, String, ForceNew) Specifies the cloud service provider.
  The value can be: **aws**, **azure**, **aliyun**, **googlecloud**,
  **tencentcloud**, **openshift**, **huaweicloud** and **privatek8s**.

  Changing this parameter will create a new resource.

* `country` - (Optional, String) Specifies the country name.

* `city` - (Optional, String) Specifies the city code.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `manage_type` - The cluster manage type. The value can be **grouped** and **discrete**.

## Import

The UCS cluster can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ucs_cluster.test b84c0d09-26cc-11ee-b6b2-0255ac100263
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `annotations`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_ucs_cluster" "test" {
    ...

  lifecycle {
    ignore_changes = [
      annotations
    ]
  }
}
```
