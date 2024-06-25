---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_connectivity"
description: |-
  Manages CSS logstash connectivity resource within HuaweiCloud.
---

# huaweicloud_css_logstash_connectivity

Manages CSS logstash connectivity resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_css_logstash_connectivity" "test" {
  cluster_id = var.cluster_id

  address_and_ports {
    address = "192.168.0.11"
    port    = "9600"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the CSS logstash cluster.

* `address_and_ports` - (Required, List, NonUpdatable) Specifies the list of addresses and ports.
  The [address_and_ports](#css_logstash_address_and_ports) structure is documented below.

<a name="css_logstash_address_and_ports"></a>
The `address_and_ports` block supports:

* `address` - (Required, String, NonUpdatable) Specifies the ip address.

* `port` - (Required, Int, NonUpdatable) Specifies the port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `connectivity_results` - The connectivity results.
  The [connectivity_results](#css_logstash_connectivity_results) structure is documented below.

<a name="css_logstash_connectivity_results"></a>
The `connectivity_results` block supports:

* `address` - The ip address.

* `port` - The port.

* `status` - The connectivity test result.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
