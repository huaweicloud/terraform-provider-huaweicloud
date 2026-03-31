---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_eip_associate"
description: |-
  Manages an EIP association for the DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_eip_associate

Use this resource to associate an EIP to a DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "eip_id" {}

resource "huaweicloud_dws_cluster_eip_associate" "test" {
  cluster_id = var.dws_cluster_id
  eip_id     = var.eip_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the cluster (to which the EIP associated) is
  located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster to which the EIP is associated.

* `eip_id` - (Required, String, NonUpdatable) Specifies the ID of the EIP to be associated with the DWS cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also `cluster_id`.

* `public_ip` - The public IP address of the cluster endpoint.

* `public_port` - The public port of the cluster endpoint.

## Import

The resource can be imported using the `id` (also `cluster_id`), e.g.

```bash
$ terraform import huaweicloud_dws_cluster_eip_associate.test <id>
```
