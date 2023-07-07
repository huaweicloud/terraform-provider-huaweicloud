---
subcategory: "Distributed Cache Service (DCS)"
---

# huaweicloud_dcs_instance_shards

Use this data source to get the list of DCS instance shards.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_instance_shards" "test" {
  instance_id  = var.instance_id
  replica_role = "slave"

  shard_names = [
    "group-1",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `shard_names` - (Optional, List) Specifies the list of the shard names.

* `replica_ips` - (Optional, List) Specifies the list of the replica ips.

* `replica_role` - (Optional, String) Specifies the role of the replica. Value options: **master**, **slave**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `shards` - Indicates the list of DCS instance replicas.
  The [Shard](#DcsInstanceShard_Shard) structure is documented below.

<a name="DcsInstanceShard_Shard"></a>
The `Shard` block supports:

* `shard_id` - Indicates the ID of the shard.

* `shard_name` - Indicates the name of the shard.

* `replicas` - Indicates the list of replicas in the shard.
  The [Replica](#DcsInstanceShard_ShardReplica) structure is documented below.

<a name="DcsInstanceShard_ShardReplica"></a>
The `ShardReplica` block supports:

* `id` - Indicates the ID of the replica.

* `ip` - Indicates the IP of the replica.

* `role` - Indicates the role of the replica.

* `node_id` - Indicates the ID of the node.

* `status` - Indicates the status of the replica.
