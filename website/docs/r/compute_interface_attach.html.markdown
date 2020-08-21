---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_interface_attach"
sidebar_current: "docs-huaweicloud-resource-compute-interface-attach"
description: |-
  Attaches a Network Interface to an Instance.
---

# huaweicloud\_compute\_interface\_attach

Attaches a Network Interface to an Instance.
This is an alternative to `huaweicloud_compute_interface_attach_v2`

## Example Usage

### Basic Attachment

```hcl
resource "huaweicloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "huaweicloud_compute_interface_attach" "ai_1" {
  instance_id = "${huaweicloud_compute_instance.instance_1.id}"
  network_id  = "${huaweicloud_networking_network_v2.network_1.id}"
}

```

### Attachment Specifying a Fixed IP

```hcl
resource "huaweicloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "huaweicloud_compute_interface_attach" "ai_1" {
  instance_id = "${huaweicloud_compute_instance.instance_1.id}"
  network_id  = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip    = "10.0.10.10"
}

```


### Attachment Using an Existing Port

```hcl
resource "huaweicloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = "${huaweicloud_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}


resource "huaweicloud_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "huaweicloud_compute_interface_attach" "ai_1" {
  instance_id = "${huaweicloud_compute_instance.instance_1.id}"
  port_id     = "${huaweicloud_networking_port_v2.port_1.id}"
}

```

### Attaching Multiple Interfaces

```hcl
resource "huaweicloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_port_v2" "ports" {
  count          = 2
  name           = "${format("port-%02d", count.index + 1)}"
  network_id     = "${huaweicloud_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}

resource "huaweicloud_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "huaweicloud_compute_interface_attach" "attachments" {
  count          = 2
  instance_id = "${huaweicloud_compute_instance.instance_1.id}"
  port_id     = "${huaweicloud_networking_port_v2.ports.*.id[count.index]}"
}
```

Note that the above example will not guarantee that the ports are attached in
a deterministic manner. The ports will be attached in a seemingly random
order.

If you want to ensure that the ports are attached in a given order, create
explicit dependencies between the ports, such as:

```hcl
resource "huaweicloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_port_v2" "ports" {
  count          = 2
  name           = "${format("port-%02d", count.index + 1)}"
  network_id     = "${huaweicloud_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}

resource "huaweicloud_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "huaweicloud_compute_interface_attach" "ai_1" {
  instance_id = "${huaweicloud_compute_instance.instance_1.id}"
  port_id     = "${huaweicloud_networking_port_v2.ports.*.id[0]}"
}

resource "huaweicloud_compute_interface_attach" "ai_2" {
  instance_id = "${huaweicloud_compute_instance.instance_1.id}"
  port_id     = "${huaweicloud_networking_port_v2.ports.*.id[1]}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the interface attachment.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new attachment.

* `instance_id` - (Required) The ID of the Instance to attach the Port or Network to.

* `port_id` - (Optional) The ID of the Port to attach to an Instance.
   _NOTE_: This option and `network_id` are mutually exclusive.

* `network_id` - (Optional) The ID of the Network to attach to an Instance. A port will be created automatically.
   _NOTE_: This option and `port_id` are mutually exclusive.

* `fixed_ip` - (Optional) An IP address to assosciate with the port.
   _NOTE_: This option cannot be used with port_id. You must specifiy a network_id. The IP address must lie in a range on the supplied network.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `instance_id` - See Argument Reference above.
* `port_id` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `fixed_ip`  - See Argument Reference above.

## Import

Interface Attachments can be imported using the Instance ID and Port ID
separated by a slash, e.g.

```
$ terraform import huaweicloud_compute_interface_attach.ai_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
