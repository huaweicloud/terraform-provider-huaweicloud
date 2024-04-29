---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_protected_object"
description: ""
---

# huaweicloud_cnad_advanced_protected_object

Manages a CNAD advanced protected object resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "object_id" {}
variable "object_ip_address" {}
variable "type" {}

resource "huaweicloud_cnad_advanced_protected_object" "test" {
  instance_id = var.instance_id
  
  protected_objects {
    id         = var.object_id
    ip_address = var.object_ip_address
    type       = var.type
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the CNAD advanced instance ID in which to bind protected
  objects. You can find it through data source `huaweicloud_cnad_advanced_instances`.

  Changing this parameter will create a new resource.

* `protected_objects` - (Required, List) Specifies the advanced protected objects.
  The [Protected_Object](#advancedProtectedObject_protected_objects) structure is documented below.

<a name="advancedProtectedObject_protected_objects"></a>
The `Protected_Object` block supports:

* `id` - (Required, String) Specifies the ID of the protected object. The field must be the ID of the Elastic IP,
  which in the same region with the CNAD advanced instance. You can find it through data source
  `huaweicloud_cnad_advanced_available_objects`.

* `ip_address` - (Required, String) Specifies the IP of the protected object. The field must be the IP of the Elastic
  IP, which in the same region with the CNAD advanced instance. You can find it through data source
  `huaweicloud_cnad_advanced_available_objects`. This field and `id` must belong to the same protected object.

* `type` - (Required, String) Specifies the type of the protected object. Valid values are **VPN**, **NAT**, **VIP**,
  **CCI**, **EIP**, **ELB**, **REROUTING_IP**, **SubEni** and **NetInterFace**. You can find it through data source
  `huaweicloud_cnad_advanced_available_objects`. This field and `id` must belong to the same protected object.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The CNAD advanced instance resource ID.

* `protected_objects` - The advanced protected objects.
  The [Protected_Object](#advancedProtectedObject_GetAdvancedProtectedObjects) structure is documented below.

<a name="advancedProtectedObject_GetAdvancedProtectedObjects"></a>
The `Protected_Object` block supports:

* `name` - The name of the protected object.

* `instance_id` - The instance ID which the protected object belongs to.

* `instance_name` - The instance name which the protected object belongs to.

* `instance_version` - The instance version which the protected object belongs to.
  + **cnad_pro**: Professional Edition.
  + **cnad_ip**: Standard Edition.
  + **cnad_ep**: Platinum Edition.
  + **cnad_full_high**: Unlimited Protection Advanced Edition.
  + **cnad_vic**: On demand Version.
  + **cnad_intl_ep**: International Station Platinum Edition.

* `region` - The region to which the protected object belongs.

* `status` - The status of the protected object.
  + **0**: Normal.
  + **1**: Cleaning in progress.
  + **2**: In a black hole.

* `block_threshold` - The blocking threshold of the protected object.

* `clean_threshold` - The cleaning threshold of the protected object.

* `policy_name` - The policy name which the protected object binding.

## Import

The CNAD advanced protected object can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cnad_advanced_protected_object.test <id>
```
