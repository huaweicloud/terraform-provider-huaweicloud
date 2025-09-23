---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshot"
description: ""
---

# huaweicloud_evs_snapshot

Provides an EVS snapshot resource.

## Example Usage

```hcl
resource "huaweicloud_evs_volume" "test" {
  name        = "volume"
  description = "my volume"
  volume_type = "SATA"
  size        = 20

  availability_zone = "cn-north-4a"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_evs_snapshot" "test" {
  name        = "snapshot-001"
  description = "Daily backup"
  volume_id   = huaweicloud_evs_volume.test.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the evs snapshot resource. If omitted, the
  provider-level region will be used. Changing this creates a new EVS snapshot resource.

* `volume_id` - (Required, String, ForceNew) The id of the snapshot's source disk. Changing the parameter creates a new
  snapshot.

* `name` - (Required, String) The name of the snapshot. The value can contain a maximum of 255 bytes.

* `metadata` - (Optional, Map, ForceNew) Specifies the user-defined metadata key-value pair. Changing the parameter
  creates a new snapshot.

* `description` - (Optional, String) The description of the snapshot. The value can contain a maximum of 255 bytes.

* `force` - (Optional, Bool) Specifies the flag for forcibly creating a snapshot. Default to false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the snapshot.

* `status` - The status of the snapshot.

* `size` - The size of the snapshot in GB.

* `created_at` - The time when the snapshot was created.

* `updated_at` - The time when the snapshot was updated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 3 minutes.

## Import

EVS snapshot can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_evs_snapshot.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `metadata`, `force`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the snapshot. Also, you can
ignore changes as below.

```hcl
resource "huaweicloud_evs_snapshot" "test" {
    ...

  lifecycle {
    ignore_changes = [
      metadata, force,
    ]
  }
}
```
