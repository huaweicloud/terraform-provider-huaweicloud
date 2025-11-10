---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_volume_auto_expand_configuration"
description: |-
  Use this resource to configure the volume auto-expansion for the specified Kafka instance within HuaweiCloud.
---

# huaweicloud_dms_kafka_volume_auto_expand_configuration

Use this resource to configure the volume auto-expansion for the specified Kafka instance within HuaweiCloud.

-> 1. This resource is unavailable for single-node instance Kafka.
   <br/> 2. This resource is only a one-time action resource for configuring volume auto-expansion. Deleting this
   resource will not clear the configuration, but will only remove the resource information from the tfstate file.

## Example Usage

### Enable volume auto-expansion

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_kafka_volume_auto_expand_configuration" "test" {
  instance_id               = var.instance_id
  auto_volume_expand_enable = true
  expand_threshold          = 80
  expand_increment          = 10
  max_volume_size           = 400
}
```

### Disable volume auto-expansion

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_kafka_volume_auto_expand_configuration" "test" {
  instance_id               = var.instance_id
  auto_volume_expand_enable = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the Kafka instance to be configured volume
  auto-expansion is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance.

* `auto_volume_expand_enable` - (Optional, Bool, NonUpdatable) Specifies whether to enable disk auto-expansion.  
   Defaults to **false**.

* `expand_threshold` - (Optional, Int, NonUpdatable) Specifies the threshold for triggering disk auto-expansion,
  in percentage (%).  
  The valid value ranges from `20` to `80`.  
  This parameter is **required** when `auto_volume_expand_enable` is set to **true**.

* `expand_increment` - (Optional, Int, NonUpdatable) Specifies the percentage of storage space to be expanded
  out of the total instance storage space, in percentage (%).  
  The valid value ranges from `10` to `100`.
  This parameter is **required** when `auto_volume_expand_enable` is set to **true**.

* `max_volume_size` - (Optional, Int, NonUpdatable) Specifies the maximum volume size for disk auto-expansion, in GB.  
  This parameter is **required** when `auto_volume_expand_enable` is set to **true**.  
  The value must meet the following requirements:
  + The value must be divisible by `100`.
  + The value must be greater than the current instance disk capacity.
  + The value must be less than the number of nodes multiplied by `30000`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
