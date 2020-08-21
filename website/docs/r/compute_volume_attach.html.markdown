---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_volume_attach"
sidebar_current: "docs-huaweicloud-resource-compute-volume-attach"
description: |-
  Attaches a Volume to an Instance.
---

# huaweicloud\_compute\_volume_attach

Attaches a Volume to an Instance.
This is an alternative to `huaweicloud_compute_volume_attach_v2`

## Example Usage

### Basic attachment of a single volume to a single instance

```hcl
resource "huaweicloud_evs_volume" "test" {
  name = "volume_1"
  availability_zone = "cn-norht-1a"
  volume_type = "SAS"
  size = 10
}

resource "huaweicloud_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "huaweicloud_compute_volume_attach" "va_1" {
  instance_id = "${huaweicloud_compute_instance.instance_1.id}"
  volume_id = "${huaweicloud_evs_volume.test.id}"
}
```

### Attaching multiple volumes to a single instance

```hcl
resource "huaweicloud_evs_volume" "test" {
  count = 2
  name = "volume_1"
  availability_zone = "cn-norht-1a"
  volume_type = "SAS"
  size = 10
}

resource "huaweicloud_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "huaweicloud_compute_volume_attach" "attachments" {
  count       = 2
  instance_id = "${huaweicloud_compute_instance.instance_1.id}"
  volume_id   = "${element(huaweicloud_evs_volume.test.*.id, count.index)}"
}

output "volume devices" {
  value = "${huaweicloud_compute_volume_attach.attachments.*.device}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    A Compute client is needed to create a volume attachment. If omitted, the
    `region` argument of the provider is used. Changing this creates a
    new volume attachment.

* `instance_id` - (Required) The ID of the Instance to attach the Volume to.

* `volume_id` - (Required) The ID of the Volume to attach to an Instance.

* `device` - (Optional) The device of the volume attachment (ex: `/dev/vdc`).
  _NOTE_: Being able to specify a device is dependent upon the hypervisor in
  use. There is a chance that the device specified in Terraform will not be
  the same device the hypervisor chose. If this happens, Terraform will wish
  to update the device upon subsequent applying which will cause the volume
  to be detached and reattached indefinitely. Please use with caution.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `instance_id` - See Argument Reference above.
* `volume_id` - See Argument Reference above.
* `device` - See Argument Reference above. _NOTE_: The correctness of this
  information is dependent upon the hypervisor in use. In some cases, this
  should not be used as an authoritative piece of information.
* `pci_address` - PCI address of the block device.

## Import

Volume Attachments can be imported using the Instance ID and Volume ID
separated by a slash, e.g.

```
$ terraform import huaweicloud_compute_volume_attach.va_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
