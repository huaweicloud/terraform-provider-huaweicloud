---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_keypair"
description: |-
  Manages a keypair resource within HuaweiCloud.
---

# huaweicloud_kps_keypair

Manages a keypair resource within HuaweiCloud.

By default, keypair use the SSH-2 (RSA, 2048) algorithm for encryption and decryption.

Keys imported support the following cryptographic algorithms:

 * RSA-1024
 * RSA-2048
 * RSA-4096

## Example Usage

### Create a new KPS keypair

```hcl
variable "kms_key_id" {}
variable "kms_key_name" {}
variable "key_file" {}

resource "huaweicloud_kps_keypair" "test" {
  name            = "test-name"
  scope           = "user"
  encryption_type = "kms"
  kms_key_id      = var.kms_key_id
  kms_key_name    = var.kms_key_name
  description     = "test description"
  key_file        = var.key_file
}
```

### Import an existing KPS keypair

```hcl
variable "public_key" {}
variable "private_key" {}

resource "huaweicloud_kps_keypair" "test" {
  name            = "test-name"
  scope           = "account"
  encryption_type = "default"
  description     = "test description"
  public_key      = var.public_key
  private_key     = var.private_key
}
```

### Import an existing KPS keypair without private key

```hcl
variable "public_key" {}

resource "huaweicloud_kps_keypair" "test" {
  name        = "test-name"
  scope       = "account"
  description = "test description"
  public_key  = var.public_key
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the keypair resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies a unique name for the keypair. The name can contain a maximum of `64`
  characters, including letters, digits, underscores (_) and hyphens (-).
  Changing this parameter will create a new resource.

* `scope` - (Optional, String, ForceNew) Specifies the scope of keypair. The options are as follows:
  + **account**: Tenant-level, available to all users under the same account.
  + **user**: User-level, only available to user.

  Defaults to `user`. Changing this parameter will create a new resource.

* `user_id` - (Optional, String) Specifies the user ID to which the keypair belongs.

  -> 1. If the `scope` set to **user**, this parameter value must be the ID of the user who creates the resource.
  <br/>2. Due to API restrictions, `private_key` and `encryption_type` must be configured when editing this field.

* `encryption_type` - (Optional, String) Specifies encryption mode. The options are as follows:
  + **default**: The default encryption mode. Applicable to sites where KMS is not deployed.
  + **kms**: KMS encryption mode.

  -> 1. Please configure this field to **default** if the KMS service is not available at the site.
  <br/>2. Due to API restrictions, `private_key` must be configured when editing this field.

* `kms_key_id` - (Optional, String) Specifies the KMS key ID to encrypt private keys.

* `kms_key_name` - (Optional, String) Specifies the KMS key name to encrypt private keys.

-> 1. At least one of `kms_key_id` or `kms_key_name` must be set when `encryption_type` is set to **kms**.
  <br/>2. Due to API restrictions, `private_key` and `encryption_type` must be configured when editing `kms_key_id` or
  `kms_key_name`.

* `description` - (Optional, String) Specifies the description of keypair.

* `public_key` - (Optional, String, ForceNew) Specifies the imported OpenSSH-formatted public key.
  It is required when import keypair. Changing this parameter will create a new resource.

* `private_key` - (Optional, String) Specifies the imported OpenSSH-formatted private key.

  -> 1. Setting this field to empty during editing will clear the private key.
  <br/>2. Due to API restrictions, `encryption_type` must be configured when configuring this field.

* `key_file` - (Optional, String, ForceNew) Specifies the path of the created private key.
  The private key file (**.pem**) is created only when creating a KPS keypair.
  Importing an existing keypair will not obtain the private key information.

  Changing this parameter will create a new resource.

  ->**NOTE:** If the private key file already exists, it will be overwritten after a new keypair is created.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to name.

* `created_at` - The keypair creation time.

* `fingerprint` - Fingerprint information about a keypair.

* `is_managed` - Whether the private key is managed by HuaweiCloud.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

Keypair can be imported using the `name`, e.g.

```bash
$ terraform import huaweicloud_kps_keypair.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `encryption_type`, `kms_key_id`,
`kms_key_name`, `key_file` and `private_key`. It is generally recommended running `terraform plan` after importing a keypair.
You can then decide if changes should be applied to the keypair, or the resource definition
should be updated to align with the keypair. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_kps_keypair" "test" {
    ...

  lifecycle {
    ignore_changes = [
      encryption_type, kms_key_id, kms_key_name, key_file, private_key
    ]
  }
}
```
