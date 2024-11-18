---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_group"
description: |-
  Manages a log group resource within HuaweiCloud.
---

# huaweicloud_lts_group

Manages a log group resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_lts_group" "test" {
  group_name  = "log_group1"
  ttl_in_days = 30
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the log group resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `group_name` - (Required, String, ForceNew) Specifies the log group name. Changing this parameter will create a new resource.

* `ttl_in_days` - (Required, Int) Specifies the log expiration time(days).  
  The value is range from `1` to `365`.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the log group.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the log group belongs.
  Changing this parameter will create a new resource.  
  This parameter is valid only when the enterprise project function is enabled, if omitted, default enterprise project
  will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The log group ID.

* `created_at` - The creation time of the log group.

## Import

The log group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_lts_group.test <id>
```
