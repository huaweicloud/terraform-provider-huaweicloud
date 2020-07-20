---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_cassandra_instance"
sidebar_current: "docs-huaweicloud-resource-gaussdb-cassandra-instance"
description: |-
  GaussDB for Cassandra instance management
---

# huaweicloud\_gaussdb\_cassandra\_instance

GaussDB for Cassandra instance management within HuaweiCoud.

## Example Usage

### create a gaussdb for cassandra instance with tags

```hcl
resource "huaweicloud_gaussdb_cassandra_instance" "instance_1" {
  name        = "gaussdb_cassandra_instance_1"
  password    = var.password
  flavor      = "geminidb.cassandra.xlarge.4"
  volume_size = 100
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### create a gaussdb cassandra instance with backup strategy

```hcl
resource "huaweicloud_gaussdb_cassandra_instance" "instance_1" {
  name        = "gaussdb_cassandra_instance_1"
  password    = var.password
  flavor      = "geminidb.cassandra.xlarge.4"
  volume_size = 100
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Specifies the AZ name.
  Changing this parameter will create a new resource.

* `name` - (Required) Specifies the instance name, which can be the same
  as an existing instance name. The value must be 4 to 64 characters in
  length and start with a letter. It is case-sensitive and can contain
  only letters, digits, hyphens (-), and underscores (_).
  Changing this parameter will create a new resource.

* `flavor` - (Required) Specifies the instance specifications. For details, 
  see [DB Instance Specifications](https://support.huaweicloud.com/intl/en-us/productdesc-geminidb/geminidb_01_0006.html)
  Changing this parameter will create a new resource.

* `node_num` - (Optional) Specifies the number of nodes, ranges from 3 to 12. Defaults to 3.
  Changing this parameter will create a new resource.

* `volume_size` - (Required) Specifies the storage space in GB. The value must be a multiple of 10.
  For a GaussDB Cassandra DB instance, the minimum storage space is 100 GB, and the maximum
  storage space is related to the instance performance specifications. For details, 
  see [DB Instance Specifications](https://support.huaweicloud.com/intl/en-us/productdesc-geminidb/geminidb_01_0006.html)
  Changing this parameter will create a new resource.

* `password` - (Required) Specifies the database password. The value must be 8 to 32 characters
  in length, including uppercase and lowercase letters, digits, and special characters,
  such as ~!@#%^*-_=+?
  Changing this parameter will create a new resource.

  You are advised to enter a strong password to improve security, preventing security risks
  such as brute force cracking.
  Changing this parameter will create a new resource.

* `vpc_id` -  (Required) Specifies the VPC ID.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required) Specifies the network ID of a subnet.
  Changing this parameter will create a new resource.

* `security_group_id` - (Optional) Specifies the security group ID. Required if the selected subnet doesn't enable network ACL.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional) Specifies the enterprise project id, Only valid for users who
  have enabled the enterprise multi-project service.
  Changing this parameter will create a new resource.

* `ssl` - (Optional) Specifies whether to enable or disable SSL. Defaults to false.
  Changing this parameter will create a new resource.

* `datastore` - (Optional) Specifies the database information. Structure is documented below.
  Changing this parameter will create a new resource.

* `backup_strategy` - (Optional) Specifies the advanced backup policy. Structure is documented below.
  Changing this parameter will create a new resource.

* `tags` - (Optional) The key/value pairs to associate with the instance.

The `datastore` block supports:

* `engine` - (Optional) Specifies the database engine. Only "GeminiDB-Cassandra" is supported now.

* `version` - (Optional) Specifies the database version. Only "3.11" is supported now.

* `storage_engine` - (Optional) Specifies the storage engine. Only "rocksDB" is supported now.


The `backup_strategy` block supports:

* `start_time` - (Required) Specifies the backup time window. Automated backups
  will be triggered during the backup time window. It must be a valid value in
  the "hh:mm-HH:MM" format. The current time is in the UTC format.
  The HH value must be 1 greater than the hh value. The values of mm and MM
  must be the same and must be set to 00. Example value: 08:00-09:00, 03:00-04:00.

* `keep_days` - (Optional) Specifies the number of days to retain the generated
   backup files. The value ranges from 0 to 35.
   If this parameter is set to 0, the automated backup policy is not set.
   If this parameter is not transferred, the automated backup policy is enabled by default.
   Backup files are stored for seven days by default.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `status` - Indicates the DB instance status.
* `region` - Indicates the region where the DB instance is deployed.
* `port` - Indicates the database port.
* `mode` - Indicates the instance type.
* `db_user_name` - Indicates the default username.
* `nodes` - Indicates the instance nodes information. Structure is documented below.
- `private_ips` - Indicates the IP address list of the db.

The `nodes` block contains:

- `id` - Indicates the node ID.
- `name` - Indicates the node name.
- `status` - Indicates the node status.
- `private_ip` - Indicates the private IP address of a node.

## Import

GaussDB Cassandra instance can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_gaussdb_cassandra_instance.instance_1 2e045d8b-b226-4aa2-91b9-7e76357655c06
```
