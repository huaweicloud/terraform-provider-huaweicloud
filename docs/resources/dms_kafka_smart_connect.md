---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_smart_connect"
description: ""
---

# huaweicloud_dms_kafka_smart_connect

Manage DMS kafka smart connect resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_kafka_smart_connect" "test" {
  instance_id       = var.instance_id
  storage_spec_code = "dms.physical.storage.ultra.v2"
  bandwidth         = "100MB"
  node_count        = 2
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the kafka instance.

  Changing this parameter will create a new resource.

* `storage_spec_code` - (Optional, String, ForceNew) Specifies the storage specification code of the connector.

  Changing this parameter will create a new resource.

* `bandwidth` - (Optional, String, ForceNew) Specifies the bandwidth of the connector.

  Changing this parameter will create a new resource.

* `node_count` - (Optional, Int, ForceNew) Specifies the node count of the connector. Defaults to 2 and minimum is 2.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID.

## Import

The kafka smart connect can be imported using the kafka `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_kafka_smart_connect.test <instance_id>/<id>
```
