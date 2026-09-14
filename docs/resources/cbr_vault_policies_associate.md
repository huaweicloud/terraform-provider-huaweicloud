---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_vault_policies_associate"
description: |-
  Manages a CBR vault policies associate resource within HuaweiCloud.
---

# huaweicloud_cbr_vault_policies_associate

Manages a CBR vault policies associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "vault_id" {}
variable "backup_policy_id" {}
variable "replication_policy_id" {}
variable "destination_vault_id" {}

resource "huaweicloud_cbr_vault_policies_associate" "test" {
  vault_id = var.vault_id

  policies {
    id = var.backup_policy_id
  }
  policies {
    id                   = var.replication_policy_id
    destination_vault_id = var.destination_vault_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted,
  the provider-level region will be used. Changing this will create new resource.

* `vault_id` - (Required, String, NonUpdatable) Specifies the ID of the CBR vault to which the policies will be
  associated.

* `policies` - (Required, List) Specifies the policy details to associate with the CBR vault.
  The [policies](#cbr_vault_policies_associate_policies) structure is documented below.

<a name="cbr_vault_policies_associate_policies"></a>
The `policies` block supports:

* `id` - (Required, String) Specifies the ID of the policy to associate with the vault.

* `destination_vault_id` - (Optional, String) Specifies the destination vault ID when associating a replication policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the vault ID.

## Import

The CBR vault policies associate can be imported using the `id` (vault ID), e.g.

```bash
$ terraform import huaweicloud_cbr_vault_policies_associate.test <id>
```
