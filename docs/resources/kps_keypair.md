---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_keypair"
description: ""
---

# huaweicloud_kps_keypair

Manages a keypair resource within HuaweiCloud.

By default, key pairs use the SSH-2 (RSA, 2048) algorithm for encryption and decryption.

Keys imported support the following cryptographic algorithms:

 * RSA-1024
 * RSA-2048
 * RSA-4096

## Example Usage

### Create a new keypair and export private key to current folder

```hcl
resource "huaweicloud_kps_keypair" "test-keypair" {
  name     = "my-keypair"
  key_file = "private_key.pem"
}
```

### Create a new keypair which scope is Tenant-level and the private key is managed by HuaweiCloud

```hcl
resource "huaweicloud_kms_key" "test" {
  key_alias = "kms_test"
}

resource "huaweicloud_kps_keypair" "test-keypair" {
  name            = "my-keypair"
  scope           = "account"
  encryption_type = "kms"
  kms_key_name    = huaweicloud_kms_key.test.key_alias
}
```

### Import an existing keypair

```hcl
variable "public_key" {}
variable "kms_key_name" {}
variable "private_key" {}

resource "huaweicloud_kps_keypair" "test-keypair" {
  name            = "my-keypair"
  public_key      = var.public_key
  encryption_type = "kms"
  kms_key_name    = var.kms_key_name
  private_key     = var.private_key
}
```

### Import a keypair without a platform-managed private key

```hcl
variable "public_key" {}

resource "huaweicloud_kps_keypair" "test-keypair" {
  name       = "my-keypair"
  public_key = var.public_key
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the keypair resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies a unique name for the keypair. The name can contain a maximum of `64`
  characters, including letters, digits, underscores (_) and hyphens (-).
  Changing this parameter will create a new resource.

* `scope` - (Optional, String, ForceNew) Specifies the scope of key pair. The options are as follows:
  - **account**: Tenant-level, available to all users under the same account.
  - **user**: User-level, only available to that user.
  The default value is `user`.
  Changing this parameter will create a new resource.

* `user_id` - (Optional, String) Specifies the user ID to which the keypair belongs.

  -> If the `scope` set to **user**, this parameter value must be the ID of the user who creates the resource.

* `encryption_type` - (Optional, String) Specifies encryption mode if manages the private key by HuaweiCloud.
  The options are as follows:
  - **default**: The default encryption mode. Applicable to sites where KMS is not deployed.
  - **kms**: KMS encryption mode.

* `kms_key_id` - (Optional, String) Specifies the KMS key ID to encrypt private keys.

* `kms_key_name` - (Optional, String) Specifies the KMS key name to encrypt private keys.

-> When the `encryption_type` set to **kms**, exactly one of `kms_key_id` or `kms_key_name` must be set.

* `description` - (Optional, String) Specifies the description of key pair.

* `public_key` - (Optional, String, ForceNew) Specifies the imported OpenSSH-formatted public key.
  It is required when import keypair. Changing this parameter will create a new resource.

* `private_key` - (Optional, String) Specifies the imported OpenSSH-formatted private key.

* `key_file` - (Optional, String, ForceNew) Specifies the path of the created private key.
  The private key file (**.pem**) is created only after the resource is created.
  Changing this parameter will create a new resource.

  ->**NOTE:** If the private key file already exists, it will be overwritten after a new keypair is created.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to name.

* `created_at` - The key pair creation time.

* `fingerprint` - Fingerprint information about an key pair.

* `is_managed` - Whether the private key is managed by HuaweiCloud.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

Keypairs can be imported using the `name`, e.g.

```bash
$ terraform import huaweicloud_kps_keypair.my-keypair test-keypair
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `encryption_type`, `kms_key_id`,
`kms_key_name` and `private_key`. It is generally recommended running `terraform plan` after importing a key pair.
You can then decide if changes should be applied to the key pair, or the resource definition
should be updated to align with the key pair. Also you can ignore changes as below.

```hcl
resource "huaweicloud_kps_keypair" "test" {
    ...

  lifecycle {
    ignore_changes = [
      encryption_type, kms_key_id, kms_key_name, private_key
    ]
  }
}
```
