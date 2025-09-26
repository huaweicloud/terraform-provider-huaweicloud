---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_policy_copy"
description: |-
  Manages a WAF policy copy resource within HuaweiCloud.
---

# huaweicloud_waf_policy_copy

Manages a WAF policy copy resource within HuaweiCloud.

-> This resource is only a one-time action resource using to copy WAF policy. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "src_policy_id" {}
variable "dest_policy_name" {}

resource "huaweicloud_waf_policy_copy" "test" {
  src_policy_id    = var.src_policy_id
  dest_policy_name = var.dest_policy_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this setting will create a new resource.

* `src_policy_id` - (Required, String, NonUpdatable) Specifies the source policy ID.

* `dest_policy_name` - (Required, String, NonUpdatable) Specifies the new policy name to be copied.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**. If you want to query resources under all enterprise projects, set this parameter to
  **all_granted_eps**.

## Attribute Reference

The following attributes are exported:

* `id` - The resource ID, which is the policy ID.
