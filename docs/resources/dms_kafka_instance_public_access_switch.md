---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_instance_public_access_switch"
description: |-
  Use this resource to switch public access of a Kafka instance within HuaweiCloud.
---

# huaweicloud_dms_kafka_instance_public_access_switch

Use this resource to switch public access of a Kafka instance within HuaweiCloud.

-> This resource is a one-time action resource for switching public access of Kafka instance. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Enable public access with specified bandwidth

```hcl
variable "instance_id" {}
variable "publicip_id" {}

resource "huaweicloud_dms_kafka_instance_public_access_switch" "test" {
  instance_id = var.instance_id
  publicip_id = var.publicip_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the public access switch are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance.

* `public_boundwidth` - (Optional, Int, NonUpdatable) Specifies the public bandwidth of the Kafka instance.

* `eip_address` - (Optional, String, NonUpdatable) Specifies the elastic IP address of the Kafka instance.

* `publicip_id` - (Optional, String, NonUpdatable) Specifies the public IP ID of the Kafka instance.

-> At least one of `public_boundwidth`, `eip_address`, or `publicip_id` must be specified when the public access feature
   is currently disabled. If the public access feature is currently enabled, not specifying any of these parameters will
   disable the public access.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
