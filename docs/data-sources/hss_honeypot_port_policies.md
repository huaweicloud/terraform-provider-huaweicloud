---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_honeypot_port_policies"
description: |-
  Use this data source to get the list of dynamic port honeypot policies.
---

# huaweicloud_hss_honeypot_port_policies

Use this data source to get the list of dynamic port honeypot policies.

## Example Usage

```hcl
data "huaweicloud_hss_honeypot_port_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - All auto launch items that match the filter parameters.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `policy_id` - The dynamic port honeypot policy ID.

* `policy_name` - The dynamic port honeypot policy name.

* `host_num` - The host number.

* `is_default` - Whether the dynamic port honeypot policy is default policy.

* `port_list` - The port list.

* `os_type` - The OS type.

* `status` - The protection status of the dynamic port honeypot policy.
  The valid values are as follows:
  + **applying**: The protection is taking effect.
  + **success**: The protection has taken effect.
  + **disable**: The protection does not take effect.
