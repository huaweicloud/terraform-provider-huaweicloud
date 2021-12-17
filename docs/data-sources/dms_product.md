---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_product

Use this data source to get the ID of an available HuaweiCloud dms product. This is an alternative
to `huaweicloud_dms_product_v1`

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

* `region` - (Optional, String) The region in which to obtain the dms products. If omitted, the provider-level region
  will be used.

* `engine` - (Required, String) Indicates the name of a message engine. The valid values are __kafka__, __rabbitmq__.

* `instance_type` - (Required, String) Indicates an instance type. The valid values are __single__ and __cluster__.

* `version` - (Optional, String) Indicates the version of a message engine.

* `availability_zones` - (Optional, List) Indicates the list of availability zones with available resources.

* `vm_specification` - (Optional, String) Indicates VM specifications.

* `storage` - (Optional, String) Indicates the storage capacity of the resource.
  The default value is the storage capacity of the product.

* `bandwidth` - (Optional, String) Indicates the baseline bandwidth of a DMS instance.
  The valid values are __100MB__, __300MB__, __600MB__ and __1200MB__.

* `partition_num` - (Optional, String) Indicates the maximum number of topics that can be created for a Kafka instance.
  The valid values are __300__, __900__ and __1800__.

* `storage_spec_code` - (Optional, String) Indicates an I/O specification.
  The valid values are __dms.physical.storage.high__ and __dms.physical.storage.ultra__.

* `node_num` - (Optional, String) Indicates the number of nodes in a cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.
