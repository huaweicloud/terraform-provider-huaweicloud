---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_rule_restriction"
description: |-
  Manages a Workspace restriction resource to restrict application rules within HuaweiCloud.
---

# huaweicloud_workspace_application_rule_restriction

Manages a Workspace restriction resource to restrict application rules within HuaweiCloud.

## Example Usage

```hcl
variable "rule_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_application_rule_restriction" "test" {
  rule_ids = var.rule_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the restricted application rules are located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `rule_ids` - (Required, List) Specifies the list of application rule IDs to be restricted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
