---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_pod_identity_association"
description: -|
  Manages a CCE cluster pod identity association resource within HuaweiCloud.
---

# huaweicloud_cce_cluster_pod_identity_association

Manages a CCE cluster pod identity association resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "namespace" {}
variable "service_account" {}
variable "agency_name" {}

resource "huaweicloud_cce_cluster_pod_identity_association" "test" {
  cluster_id      = var.cluster_id
  namespace       = var. namespace
  service_account = var.service_account
  agency_name     = var.agency_name

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the pod identity association.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `cluster_id` - (Required, String, NoneUpdatable) Specifies the cluster ID that the pod identity association belongs.

* `namespace` - (Required, String, NoneUpdatable) Specifies the namespace of the service account associated with
  a pod identity association.

* `service_account` - (Required, String, NoneUpdatable) Specifies the name of the service account associated with
  a pod identity association. Only one pod identity association can be created for a service account.

* `agency_name` - (Required, String) Specifies the name of the agency to be associated with a pod identity association.
  The agency can be a general agency or a trust agency.

* `tags` - (Optional, Map, NoneUpdatable) Specifies the key/value pairs to associate with the pod identity association.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the uuid of the pod identity association.

* `created_at` - The creation time of the pod identity association.

* `updated_at` - The last update time of the pod identity association.

## Import

The CCE pod identity association can be imported using the `cluster_id` and `id`, separated by a slash (`/`), e.g.

```bash
$ terraform import huaweicloud_cce_cluster_pod_identity_association.test <cluster_id>/<id>
```
