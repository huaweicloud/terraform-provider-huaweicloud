---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_ingress_port"
description: |-
  Manages an APIG instance custom ingress port resource within HuaweiCloud.
---

# huaweicloud_apig_instance_ingress_port

Manages an APIG instance custom ingress port resource within HuaweiCloud.

## Example Usage

### Create a custom ingress port with HTTP protocol

```hcl
variable "instance_id" {}

resource "huaweicloud_apig_instance_ingress_port" "test" {
  instance_id = var.instance_id
  protocol    = "HTTP"
  port        = 8080
}
```

### Create a custom ingress port with HTTPS protocol

```hcl
variable "instance_id" {}

resource "huaweicloud_apig_instance_ingress_port" "test" {
  instance_id = var.instance_id
  protocol    = "HTTPS"
  port        = 8443
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the ingress port is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the custom
  ingress port belongs.

* `protocol` - (Required, String, NonUpdatable) Specifies the protocol of the custom ingress port.  
  The valid values are as follows:
  + **HTTP**: The custom ingress port uses HTTP protocol.
  + **HTTPS**: The custom ingress port uses HTTPS protocol.

* `port` - (Required, Int, NonUpdatable) Specifies the port number of the custom ingress port.  
  The valid value is range from `1,024` to `49,151`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the ingress port ID.

* `status` - The status of the custom ingress port.
  + **normal**: The custom ingress port status is normal.
  + **abnormal**: The custom ingress port status is abnormal and cannot be used.

## Import

The resource can be imported using the `instance_id` and `id` (also `ingress_port_id`), separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_apig_instance_ingress_port.test <instance_id>/<id>
```

If the value of `ingress_port_id` is unknown, it can be imported using `instance_id`, `protocol`, and `port`, separated
by slashes (/), e.g.

```bash
$ terraform import huaweicloud_apig_instance_ingress_port.test <instance_id>/<protocol>/<port>
```
