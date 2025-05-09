---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project_action"
description: |-
 Use this resource to operate the enterprise project within HuaweiCloud.
---

# huaweicloud_enterprise_project_action

Use this resource to operate the enterprise project within HuaweiCloud.

-> This resource is only a one-time action resource for operating the enterprise project. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "enterprise_project_id" {}

resource "huaweicloud_enterprise_project_action" "test" {
  enterprise_project_id = var.enterprise_project_id
  action                = "disable"
}
```

## Argument Reference

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the ID of enterprise project to be operated.

* `action` - (Required, String, NonUpdatable) Specifies the action type.  
  The valid values are as follows:
  + **enable**
  + **disable**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
