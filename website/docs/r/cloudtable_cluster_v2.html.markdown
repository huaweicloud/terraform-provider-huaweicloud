---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cloudtable_cluster_v2"
sidebar_current: "docs-huaweicloud-resource-cloudtable-cluster-v2"
description: |-
  cloud table cluster management
---

# huaweicloud\_cloudtable\_cluster\_v2

cloud table cluster management

## Example Usage

### create a CloudTable cluster

```hcl
resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  name = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_cloudtable_cluster_v2" "cluster" {
  availability_zone = "{{ availability_zone }}"
  name = "terraform-test-cluster"
  rs_num = 2
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup.id}"
  subnet_id = "{{ network_id }}"
  vpc_id = "{{ vpc_id }}"
  storage_type = "COMMON"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` -
  (Required)
  Availability zone (AZ).  Changing this parameter will create a new resource.

* `name` -
  (Required)
  Cluster name. The value must be between 4 and 64 characters long and
  start with a letter. Only letters, digits, and hyphens (-) are
  allowed. It is case insensitive.  Changing this parameter will create a new resource.

* `rs_num` -
  (Required)
  Number of computing units. Value range: 2 to 10.  Changing this parameter will create a new resource.

* `security_group_id` -
  (Required)
  Security group ID.  Changing this parameter will create a new resource.

* `storage_type` -
  (Required)
  Storage I/O type. The value are ULTRAHIGH and COMMON.  Changing this parameter will create a new resource.

* `subnet_id` -
  (Required)
  Subnet ID.  Changing this parameter will create a new resource.

* `vpc_id` -
  (Required)
  VPC of the cluster.  Changing this parameter will create a new resource.

- - -

* `enable_iam_auth` -
  (Optional)
  Whether to enable IAM authentication for OpenTSDB.  Changing this parameter will create a new resource.

* `lemon_num` -
  (Optional)
  Number of Lemon nodes Value range: 2 to 10.  Changing this parameter will create a new resource.

* `opentsdb_num` -
  (Optional)
  Number of OpenTSDB nodes Value range: 2 to 10.  Changing this parameter will create a new resource.

* `tags` -
  (Optional)
  Enterprise project information. Structure is documented below. Changing this parameter will create a new resource.

The `tags` block supports:

* `key` -
  (Optional)
  Tag value.  Changing this parameter will create a new resource.

* `value` -
  (Optional)
  Tag key.  Changing this parameter will create a new resource.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created` -
  Time when the cluster was created.

* `hbase_public_endpoint` -
  HBase link of the public network.

* `lemon_link` -
  Lemon link of the intranet.

* `open_tsdb_link` -
  OpenTSDB link of the intranet.

* `opentsdb_public_endpoint` -
  OpenTSDB link of the public network.

* `storage_quota` -
  Storage quota.

* `used_storage_size` -
  Used storage space.

* `zookeeper_link` -
  ZooKeeper link of the intranet.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 30 minute.
