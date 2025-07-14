---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_drills"
description: |-
  Use this data source to query SDRS disaster recovery drills within HuaweiCloud.
---

# huaweicloud_sdrs_drills

Use this data source to query SDRS disaster recovery drills within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_sdrs_drills" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_group_id` - (Optional, String) Specifies the ID of a protection group.
  The value of this parameter can query from datasource `huaweicloud_sdrs_protection_groups`.

* `name` - (Optional, String) Specifies the DR drill name. Fuzzy search is supported.

* `status` - (Optional, String) Specifies the DR drill status.
  For details, see [DR Drill Status](https://support.huaweicloud.com/intl/en-us/api-sdrs/en-us_topic_0126152933.html).

* `drill_vpc_id` - (Optional, String) Specifies the ID of the VPC used for a DR drill.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `disaster_recovery_drills` - The DR drills.

  The [disaster_recovery_drills](#disaster_recovery_drills_struct) structure is documented below.

<a name="disaster_recovery_drills_struct"></a>
The `disaster_recovery_drills` block supports:

* `name` - Tthe DR drill name.

* `status` - The DR drill status.

* `drill_vpc_id` - The ID of the VPC used for a DR drill.

* `created_at` - The time when a DR drill was created.
  The default format is as follows: "yyyy-MM-dd HH:mm:ss.SSS", for example, **2019-04-01 12:00:00.000**.

* `updated_at` - The time when a DR drill was updated.
  The default format is as follows: "yyyy-MM-dd HH:mm:ss.SSS", for example, **2019-04-01 12:00:00.000**.

* `server_group_id` - The ID of a protection group.

* `drill_servers` - The drill servers.

  The [drill_servers](#disaster_recovery_drills_drill_servers_struct) structure is documented below.

* `id` - The DR drill ID.

<a name="disaster_recovery_drills_drill_servers_struct"></a>
The `drill_servers` block supports:

* `protected_instance` - The protected instance ID of the drill server.

* `drill_server_id` - The drill server ID.
