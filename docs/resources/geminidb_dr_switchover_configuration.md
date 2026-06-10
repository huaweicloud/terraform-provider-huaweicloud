---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_dr_switchover_configuration"
description: |-
  Manages disaster recovery switchover settings for a GeminiDB Cassandra instance within HuaweiCloud.
---

# huaweicloud_geminidb_dr_switchover_configuration

Manages disaster recovery switchover settings for a GeminiDB Cassandra instance within HuaweiCloud.

## Example Usage

```hcl
var "instance_id" {}

resource "huaweicloud_geminidb_dr_switchover_configuration" "test" {
  instance_id      = var.instance_id
  switchover_ratio = 60
  sync_delay       = 300
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new geminidb instance resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GeminiDB Redis instance.

* `switchover_ratio` - (Optional, Int) Specifies the percentage of faulty nodes for disaster recovery switchover.
* The valid values are **50**, **60**, **70**, **80**, **90**, and **100**. Default value is **100**.

* `sync_delay` - (Optional, Int) Specifies the data synchronization delay threshold in seconds. When the synchronization
* delay between the disaster recovery instance and the primary instance exceeds this value, disaster recovery switchover
* will not be triggered. The minimum value is **10** seconds. If not specified, no delay judgment will be performed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as the instance ID.

## Import

The GeminiDB configuring DR switchover for an Instance can be imported using the `instance_id`, e.g.

```
$ terraform import huaweicloud_geminidb_dr_switchover_configuration.test <instance_id>
```
