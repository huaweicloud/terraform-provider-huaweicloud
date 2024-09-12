---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_cns"
description: |-
  Use this data source to query the list of CNs under specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_cns

Use this data source to query the list of CNs under specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}

data "huaweicloud_dws_cluster_cns" "test" {
  cluster_id = var.dws_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID to which the CNs belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `min_num` - The minimum number of CNs supported by the cluster.

* `max_num` - The maximum number of CNs supported by the cluster.

* `cns` - The list of the CNs under specified DWS cluster.

  The [cns](#cns_struct) structure is documented below.

<a name="cns_struct"></a>
The `cns` block supports:

* `id` - The ID of the CN.

* `name` - The name of the CN.

* `private_ip` - The private IP address of the CN.
