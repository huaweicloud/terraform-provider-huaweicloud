---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_secret"
description: |
  Manages CSMS(Cloud Secret Management Service) secrets within HuaweiCloud.
---

# huaweicloud_csms_secret

Manages CSMS(Cloud Secret Management Service) secrets within HuaweiCloud.

## Example Usage

### Encrypt Plaintext

```hcl
resource "huaweicloud_csms_secret" "test1" {
  name        = "test_secret"
  secret_text = "this is a password"
}
```

### Encrypt JSON Data

```hcl
resource "huaweicloud_csms_secret" "test2" {
  name        = "mysql_admin"
  secret_text = jsonencode({
    username = "admin"
    password = "123456"
  })
}
```

### Encrypt String Binary

```hcl
variable "secret_binary" {}

resource "huaweicloud_csms_secret" "test3" {
  name          = "test_binary"
  secret_binary = var.secret_binary
}
```

### The secret associated event

```hcl
variable "name" {}
variable "secret_type" {}
variable "secret_text" {}

resource "huaweicloud_csms_event" "test" {
  ...
}

resource "huaweicloud_csms_secret" "test" {
  name                = var.name
  secret_type         = var.secret_type
  secret_text         = var.secret_text
  event_subscriptions = [huaweicloud_csms_event.test.name]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CSMS secrets.
  If omitted, the provider-level region will be used. Changing this setting will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the secret name. The maximum length is 64 characters.
  Only digits, letters, underscores(_), hyphens(-) and dots(.) are allowed.

  Changing this parameter will create a new resource.

* `secret_text` - (Optional, String) Specifies the plaintext of a text secret. CSMS encrypts the plaintext and stores
  it in the initial version of the secret. The maximum size is 32 KB.

* `secret_binary` - (Optional, String) Specifies the plaintext of a binary secret encoded using Base64. CSMS encrypts
  the plaintext and stores it in the initial version of the secret. The maximum size is 32 KB.

-> 1. One of the fields `secret_text` and `secret_binary` must be configured, and can not be specified both together. The
  `secret_text` and `secret_binary` are sensitive, and we store their hashes in the state file.
  <br/>2. Whenever the `secret_text` or `secret_binary` parameters are changed, the latest version is incremented.

* `expire_time` - (Optional, Int) Specifies the expiration time of a secret, `expire_time` can only be edited
  when `status` is **ENABLED**. The time is in the format of timestamp, that is, the offset milliseconds
  from 1970-01-01 00:00:00 UTC to the specified time. The time must be greater than the current time.

  -> Due to API reasons, please ensure that the last three digits of the millisecond timestamp are `0`, otherwise changes
  will be triggered. For example, `1729243021000`.

* `kms_key_id` - (Optional, String) Specifies the ID of the KMS key used to encrypt secrets.
  If this parameter is not specified when creating the secret, the default master key **csms/default** will be used.
  The default key is automatically created by the CSMS.
  Use this data source
  [huaweicloud_kms_keys](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs/data-sources/kms_keys)
  to get the KMS key.

* `description` - (Optional, String) Specifies the description of a secret.

* `secret_type` - (Optional, String, ForceNew) Specifies the type of the secret.
  Currently, only supported **COMMON**. The default value is **COMMON**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the secret belongs.
  If omitted, the default enterprise project will be used.
  If the enterprise project function is not enabled, ignore this parameter.

  Changing this parameter will create a new resource.

* `event_subscriptions` - (Optional, List) Specifies the event list associated with the secret.
  Currently, only one event can be associated.

* `tags` - (Optional, Map) Specifies the tags of a CSMS secrets, key/value pair format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is constructed from the secret ID and name, separated by a slash.

* `secret_id` - The secret ID in UUID format.

* `latest_version` - The latest version id.

* `version_stages` - The secret version status list.

* `status` - The CSMS secret status. Values can be: **ENABLED**, **DISABLED**, **PENDING_DELETE** and **FROZEN**.

* `create_time` - Time when the CSMS secrets created, in UTC format.

## Import

CSMS secret can be imported using the ID and the name of secret, separated by a slash, e.g.

```bash
terraform import huaweicloud_csms_secret.test <id>/<name>
```
