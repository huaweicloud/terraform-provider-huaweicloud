---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_drill"
description: ""
---

# huaweicloud_sdrs_drill

Manages an SDRS DR drill resource within HuaweiCloud.

-> This resource requires the protection group has at least one protected instance or replication pair.
And there has some usage restrictions for this resource,
refer to [document](https://support.huaweicloud.com/intl/en-us/qs-sdrs/en-us_topic_0122528555.html).

## Example Usage

```hcl
variable "group_id" {}

resource "huaweicloud_sdrs_drill" "test" {
  name     = "test-drill"
  group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of a DR drill. The name can contain a maximum of 64 characters.
  The value can contain only letters (a to z and A to Z), digits (0 to 9), dots (.), underscores (_), and hyphens (-).

* `group_id` - (Required, String, ForceNew) Specifies the ID of a protection group.

  Changing this parameter will create a new resource.

* `drill_vpc_id` - (Optional, String, ForceNew) Specifies the drill VPC ID.
  
  If specified, make sure the drill VPC CIDR block consistent with the VPC of the protection group, and make sure the
  drill VPC has subnets.
  
  If not specified, the system automatically creates a drill VPC.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of a DR drill.
  + **creating**: The DR drill is being created.
  + **available**: The DR drill is available.
  + **deleting**: The DR drill is being deleted.
  + **error-deleting**: Failed to delete the DR drill.
  + **error**: Failed to create the DR drill.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The SDRS DR drill can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_sdrs_drill.test <id>
```
