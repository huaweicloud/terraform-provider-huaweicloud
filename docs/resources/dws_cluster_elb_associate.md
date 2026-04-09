---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_elb_associate"
description: |-
  Use this resource to associate an ELB to a DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_elb_associate

Use this resource to associate an ELB to a DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "elb_id" {}

resource "huaweicloud_dws_cluster_elb_associate" "test" {
  cluster_id = var.dws_cluster_id
  elb_id     = var.elb_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the cluster (to which the ELB associated) is
  located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster to which the ELB is associated.

* `elb_id` - (Required, String, NonUpdatable) Specifies the ID of the ELB to be associated with the DWS cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also `cluster_id`.

* `public_ip` - The public IP address of the ELB loadbalancer.

* `private_ip` - The private IP address of the ELB loadbalancer.

## Import

The resource can be imported using the `id` (also `cluster_id`), e.g.

```bash
$ terraform import huaweicloud_dws_cluster_elb_associate.test <id>
```
