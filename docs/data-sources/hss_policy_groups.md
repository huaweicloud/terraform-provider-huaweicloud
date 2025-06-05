---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_policy_groups"
description: |-
  Use this data source to get the list of HSS policy groups within HuaweiCloud.
---

# huaweicloud_hss_policy_groups

Use this data source to get the list of HSS policy groups within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_policy_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `group_id` - (Optional, String) Specifies the policy group ID.

* `group_name` - (Optional, String) Specifies the policy group name.

* `container_mode` - (Optional, Bool) Specifies whether to query the container edition policy.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The policy group list.
  The [data_list](#policy_group_structure) structure is documented below.

<a name="policy_group_structure"></a>
The `data_list` block supports:

* `group_id` - The policy group ID.

* `group_name` - The policy group name.

* `description` - The description of the policy group.

* `deletable` - Whether a policy group can be deleted.

* `host_num` - The number of associated servers.

* `default_group` - Whether a policy group is the default policy group.

* `support_os` - The supported OS. The valid value are **Linux** and **Windows**.

* `support_version` - The supported versions. The valid values are as follows:
  + **hss.version.basic**: Indicates policy group of the basic edition.
  + **hss.version.advanced**: Indicates policy group of the professional edition.
  + **hss.version.enterprise**: Indicates policy group of the enterprise edition.
  + **hss.version.premium**: Indicates policy group of the premium edition.
  + **hss.version.wtp**: Indicates policy group of the WTP edition.
  + **hss.version.container.enterprise**: Indicates policy group of the container edition.
