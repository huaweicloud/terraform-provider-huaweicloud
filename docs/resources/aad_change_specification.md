---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_change_specification"
description: |-
  Use this resource to modify Advanced Anti-DDos specification within HuaweiCloud.
---

# huaweicloud_aad_change_specification

Use this resource to modify Advanced Anti-DDos specification within HuaweiCloud.

-> This resource is only a one-time action resource for updating AAD instance specification. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "basic_bandwidth" {}
variable "elastic_bandwidth" {}
variable "service_bandwidth" {}
variable "port_num" {}
variable "bind_domain_num" {}

resource "huaweicloud_aad_change_specification" "test" {
  instance_id = var.instance_id
  
  upgrade_data {
    basic_bandwidth   = var.basic_bandwidth
    elastic_bandwidth = var.elastic_bandwidth
    service_bandwidth = var.service_bandwidth
    port_num          = var.port_num
    bind_domain_num   = var.bind_domain_num
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, NonUpdatable) Specifies the AAD instance ID.

* `upgrade_data` - (Required, List, NonUpdatable) Specifies the upgrade data.

  The [upgrade_data](#upgrade_data_struct) structure is documented below.

<a name="upgrade_data_struct"></a>
The `upgrade_data` block supports:

* `basic_bandwidth` - (Optional, String, NonUpdatable) Specifies the basic bandwidth (Gbps).

* `elastic_bandwidth` - (Optional, String, NonUpdatable) Specifies the elastic bandwidth (Gbps).

* `service_bandwidth` - (Optional, Int, NonUpdatable) Specifies the service bandwidth (Mbps).

* `port_num` - (Optional, Int, NonUpdatable) Specifies the port number.

* `bind_domain_num` - (Optional, Int, NonUpdatable) Specifies the bind domain number.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `instance_id`).

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Defaults to `10` minutes.
