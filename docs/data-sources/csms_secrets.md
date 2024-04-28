---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_secrets"
description: ""
---

# huaweicloud_csms_secrets

Use this data source to get the list of CSMS secrets.

## Example Usage

```hcl
variable "secret_name" {}

data "huaweicloud_csms_secrets" "test" {
  name = var.secret_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the secret.

* `secret_id` - (Optional, String) Specifies the ID of the secret.

* `status` - (Optional, String) Specifies the secret status. Valid values are **ENABLED**, **DISABLED**,
  **PENDING_DELETE** and **FROZEN**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `event_name` - (Optional, String) Specifies the event name related to the secret.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `secrets` - Indicates the secrets list.
  The [secrets](#CSMS_secrets) structure is documented below.

<a name="CSMS_secrets"></a>
The `secrets` block supports:

* `id` - The secret ID.

* `name` - Indicates the secret name.

* `status` - Indicates the secret status.

* `kms_key_id` - Indicates the ID of KMS key used to encrypt secret.

* `description` - Indicates the secret description.

* `created_at` - Indicates the time when the secret created, in UTC format.

* `updated_at` - Indicates the time when the secret updated, in UTC format.

* `scheduled_deleted_at` - Indicates the time when the secret is scheduled to be deleted, in UTC format.

* `secret_type` - Indicates the secret type. Valid values are **COMMON** and **RDS**.

* `auto_rotation` - Indicates whether to enable the secret automatic rotation.

* `rotation_period` - Indicates the secret rotation period. Valid when `auto_rotation` is **true**.

* `rotation_config` - Indicates the secret rotation config. Valid when `auto_rotation` is **true**.

* `rotation_at` - Indicates the secret rotation time, in UTC format. Valid when `auto_rotation` is **true**.

* `next_rotation_at` - Indicates the secret next rotation time, in UTC format. Valid when `auto_rotation` is **true**.

* `event_subscriptions` - Indicates the list of events subscribed to by secret.

* `enterprise_project_id` - Indicates the enterprise project ID.
