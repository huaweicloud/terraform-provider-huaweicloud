---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_hda_upgrade_records"
description: |-
  Use this data source to get HDA upgrade records of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_hda_upgrade_records

Use this data source to get HDA upgrade records of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_app_hda_upgrade_records" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the HDA upgrade records are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of HDA upgrade records that matched filter parameters.  
  The [records](#workspace_app_hda_upgrade_records) structure is documented below.

<a name="workspace_app_hda_upgrade_records"></a>
The `records` block supports:

* `server_id` - The ID of the server.

* `machine_name` - The machine name of the server.

* `server_name` - The name of the server.

* `server_group_name` - The name of the server group.

* `sid` - The SID of the server.

* `current_version` - The current version of the access agent.

* `target_version` - The target version of the access agent.

* `upgrade_status` - The HDA upgrade status.
  + **SUCCESS**: Upgrade completed successfully
  + **FAILED**: Upgrade failed
  + **PENDING**: Upgrade pending
  + **RUNNING**: Upgrade in progress

* `upgrade_time` - The upgrade time, in RFC3339 format.
