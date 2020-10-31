---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud\_dms\_product

Use this data source to get the ID of an available HuaweiCloud dms product.
This is an alternative to `huaweicloud_dms_product_v1`

## Example Usage

```hcl

data "huaweicloud_dms_product" "product1" {
  engine            = "kafka"
  version           = "1.1.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the dms products. If omitted, the provider-level region will work as default.

* `engine` - (Required) Indicates the name of a message engine.

* `version` - (Optional) Indicates the version of a message engine.

* `instance_type` - (Required) Indicates an instance type. Options: "single" and "cluster"

* `vm_specification` - (Optional) Indicates VM specifications.

* `storage` - (Optional) Indicates the message storage space.

* `bandwidth` - (Optional) Indicates the baseline bandwidth of a Kafka instance.

* `partition_num` - (Optional) Indicates the maximum number of topics that can be created for a Kafka instance.

* `storage_spec_code` - (Optional) Indicates an I/O specification.

* `io_type` - (Optional) Indicates an I/O type.

* `node_num` - (Optional) Indicates the number of nodes in a cluster.


## Attributes Reference

`id` is set to the ID of the found product. In addition, the following attributes
are exported:

* `engine` - See Argument Reference above.
* `version` - See Argument Reference above.
* `instance_type` - See Argument Reference above.
* `vm_specification` - See Argument Reference above.
* `bandwidth` - See Argument Reference above.
* `partition_num` - See Argument Reference above.
* `storage_spec_code` - See Argument Reference above.
* `io_type` - See Argument Reference above.
* `node_num` - See Argument Reference above.
