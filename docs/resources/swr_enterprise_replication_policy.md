---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_replication_policy"
description: |-
  Manages a SWR enterprise instance replication policy resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_replication_policy

Manages a SWR enterprise instance replication policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "name" {}
variable "dest_registry_id" {}

resource "huaweicloud_swr_enterprise_replication_policy" "test"{
  instance_id     = var.instance_id
  name            = var.name
  enabled         = true
  repo_scope_mode = "regular"
  description     = "demo"

  dest_registry {
    id = var.dest_registry_id
  }

  filters {
    type  = "name"
    value = "**/**"
  }
    
  filters {
    type  = "tag"
    value = "**"
  }

  trigger {
    trigger_settings {
      cron = "0 0 0 1 * ?"
    }

    type = "scheduled"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `name` - (Required, String) Specifies the policy name.

* `repo_scope_mode` - (Required, String) Specifies the repo scope mode.
  Values can be **regular** and **selection**.

* `enabled` - (Required, Bool) Specifies whether the policy is enabled. Defaults to **false**.

* `filters` - (Required, List) Specifies the source resource filter.
  The [filters](#block--filters) structure is documented below.

* `trigger` - (Required, List) Specifies the trigger config.
  The [trigger](#block--trigger) structure is documented below.

* `description` - (Optional, String) Specifies the description of policy.

* `dest_namespace` - (Optional, String) Specifies the destination namespace name.

* `override` - (Optional, Bool) Specifies whether to override the repository. Defaults to **false**.

* `src_registry` - (Optional, List) Specifies the source registry infos.
  The [registry](#block--registry) structure is documented below.

* `dest_registry` - (Optional, List) Specifies the destination registry infos.
  The [registry](#block--registry) structure is documented below.

-> At least one of the `src_registry` and `dest_registry` must be specified.

<a name="block--registry"></a>
The `registry` block supports:

* `id` - (Required, Int) Specifies the registry ID.

<a name="block--filters"></a>
The `filters` block supports:

* `type` - (Required, String) Specifies the filter type.
  Values can be **name** and **tag**.

* `value` - (Required, String) Specifies the regular expression of the filter.

<a name="block--trigger"></a>
The `trigger` block supports:

* `type` - (Required, String) Specifies the trigger type.
  Values can be **manual**, **scheduled** and **event_based**.

* `trigger_settings` - (Optional, List) Specifies the trigger settings.
  The [trigger_settings](#block--trigger--trigger_settings) structure is documented below.

  -> Only when `trigger.0.type` is **scheduled**, `trigger_settings` should be specified.

<a name="block--trigger--trigger_settings"></a>
The `trigger_settings` block supports:

* `cron` - (Optional, String) Specifies the scheduled setting.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `policy_id` - Indicates the policy ID.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

## Import

The replication policy can be imported using `instance_id` and `policy_id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_replication_policy.test <instance_id>/<policy_id>
```
