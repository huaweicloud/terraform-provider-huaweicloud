---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_cluster_synchronized_iam_objects"
description: |-
  Use this data source to query synchronized IAM users and user groups of an MRS cluster within HuaweiCloud.
---

# huaweicloud_mapreduce_cluster_synchronized_iam_objects

Use this data source to query synchronized IAM users and user groups of an MRS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_mapreduce_cluster_synchronized_iam_objects" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the synchronized IAM users and user groups are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `user_names` - The synchronized IAM user name list.

* `group_names` - The synchronized IAM user group name list.
