---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_replication_pair"
description: ""
---

# huaweicloud_sdrs_replication_pair

Manages an SDRS replication pair resource within HuaweiCloud.

## Example Usage

```hcl
variable "protection_group_id" {}
variable "volume_id" {}

resource "huaweicloud_sdrs_replication_pair" "test" {
  name                 = "test-replication-pair"
  group_id             = var.protection_group_id
  volume_id            = var.volume_id
  description          = "test description"
  delete_target_volume = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of a replication pair. The name can contain a maximum of 64 characters.
  The value can contain only letters (a to z and A to Z), digits (0 to 9), dots (.), underscores (_), and hyphens (-).

* `group_id` - (Required, String, ForceNew) Specifies the ID of a protection group.

  Changing this parameter will create a new resource.

* `volume_id` - (Required, String, ForceNew) Specifies the ID of the production site disk.
  When the provider is successfully invoked, the disaster recovery site disk will be automatically created.

  Changing this parameter will create a new resource.

* `delete_target_volume` - (Optional, Bool) Specifies whether to delete the disaster recovery site disk.
  The default value is **false**.

* `description` - (Optional, String, ForceNew) Specifies the description of a replication pair. The value can contain
  a maximum of 64 characters and angle brackets (<) and (>) are not allowed.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `fault_level` - The fault level of a replication pair.
  + 0: No fault occurs.
  + 2: The disk of the current production site does not have read/write permissions. In this case, you are advised to
  perform a failover.
  + 5: The replication link is disconnected. In this case, a failover is not allowed. Contact the customer service to
  obtain service support.

* `replication_model` - The replication mode of a replication pair. The default value is **hypermetro**,
  indicating synchronous replication.

* `status` - The status of a replication pair. For details,
  see [Replication Pair Status](https://support.huaweicloud.com/intl/en-us/api-sdrs/en-us_topic_0126152932.html).

* `target_volume_id` - The ID of the disk in the protection availability zone.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The SDRS replication pair can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_sdrs_replication_pair.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `delete_target_volume`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_sdrs_replication_pair" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      delete_target_volume,
    ]
  }
}
```
