---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_network_security_groups"
description: |-
  Use this data source to get the list of HSS container network security groups within HuaweiCloud.
---

# huaweicloud_hss_container_network_security_groups

Use this data source to get the list of HSS container network security groups within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_network_security_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need to set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `security_groups` - The list of security groups.

  The [security_groups](#security_groups_struct) structure is documented below.

<a name="security_groups_struct"></a>
The `security_groups` block supports:

* `security_group_id` - The security group ID.

* `security_group_name` - The security group name.

* `security_group_description` - The security group description.
