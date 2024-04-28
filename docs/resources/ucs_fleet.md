---
subcategory: "Ubiquitous Cloud Native Service (UCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ucs_fleet"
description: ""
---

# huaweicloud_ucs_fleet

Manages a UCS fleet resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
resource "huaweicloud_ucs_fleet" "test" {
  name        = "fleet_1"
  description = "created by terraform"
}
```

### Fleet with Permissions

```hcl
variable "policy_id_1" {}
variable "policy_id_2" {}
variable "policy_id_3" {}

resource "huaweicloud_ucs_fleet" "test" {
  name        = "fleet_1"
  description = "created by terraform"

  permissions {
    namespaces = ["*"]
    policy_ids = [var.policy_id_1]
  }

  permissions {
    namespaces = ["default", "kube-system"]
    policy_ids = [var.policy_id_2,var.policy_id_3]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the UCS fleet.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the UCS fleet.

* `permissions` - (Optional, List) Specifies the permissions of the UCS fleet. The structure is documented below.

The `permissions` block supports:

* `namespaces` - (Optional, List) Specifies the list of namespaces.
  The elements can be: **\***, **default**, **kube-system** and **kube-public**.

* `policy_ids` - (Optional, List) Specifies the list of policy IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `cluster_ids` - Indicates the list of cluster IDs to add to the UCS fleet.

## Import

The UCS fleet can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ucs_fleet.test dbd042ec-2474-11ee-9d1c-0255ac1000ce
```
