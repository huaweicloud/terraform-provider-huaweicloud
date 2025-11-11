---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_replication_policies"
description: |-
  Use this data source to get the list of SWR enterprise instance replication policies.
---

# huaweicloud_swr_enterprise_replication_policies

Use this data source to get the list of SWR enterprise instance replication policies.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_replication_policies" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `registry_id` - (Optional, Int) Specifies the registry ID.

* `name` - (Optional, String) Specifies the policy name.

* `order_column` - (Optional, String) Specifies the order column.
  Values can be **created_at** or **updated_at**. Default to **created_at**.

* `order_type` - (Optional, String) Specifies the order type. Values can be **desc** or **asc**.
  `order_column` is required if `order_type` is specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - Indicates the policies.

  The [policies](#policies_struct) structure is documented below.

* `total` - Indicates the total counts of policies.

<a name="policies_struct"></a>
The `policies` block supports:

* `id` - Indicates the policy ID.

* `name` - Indicates the policy name.

* `description` - Indicates the policy description.

* `repo_scope_mode` - Indicates the repo scope mode.

* `override` - Indicates whether the repository is overrided.

* `enabled` - Indicates whether the repository is enabled.

* `src_registry` - Indicates the source registry infos.

  The [src_registry](#policies_src_registry_struct) structure is documented below.

* `dest_registry` - Indicates the destination registry infos.

  The [dest_registry](#policies_dest_registry_struct) structure is documented below.

* `dest_namespace` - Indicates the destination namespace.

* `filters` - Indicates the source resource filter.

  The [filters](#policies_filters_struct) structure is documented below.

* `trigger` - Indicates the trigger config.

  The [trigger](#policies_trigger_struct) structure is documented below.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

<a name="policies_src_registry_struct"></a>
The `src_registry` block supports:

* `id` - The source registry ID.

<a name="policies_dest_registry_struct"></a>
The `dest_registry` block supports:

* `id` - Indicates the destination registry ID.

<a name="policies_filters_struct"></a>
The `filters` block supports:

* `value` - Indicates the regular expression of the filter.

* `type` - Indicates the filter type.

<a name="policies_trigger_struct"></a>
The `trigger` block supports:

* `type` - Indicates the trigger type.

* `trigger_settings` - Indicates the trigger settings.

  The [trigger_settings](#trigger_trigger_settings_struct) structure is documented below.

<a name="trigger_trigger_settings_struct"></a>
The `trigger_settings` block supports:

* `cron` - Indicates the scheduled setting.
