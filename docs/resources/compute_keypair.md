---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_keypair"
description: ""
---

# huaweicloud_compute_keypair

!> **WARNING:** It has been deprecated, use `huaweicloud_kps_keypair` instead.

Manages a keypair resource within HuaweiCloud.

## Example Usage

### Create a new keypair and export private key to current folder

```hcl
resource "huaweicloud_compute_keypair" "test-keypair" {
  name     = "my-keypair"
  key_file = "private_key.pem"
}
```

### Import an exist keypair

```hcl
resource "huaweicloud_compute_keypair" "test-keypair" {
  name       = "my-keypair"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAlJq5Pu+eizhou7nFFDxXofr2ySF8k/yuA9OnJdVF9Fbf85Z59CWNZBvcAT... root@terra-dev"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the keypair resource. If omitted, the
  provider-level region will be used. Changing this creates a new keypair.

* `name` - (Required, String, ForceNew) Specifies a unique name for the keypair. Changing this creates a new keypair.

* `public_key` - (Optional, String, ForceNew) Specifies the imported OpenSSH-formatted public key. Changing this creates
  a new keypair.
  This parameter and `key_file` are alternative.

* `key_file` - (Optional, String, ForceNew) Specifies the path of the created private key.
  The private key file (**.pem**) is created only after the resource is created.
  By default, the private key file will be created in the same folder as the current script file.
  If you need to create it in another folder, please specify the path for `key_file`.
  Changing this creates a new keypair.

  ~>**NOTE:** If the private key file already exists, it will be overwritten after a new keypair is created.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Keypairs can be imported using the `name`, e.g.

```bash
$ terraform import huaweicloud_compute_keypair.my-keypair test-keypair
```
