---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_create_privacy_rules"
description: |-
  Manages a resource to batch create WAF privacy rules within HuaweiCloud WAF.
---

# huaweicloud_waf_batch_create_privacy_rules

Manages a resource to batch create WAF privacy rules within HuaweiCloud WAF.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is a one-time action resource for batch creating privacy rules. Deleting this resource
   will not remove the created rules, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "policy_ids" {
  type = list(string)
}

resource "huaweicloud_waf_batch_create_privacy_rules" "test" {
  url                   = "/admin/xxx"
  category              = "params"
  index                 = "test-index"
  policy_ids            = var.policy_ids
  enterprise_project_id = var.enterprise_project_id
  description           = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `url` - (Required, String, NonUpdatable) Specifies the URL to be protected by the privacy blocking rule.
  You need to enter a standard URL format, such as /admin/xxx or /admin/, with an asterisk (*) at the end to represent
  a path prefix.

* `category` - (Required, String, NonUpdatable) Specifies the masked fields.
  Valid values are **params**, **cookie**, **header**, and **form**.

* `index` - (Required, String, NonUpdatable) Specifies the field name to be masked. The field name will be set according
  to the value of `category`. The masked field will not appear in the log. The length of the masked field name cannot
  exceed `2,048` bytes and can only consist of numbers, letters, underscores and hyphens.

* `policy_ids` - (Required, List, NonUpdatable) Specifies the list of policy IDs to which the rule will be applied.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the rule.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Value `0` indicates the default enterprise project.
  Value **all_granted_eps** indicates all enterprise projects to which the user has been granted access.
  Defaults to `0`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
