---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_instance"
description: ""
---

# huaweicloud_er_instance

Manages an ER instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "router_name" {}
variable "bgp_as_number" {}
variable "availability_zones" {
  type = list(string)
}

resource "huaweicloud_er_instance" "test" {
  availability_zones = var.availability_zones

  name = var.router_name
  asn  = var.bgp_as_number
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The router name.  
  The name can contain `1` to `64` characters, only English letters, Chinese characters, digits, underscore (_),
  hyphens (-) and dots (.) allowed.

* `availability_zones` - (Required, List) The availability zone list where the ER instance is located.
  The maximum number of availability zone is two. Select two AZs to configure active-active deployment for high
  availability which will ensure reliability and disaster recovery.

* `asn` - (Required, Int, ForceNew) The BGP AS number of the ER instance.  
  The valid value is range from `64,512` to `65534` or range from `4,200,000,000` to `4,294,967,294`.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description of the ER instance.  
  The description contain a maximum of `255` characters, and the angle brackets (< and >) are not allowed.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the ER instance
belongs.

* `enable_default_propagation` - (Optional, Bool) Whether to enable the propagation of the default route table.  
  The default value is **false**.

* `enable_default_association` - (Optional, Bool) Whether to enable the association of the default route table.  
  The default value is **false**.

* `auto_accept_shared_attachments` - (Optional, Bool) Whether to automatically accept the creation of shared
attachment.
  The default value is **false**.

* `default_propagation_route_table_id` - (Optional, String) The ID of the default propagation route table.

  -> Before modifying the default routing table of this instance (except cancel the default route table), make sure
     `enable_default_propagation` is set to **true**.

* `default_association_route_table_id` - (Optional, String) The ID of the default association route table.

  -> Before modifying the default routing table of this instance (except cancel the default route table), make sure
     `enable_default_association` is set to **true**.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Current status of the router.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

The router instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_er_instance.test 0ce123456a00f2591fabc00385ff1234
```
