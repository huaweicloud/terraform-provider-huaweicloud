---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_access_key"
description: |-
  Manages a permanent Access Key resource within HuaweiCloud IAM service.
---

# huaweicloud_identity_access_key

Manages a permanent Access Key resource within HuaweiCloud IAM service.

-> **NOTE:** You *must* have admin privileges in your HuaweiCloud cloud to use this resource.

## Example Usage

### Create an access key with a custom name in current execution directory

```hcl
variable "user_id" {}

resource "huaweicloud_identity_access_key" "test" {
  user_id     = var.user_id
  secret_file = abspath("./credentials.csv")
}
```

### Create an access key with the default name in current execution directory

```hcl
variable "user_id" {}

resource "huaweicloud_identity_access_key" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String, ForceNew) Specifies the IAM user ID for which access key (AK/SK) to be created.  
  Changing this creates a new resource.

* `description` - (Optional, String) Specifies the description of the access key.

* `status` - (Optional, String) Specifies the status of the access key.  
  The valid values are as follows:
  + **active**
  + **inactive**

  Defaults to **active**.

* `secret_file` - (Optional, String, ForceNew) Specifies the file name of the credentials (CSV) that can save access
  key and access secret key.  
  Defaults to **./credentials-{{user name}}.csv**. Changing this creates a new resource.

* `pgp_key` - (Optional, String, ForceNew) Specifies the PGP public key (base64 encoded) used to encrypt the storaged
  secret key.  
  Changing this creates a new resource.

  -> This input can also be the Keybase username, in the format: `keybase:some_person_that_exists`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The access key ID.

* `secret` - The access secret key.

  -> Set this value only if the credentials file fails to save, which usually occurs when the directory does not exist
     or access is denied.

* `key_fingerprint` - The fingerprint of the PGP key used to encrypt the secret.

* `encrypted_secret` - The encrypted secret, which encoded in base64.  
  The encrypted secret may be decrypted using the command line, for example:
  `terraform output encrypted_secret | base64 --decode | keybase pgp decrypt`.

* `user_name` - The name of IAM user.

* `create_time` - The creation time of the access key, in [ISO-8601](https://www.iso.org/iso-8601-date-and-time-format.html)
  UTC format.
