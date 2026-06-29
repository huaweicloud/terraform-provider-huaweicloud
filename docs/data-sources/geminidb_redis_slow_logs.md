---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_redis_slow_logs"
description: |-
  Use this data source to query the list of GeminiDB Redis instance slow logs.
---

# huaweicloud_geminidb_redis_slow_logs

Use this data source to query the list of GeminiDB Redis instance slow logs.

## Example Usage

### Query all slow logs

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_geminidb_redis_slow_logs" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

### Query slow logs by node ID

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}
variable "node_id" {}

data "huaweicloud_geminidb_redis_slow_logs" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the GeminiDB Redis instance ID.

* `start_time` - (Required, String) Specifies the query start time.
  The format is **yyyy-mm-ddThh:mm:ss+0800**, e.g. **2026-06-01T12:00:00+0800**.
  The start time cannot be `30` days earlier than the current time.

* `end_time` - (Required, String) Specifies the query end time.
  The format is **yyyy-mm-ddThh:mm:ss+0800**, e.g. **2026-06-01T12:00:00+0800**.
  The end time cannot be later than the current time.

* `operate_type` - (Optional, String) Specifies the statement type.
  The value can be **set**, **get**, **del**, **incr**, **incrby**, **incrbyfloat**, **decr**, **decrby**, **getset**,
  **append**, **mget**, **keys**, **setnx**, **setex**, **psetex**, **delvx**, **mset**, **msetnx**, **getrange**,
  **substr**, **setrange**, **strlen**, **exists**, **expire**, **pexpire**, **expireat**, **pexpireat**, **ttl**,
  **pttl**, **persist**, **type**, **scanx**, **pksetexat**, **sort**, **hdel**, **hset**, **hget**, **hgetall**,
  **hexists**, **hincrby**, **hincrbyfloat**, **hkeys**, **hlen**, **hmget**, **hmset**, **hsetnx**, **hstrlen**,
  **hvals**, **hscan**, **hscanx**, **pkhscanrange**, **pkhrscanrange**, **lindex**, **linsert**, **llen**, **lpop**,
  **lpush**, **lpushx**, **lrange**, **lrem**, **lset**, **ltrim**, **rpop**, **rpoplpush**, **rpush**, **rpushx**,
  **zadd**, **zcard**,**zscan**,**zincrby**, **zrange**, **zrevrange**, **zrangebyscore**, **zrevrangebyscore**,
  **zcount**, **zrem**, **zunionstore**, **zinterstore**, **zrank**, **zrevrank**, **zscore**, **zrangebylex**,
  **zrevrangebylex**, **zlexcount**, **zremrangebyrank**, **zremrangebyscore**, **zremrangebylex**, **zpopmax**,
  **zpopmin**, **sadd**, **spop**, **scard**, **smembers**, **sscan**, **srem**, **sunion**, **sunionstore**,
  **sinter**, **sinterstore**, **sismember**, **sdiff**, **sdiffstore**, **smove**, **srandmember**, **bitset**,
  **bitget**, **bitcount**, **bitpos**, **bitop**, **bitfield**, **pfadd**, **pfcount**, **pfmerge**, **geoadd**,
  **georadiusbymember**, **georadius**, **geohash**, **geodist**, **geopos**, **xadd**, **xack**, **xgroup**, **xdel**,
  **xtrim**, **xlen**, **xrange**, **xrevrange**, **xclaim**, **xpending**, **xinfo**, **xread**, or **xreadgroup**.

* `node_id` - (Optional, String) Specifies the node ID.

* `keywords` - (Optional, List) Specifies query the slow logs by keywords matched.
  A maximum of `10` keywords are supported.

* `max_cost_time` - (Optional, Int) Specifies the logs can be searched based on the maximum execution duration, in ms.

* `min_cost_time` - (Optional, Int) Specifies the logs can be searched based on the minimum execution duration, in ms.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slow_logs` - The list of slow logs.
  The [slow_logs](#slow_logs_struct) structure is documented below.

<a name="slow_logs_struct"></a>
The `slow_logs` block supports:

* `node_id` - The node ID.

* `node_name` - The node name.

* `whole_message` - The statement.

* `operate_type` - The statement type.

* `cost_time` - The execution time, in ms.

* `log_time` - The UTC time when a log is generated.

* `line_num` - The sequence number of a log event.
