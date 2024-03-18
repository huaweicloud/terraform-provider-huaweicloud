---
subcategory: "Data Encryption Workshop (DEW)"
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

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the CSMS secrets.
  If omitted, the provider-level region will be used. Changing this setting will create a new resource.

* `name` - (Required, String, ForceNew) The secret name. The maximum length is 64 characters.
  Only digits, letters, underscores(_), hyphens(-) and dots(.) are allowed.

* `secret_text` - (Optional, String) Specifies the plaintext of a secret in text format. The maximum size is 32 KB.

* `secret_binary` - (Optional, String) Specifies the plaintext of a binary secret in text format encoded using Base64.
  The maximum size is 32 KB.

-> **NOTE:** Exactly one of `secret_text` and `secret_binary` should be specified.
  The `secret_text` and `secret_binary` are sensitive and in the state file we store its hash.

* `expire_time` - (Optional, Int) Specifies the expiration time of a secret, `expire_time` can only be edited
  when `status` is **ENABLED**. The time is in the format of timestamp, that is, the offset milliseconds
  from 1970-01-01 00:00:00 UTC to the specified time. The time must be greater than the current time.

* `kms_key_id` - (Optional, String) The ID of the KMS key used to encrypt secrets.
  If this parameter is not specified when creating the secret, the default master key csms/default will be used.
  The default key is automatically created by the CSMS.
  Use this data source
  [huaweicloud_kms_key](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs/resources/kms_key)
  to get the KMS key.

* `description` - (Optional, String) The description of a secret.

* `tags` - (Optional, Map) The tags of a CSMS secrets, key/value pair format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is constructed from the secret ID and name, separated by slashes.

* `secret_id` - The secret ID in UUID format.

* `latest_version` - The latest version id.

* `status` - The CSMS secret status. Values can be: ENABLED, DISABLED, PENDING_DELETE and FROZEN.

* `create_time` - Time when the CSMS secrets created, in UTC format.

## Import

CSMS secret can be imported using the ID and the name of secret, separated by a slash, e.g.

```sh
terraform import huaweicloud_csms_secret.test 93cba7f5-550b-45dc-912e-277b3296fb27/test_secret
```
