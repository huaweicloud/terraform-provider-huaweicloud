---
subcategory: Cloud Connect (CC)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_connection"
description: ""
---

# huaweicloud_cc_connection

Manages a Cloud Connection resource within HuaweiCloud.

Cloud Connect (CC) is a service that enables you to quickly build ultra-fast, high-quality, and stable networks
between VPCs across regions and between VPCs and on-premises data centers.

## Example Usage

```hcl
resource "huaweicloud_cc_connection" "test" {
   name = "connection_demo"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The cloud connection name.  
  The name can contain `1` to `64` characters, only English letters, Chinese characters, digits, hyphens (-),
  underscores (_) and dots (.) are allowed.

* `description` - (Optional, String) The Description about the cloud connection.  
  The description contain a maximum of `255` characters, and the angle brackets (< and >) are not allowed.

* `enterprise_project_id` - (Optional, String) The enterprise project id of the cloud connection.  
  Value 0 indicates the default enterprise project.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the cloud connection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `domain_id` - The Domain ID.

* `status` - The status of the cloud connection.  
  The options are as follows:
    + **ACTIVE**: Device deleted.

* `used_scene` - The Scenario.  
  The options are as follows:
    + **vpc**: VPCs or virtual gateways can use this cloud connection.

* `network_instance_number` - The number of network instances associated with the cloud connection instance.

* `bandwidth_package_number` - The number of bandwidth packages associated with the cloud connection instance.

* `inter_region_bandwidth_number` - The number of inter-domain bandwidths associated with the cloud connection instance.

## Import

The cloud connection can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cc_connection.test 0ce123456a00f2591fabc00385ff1234
```
