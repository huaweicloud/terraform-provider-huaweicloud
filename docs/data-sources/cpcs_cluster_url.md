---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_cluster_url"
description: |-
  Use this data source to get the CPCS cluster URL.
---

# huaweicloud_cpcs_cluster_url

Use this data source to get the CPCS cluster URL.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cpcs_cluster_url" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the CPCS cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `uri` - The cluster service management redirect link.
