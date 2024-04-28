---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_protected_objects"
description: ""
---

# huaweicloud_cnad_advanced_protected_objects

Use this data source to get the list of CNAD advanced protected objects.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_cnad_advanced_protected_objects" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, String) Specifies the CNAD advanced instance ID. You can find it through data source
  `huaweicloud_cnad_advanced_instances`.

* `policy_id` - (Optional, String) Specifies the CNAD advanced policy ID.

* `ip_address` - (Optional, String) Specifies the CNAD advanced protected object IP which you want to query.
  The IP must be the Elastic IP.

* `protected_object_id` - (Optional, String) Specifies the CNAD advanced protected object ID which you want to query.
  The protected object ID must be the ID of the Elastic IP.

* `type` - (Optional, String) Specifies the type of the protected object. This field means the type of resource which
  bound to Elastic IP. Valid values are: **VPN**, **NAT**, **VIP**,**CCI**, **EIP**, **ELB**, **REROUTING_IP**,
  **SubEni** and **NetInterFace**.

* `is_unbound` - (Optional, Bool) Specifies whether query protected objects which policies unbound.
  If the field is set to **true**, only protection objects that can bind policies will be returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `protected_objects` - Indicates the list of the advanced protected objects.
  The [protected_objects](#GetProtectedObjectsResponseBody_protected_objects) structure is documented below.

<a name="GetProtectedObjectsResponseBody_protected_objects"></a>
The `protected_objects` block supports:

* `id` - Indicates the ID of the protected object.

* `ip_address` - Indicates the IP of the protected object.

* `type` - Indicates the type of the protected object.

* `name` - Indicates the name of the protected object.

* `instance_id` - Indicates the instance ID of the protected object.

* `instance_name` - Indicates the instance name of the protected object.

* `instance_version` - Indicates the instance version of the protected object.

* `region` - Indicates the region to which the protected object belongs.

* `status` - Indicates the status of the protected object.

* `block_threshold` - Indicates the block threshold of the protected object.

* `clean_threshold` - Indicates the clean threshold of the protected object.

* `policy_name` - Indicates the policy name of the protected object.
