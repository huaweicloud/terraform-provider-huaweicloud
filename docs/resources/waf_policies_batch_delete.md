---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_policies_batch_delete"
description: |-
  Manages a WAF policies batch delete resource within HuaweiCloud.
---

# huaweicloud_waf_policies_batch_delete

Manages a WAF policies batch delete resource within HuaweiCloud.

-> This resource is only a one-time action resource using to batch delete WAF policies. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "policy_ids" {
  type = list(string)
}

resource "huaweicloud_waf_policies_batch_delete" "test" {
  policy_ids = var.policy_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this setting will create a new resource.

* `policy_ids` - (Optional, List, NonUpdatable) Specifies the list of policy IDs to be deleted.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

The following attributes are exported:

* `id` - The resource ID in UUID format.
