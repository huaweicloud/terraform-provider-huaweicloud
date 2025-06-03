---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_secrets_by_tags"
description: |-
  Use this data source to get a list of the secrets by tags.
---

# huaweicloud_csms_secret_versions

Use this data source to get a list of the secrets by tags.

## Example Usage

```hcl
variable "resource_instances" {}
variable "action" {}

data "huaweicloud_csms_secrets_by_tags" "test" {
  resource_instances = var.resource_instances
  action             = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_instances` - (Required, String) Specifies the resource instances, the valid value is **resource_instances**.

* `action` - (Required, String) Specifies the operation type. The valid values are as follows:
  + **filter**: Indicates filtering secrets.
  + **count**: Indicates the total number of secrets.

* `tags` - (Optional, List) Specifies the list of tags, the maximum of tags is `10`.
  The [tags](#tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the key-value pair to be matched.
  The [matches](#matches_struct) structure is documented below.

* `sequence` - (Optional, String) Specifies the `36` byte sequence number of a request message.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Optional, String) Specifies the tag key.

* `values` - (Optional, List) Specifies the set of tag values, the maximum of values is `10`.
  If the tag list is empty, any value can be matched.
  A search result matches only one value.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Optional, String) Specifies the search field, the valid value is **resource_name**.

* `value` - (Optional, String) Specifies the field for fuzzy match, maximum of `255` characters are allowed.
  If it is left blank, a null value is returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of the filtered secrets.

  The [resources](#resources_struct) structure is documented below.

* `total_count` - The total number of the filtered secrets.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - The secret ID.

* `resource_name` - The secret name.

* `resource_detail` - The secret detail.

  The [resource_detail](#resource_detail_struct) structure is documented below.

* `tags` - The tag list.

  The [tags](#tags_item_struct) structure is documented below.

* `sys_tags` - The system tag list.

  The [sys_tags](#sys_tags_struct) structure is documented below.

<a name="resource_detail_struct"></a>
The `resource_detail` block supports:

* `id` - The ID of the secret.

* `name` - The secret name.

* `state` - The secret status. The valid values are as follows:
  + **ENABLED**: Indicates enabled status.
  + **DISABLED**: Indicates disabled status.
  + **PENDING_DELETE**: Indicates pending deletion status.
  + **FROZEN**: Indicates frozen state.

* `kms_key_id` - The ID of KMS key used to encrypt secret.

* `description` - The description of the secret.

* `create_time` - The creation time of the secret, the value is a timestamp.

* `update_time` - The update time of the secret, the value is a timestamp.

* `scheduled_delete_time` - The time of the secret to be scheduled deleted, the value is a timestamp.

* `secret_type` - The secret type. The valid values are as follows:
  + **COMMON**: shared secret (default), which is used to store sensitive information in an application system.
  + **RDS**: RDS secret, which is used to store RDS account information. (no longer supported, replaced by RDS-FG).
  + **RDS-FG**: RDS secret, which is used to store RDS account information.
  + **GaussDB-FG**: TaurusDB secret, which is used to store TaurusDB account information.

* `auto_rotation` - Automatic rotation. The valid values are as follows:
  + **true**: Enabled.
  + **false**: Disabled.

* `rotation_period` - The secret rotation period. Valid when `auto_rotation` is **true**.

* `rotation_config` - The secret rotation config. Valid when `auto_rotation` is **true**.

* `rotation_time` - The rotation time of the secret, the value is a timestamp.

* `next_rotation_time` - The next rotation time of the secret, the value is a timestamp.

* `event_subscriptions` - The list of events subscribed to by secret.

* `enterprise_project_id` - The enterprise project ID.

<a name="tags_item_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.

<a name="sys_tags_struct"></a>
The `sys_tags` block supports:

* `key` - The system tag key.

* `value` - The system tag value.
