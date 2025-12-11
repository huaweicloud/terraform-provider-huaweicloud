---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_webtamper_protection_statistics"
description: |-
  Use this data source to query the protection data statistics.
---

# huaweicloud_hss_webtamper_protection_statistics

Use this data source to query the protection data statistics.

## Example Usage

```hcl
data "huaweicloud_hss_webtamper_protection_statistics" "test" {}
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

* `protect_host_num` - The total number of protected servers, including servers in the protected, partially protected,
  protection suspended, protection failed, protection interrupted, and enabling states.

* `protect_success_host_num` - The total number of servers in the protected state.

* `protect_fail_host_num` - The total number of servers where protection failed.

* `anti_tampering_num` - The events in the last 168 hours.
