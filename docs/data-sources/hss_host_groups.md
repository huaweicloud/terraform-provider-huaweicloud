---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_host_groups"
description: ""
---

# huaweicloud_hss_host_groups

Use this data source to get the list of HSS host groups within HuaweiCloud.

## Example Usage

```hcl
variable group_id {}

data "huaweicloud_hss_host_groups" "test" {
  group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS host groups.
  If omitted, the provider-level region will be used.

* `group_id` - (Optional, String) Specifies the ID of the host group to be queried.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the host groups
  belong. If omitted, will query the values under the default enterprise project.

* `name` - (Optional, String) Specifies the name of the host group to be queried. This field will undergo a fuzzy
  matching query, the query result is for all host groups whose names contain this value.

* `host_num` - (Optional, String) Specifies the number of hosts in the host groups to be queried.

* `risk_host_num` - (Optional, String) Specifies the number of risky hosts in the host groups to be queried.

* `unprotect_host_num` - (Optional, String) Specifies the number of unprotected hosts in the host groups to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `groups` - All host groups that match the filter parameters.

  The [groups](#hss_groups) structure is documented below.

<a name="hss_groups"></a>
The `groups` block supports:

* `id` - The ID of the host group.

* `name` - The name of the host group.

* `host_num` - The number of hosts in the host group.

* `risk_host_num` - The number of risky hosts in the host group.

* `unprotect_host_num` - The number of unprotected hosts in the host group.

* `host_ids` - The list of host IDs in the host group.
