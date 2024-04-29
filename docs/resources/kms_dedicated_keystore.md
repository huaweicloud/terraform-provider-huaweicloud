---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_dedicated_keystore"
description: ""
---

# huaweicloud_kms_dedicated_keystore

Manages a KMS (Key Management Service) dedicated keystore resource within HuaweiCloud.

## Example Usage

```hcl
variable "hsm_cluster_id" {}
variable "hsm_ca_cert" {}

resource "huaweicloud_kms_dedicated_keystore" "test" {
  alias          = "test_name"
  hsm_cluster_id = var.hsm_cluster_id
  hsm_ca_cert    = var.hsm_ca_cert
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `alias` - (Required, String, ForceNew) Specifies the alias of a dedicated keystore. The valid length is limited from
  `1` to `255`. Only digits, uppercase letters, lowercase letters, underscores (_), hyphens (-), colons (:) and
  forward slashes (/) are allowed.

  Changing this parameter will create a new resource.

* `hsm_cluster_id` - (Required, String, ForceNew) Specifies the ID of a dedicated HSM cluster that has no dedicated keystore.
  Changing this parameter will create a new resource.

  -> The dedicated HSM cluster must be active. The cluster can be activated only after adding the master cryptographic
  machine and non-master cryptographic machine. Currently, only some encryption machine models are supported. Support
  work order consultation for details.

* `hsm_ca_cert` - (Required, String, ForceNew) Specifies the CA certificate of the dedicated HSM cluster.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The KMS dedicated keystore can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_kms_dedicated_keystore.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `hsm_ca_cert`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_kms_dedicated_keystore" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      hsm_ca_cert,
    ]
  }
}
```
