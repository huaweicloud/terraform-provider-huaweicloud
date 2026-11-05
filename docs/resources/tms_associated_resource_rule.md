---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_associated_resource_rule"
description: |-
  Manages a TMS associated resource rule within HuaweiCloud.
---

# huaweicloud_tms_associated_resource_rule

Manages a TMS associated resource rule within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_tms_associated_resource_rule" "test" {
  setting_name             = "ecs_ebs"
  region_id                = "cn-north-4"
  tag_keys                 = ["key1", "key2"]
  existing_resource_status = "enable"
  auto_delete_status       = "disable"
}
```

## Argument Reference

The following arguments are supported:

* `setting_name` - (Required, String, ForceNew) Specifies the name of the rule setting.
  Changing this parameter will create a new resource.

* `region_id` - (Required, String, ForceNew) Specifies the region ID where the rule takes effect.
  Changing this parameter will create a new resource.

* `tag_keys` - (Required, List) Specifies the tag keys associated with the rule.

* `existing_resource_status` - (Optional, String) Specifies whether the rule takes effect for existing resources.
  The valid values are as follows:
  + **enable**
  + **disable**

* `auto_delete_status` - (Optional, String) Specifies whether to automatically delete tags when the association is removed.
  The valid values are as follows:
  + **enable**
  + **disable**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Associated resource rules can be imported using the `setting_name` and `region_id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_tms_associated_resource_rule.test ecs_ebs/cn-north-4
```