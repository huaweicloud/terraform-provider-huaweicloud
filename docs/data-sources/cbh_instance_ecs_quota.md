---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_instance_ecs_quota"
description: |-
  Use this data source to get the ECS quota required for creating CBH instance within HuaweiCloud.
---

# huaweicloud_cbh_instance_ecs_quota

Use this data source to get the ECS quota required for creating CBH instance within HuaweiCloud.

## Example Usage

```hcl
variable "resource_spec_code" {}
variable "availability_zone" {}

data "huaweicloud_cbh_instance_ecs_quota" "test" {
  availability_zone  = var.availability_zone
  resource_spec_code = var.resource_spec_code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Required, String) Specifies the availability zone name.

* `resource_spec_code` - (Required, String) Specifies the specification code of the CBH instance to be created.
  The valid values are as follows:  
  + **cbh.basic.10**: `10` asset standard version.  
  + **cbh.enhance.10**: `10` asset professional version.  
  
  Please refer to the API document link for its value
  [reference](https://support.huaweicloud.com/intl/en-us/api-cbh/ListSpecifications.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status_v6` - The status of the CBH instance specification resource which support Ipv6.
  The valid values are as follows:  
  + **sellout**: This specification of resources has been sold out.  
  + **normal**: This specification of resources are commercially available normally.

* `status` - The status of the CBH instance specification resource.
  The valid values are as follows:  
  + **sellout**: This specification of resources has been sold out.  
  + **normal**: This specification of resources are commercially available normally.
