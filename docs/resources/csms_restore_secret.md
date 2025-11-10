---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_restore_secret"
description: |-
  Restores a secret from a backup content in HuaweiCloud DEW CSMS service.
---

# huaweicloud_csms_restore_secret

Restores a secret from a backup content in HuaweiCloud DEW CSMS service.

-> This resource is a one-time action resource. Deleting this resource will not affect the restored secret,
  but will only remove the resource information from the tfstate file. This resource can only restore secrets that no
  longer exist.

## Example Usage

```hcl
resource "huaweicloud_csms_restore_secret" "test" {
  secret_blob = file("your-secret-backup-file-path")
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `secret_blob` - (Required, String, NonUpdatable) Specifies the secret backup file content to restore from.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the restored secret.

* `name` - The name of the restored secret.

* `state` - The state of the restored secret. Possible values are:
  + **ENABLED**: The secret is enabled.
  + **DISABLED**: The secret is disabled.
  + **PENDING_DELETE**: The secret is pending deletion.
  + **FROZEN**: The secret is frozen.

* `kms_key_id` - The ID of the KMS key used to encrypt the secret value.

* `description` - The description of the secret.

* `create_time` - The creation time of the secret, in milliseconds since the Unix epoch.

* `update_time` - The last update time of the secret, in milliseconds since the Unix epoch.

* `scheduled_delete_time` - The scheduled deletion time of the secret, in milliseconds since the Unix epoch.

* `secret_type` - The type of the secret. Possible values are:
  + **COMMON**: The shared secret (default), which is used to store sensitive information in an application system.
  + **RDS**: The RDS secret, which is used to store RDS account information. (no longer supported, replaced by RDS-FG).
  + **RDS-FG**: The RDS secret, which is used to store RDS account information.
  + **GaussDB-FG**: The taurusDB secret, which is used to store taurusDB account information.

* `auto_rotation` - The Automatic rotation. The value can be **true** (enabled) or **false** (disabled).
  The default value is **false**.

* `rotation_period` - The rotation period of the secret.

* `rotation_config` - The rotation configuration of the secret.

* `rotation_time` - The last rotation time of the secret, in milliseconds since the Unix epoch.

* `next_rotation_time` - The next rotation time of the secret, in milliseconds since the Unix epoch.

* `event_subscriptions` - The list of events subscribed to by secrets. Currently, only one event can be subscribed to.
  When a basic event contained in an event is triggered, a notification message is sent to the notification topic
  corresponding to the event.

* `enterprise_project_id` - The ID of the enterprise project that the secret belongs to.

* `rotation_func_urn` - The URN of the FunctionGraph function used for rotation.
