---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_policy_associate"
description: ""
---

# huaweicloud_cnad_advanced_policy_associate

Manages a CNAD advanced policy associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "policy_id" {}
variable "instance_id" {}
variable "protected_object_ids" {
  type = list(string)
}

resource "huaweicloud_cnad_advanced_policy_associate" "test" {
  policy_id            = var.policy_id
  instance_id          = var.instance_id
  protected_object_ids = var.protected_object_ids
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID. You can find it through data source
  `huaweicloud_cnad_advanced_instances`.

  Changing this parameter will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the CNAD advanced policy ID in which to associate protected
  objects.

  Changing this parameter will create a new resource.

* `protected_object_ids` - (Required, List) Specifies the protected object IDs to associate. The protected object must
  have no binding policy. You can find it through data source `huaweicloud_cnad_advanced_protected_objects`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `protected_objects` - The advanced protected objects which the policy associate.
  The [protected_objects](#advancedProtectedObject_GetAdvancedProtectedObjects) structure is documented below.

<a name="advancedProtectedObject_GetAdvancedProtectedObjects"></a>
The `protected_objects` block supports:

* `id` - Indicates the ID of the protected object.

* `ip_address` - Indicates the IP of the protected object.

* `type` - Indicates the type of the protected object.

* `name` - Indicates the name of the protected object.

* `instance_id` - Indicates the instance ID of the protected object.

* `instance_name` - Indicates the instance name of the protected object.

* `instance_version` - Indicates the instance version of the protected object.
  + **cnad_pro**: Professional Edition.
  + **cnad_ip**: Standard Edition.
  + **cnad_ep**: Platinum Edition.
  + **cnad_full_high**: Unlimited Protection Advanced Edition.
  + **cnad_vic**: On demand Version.
  + **cnad_intl_ep**: International Station Platinum Edition.

* `region` - Indicates the region to which the protected object belongs.

* `status` - Indicates the status of the protected object.
  + **0**: Normal.
  + **1**: Cleaning in progress.
  + **2**: In a black hole.

* `block_threshold` - Indicates the block threshold of the protected object.

* `clean_threshold` - Indicates the clean threshold of the protected object.

* `policy_name` - Indicates the policy name of the protected object.

## Import

The CNAD advanced policy associate can be imported using the `policy_id` and `instance_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cnad_advanced_policy_associate.test <policy_id>/<instance_id>
```
