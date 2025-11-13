---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_ingress_ports"
description: |-
  Use this data source to query the custom ingress ports of the dedicated instance within HuaweiCloud.
---

# huaweicloud_apig_instance_ingress_ports

Use this data source to query the custom ingress ports of the dedicated instance within HuaweiCloud.

## Example Usage

### Query all custom ingress ports under a specified instance

```hcl
variable "instance_id" {}

data "huaweicloud_apig_instance_ingress_ports" "test" {
  instance_id = var.instance_id
}
```

### Query all custom ingress ports with HTTPS protocol under a specified instance

```hcl
variable "instance_id" {}

data "huaweicloud_apig_instance_ingress_ports" "test" {
  instance_id = var.instance_id
  protocol    = "HTTPS"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the ingress ports are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the ingress ports belong.

* `protocol` - (Optional, String) Specifies the protocol of the ingress port to be queried.  
  The valid values are as follows:
  + **HTTP**: The ingress port uses HTTP protocol.
  + **HTTPS**: The ingress port uses HTTPS protocol.

* `port` - (Optional, Int) Specifies the port number of the ingress port to be queried.  
  The valid value is range from `1,024` to `49,151`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ingress_ports` - The list of the ingress ports that matched filter parameters.  
  The [ingress_ports](#apig_ingress_ports_attr) structure is documented below.

<a name="apig_ingress_ports_attr"></a>
The `ingress_ports` block supports:

* `id` - The ID of the ingress port.

* `protocol` - The protocol of the ingress port.

* `port` - The port number of the ingress port.

* `status` - The status of the ingress port.
  + **normal**: The ingress port status is normal.
  + **abnormal**: The ingress port status is abnormal and cannot be used.
