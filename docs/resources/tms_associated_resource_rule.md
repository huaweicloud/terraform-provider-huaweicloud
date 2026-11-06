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
data "huaweicloud_tms_associated_resource_settings" "test" {}

resource "huaweicloud_tms_associated_resource_rule" "test" {
  setting_name = data.huaweicloud_tms_associated_resource_settings.test.settings[0].setting_name
  region_id    = "cn-north-4"
  tag_keys     = ["key1", "key2"]
}
```

## Argument Reference

The following arguments are supported:

* `setting_name` - (Required, String, ForceNew) Specifies the name of the rule setting.
  It must be a valid `setting_name` retrieved from the
  `huaweicloud_tms_associated_resource_settings` data source.
  Changing this parameter will create a new resource.

* `region_id` - (Required, String, ForceNew) Specifies the region ID where the rule takes effect.
  It must be one of the `region_ids` of the corresponding setting retrieved from the
  `huaweicloud_tms_associated_resource_settings` data source.
  Changing this parameter will create a new resource.

* `tag_keys` - (Optional, List) Specifies the tag keys that the rule takes effect on.
  Each element represents a tag key. An empty list means the rule takes effect on all tags.

* `existing_resource_status` - (Optional, String) Specifies whether the rule takes effect for existing resources.
  Defaults to **disable**. The valid values are as follows:
  + **enable**
  + **disable**

* `auto_delete_status` - (Optional, String) Specifies whether to automatically delete tags when the association is removed.
  Defaults to **disable**. The valid values are as follows:
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