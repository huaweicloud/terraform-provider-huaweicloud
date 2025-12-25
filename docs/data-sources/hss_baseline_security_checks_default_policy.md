---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_security_checks_default_policy"
description: |-
  Use this data source to get the HSS baseline security checks default policy within HuaweiCloud.
---

# huaweicloud_hss_baseline_security_checks_default_policy

Use this data source to get the HSS baseline security checks default policy within HuaweiCloud.

## Example Usage

```hcl
variable "support_os" {}

data "huaweicloud_hss_baseline_security_checks_default_policy" "test" {
  support_os = var.support_os
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `support_os` - (Required, String) Specifies the operating system type supported by the policy.  
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `content` - The policy detail.
