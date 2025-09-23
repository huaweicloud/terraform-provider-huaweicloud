---
subcategory: "CloudTable"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cloudtable_cluster"
description: ""
---

# huaweicloud_cloudtable_cluster

Manages a CloudTable cluster resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_name" {}
variable "vpc_id" {}
variable "network_id" {}
variable "security_group_id" {}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = var.vpc_id
  ...
}

resource "huaweicloud_cloudtable_cluster" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = var.cluster_name
  storage_type      = "ULTRAHIGH"
  vpc_id            = var.vpc_id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = var.security_group_id
  hbase_version     = "1.0.6"
  rs_num            = 4
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the cluster.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone in which to create the cluster.
  Please following [reference](https://developer.huaweicloud.com/en-us/endpoint/?CloudTable) for the values.
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the cluster name.  
  The name consists of `4` to `64` characters, including lowercase letters, numbers and hyphens (-).  
  Changing this parameter will create a new resource.

* `storage_type` - (Required, String, ForceNew) Specifies the storage type.
  The valid values are **COMMON** and **ULTRAHIGH**. Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the cluster belongs.
  Changing this parameter will create a new resource.

* `network_id` - (Required, String, ForceNew) Specifies the ID of the network to which the cluster belongs.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String, ForceNew) Specifies the security group ID of the cluster.
  Changing this parameter will create a new resource.

* `hbase_version` - (Required, String, ForceNew) Specifies the version of HBase datastore.

* `iam_auth_enabled` - (Optional, Bool, ForceNew) Specifies whether IAM authorization is enabled.
  Changing this parameter will create a new resource.

* `opentsdb_num` - (Optional, Int, ForceNew) Specifies the TSD nodes number of the cluster.
  Changing this parameter will create a new resource.

* `rs_num` - (Optional, Int, ForceNew) Specifies the compute nodes number of the cluster.
  The valid values must be at least `2`. Defaults to `2`.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The cluster ID.

* `status` - The cluster status.

* `created_at` - The time (UTC) when the cluster was created.

* `hbase_public_endpoint` - The HBase public network endpoint address.

* `open_tsdb_link` - The intranet OpenTSDB connection access address.

* `opentsdb_public_endpoint` - The OpenTSDB public network endpoint address.

* `storage_size` - The storage size, in GB.

* `storage_size_used` - The currently used storage, in GB.

* `zookeeper_link` - The intranet zookeeper connection access address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 5 minutes.

## Import

Clusters can be imported by their `id`. e.g.:

```bash
terraform import huaweicloud_cloudtable_cluster.test 4c2d38b6-6fb0-480c-8813-5f536b5ba6a4
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `availability_zone`, `network_id`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition should be updated to
align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cloudtable_cluster" "test" {
    ...

  lifecycle {
    ignore_changes = [
      availability_zone, network_id,
    ]
  }
}
```
