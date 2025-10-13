---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_ip_information"
description: |-
  Use this data source to get the IP attribution information of CDN nodes within HuaweiCloud.  
---

# huaweicloud_cdn_ip_information

Use this data source to get the IP attribution information of CDN nodes within HuaweiCloud.

## Example Usage

```hcl
variable "ip_list_str" {}

data "huaweicloud_cdn_ip_information" "test" {
  ips = var.ip_list_str
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the queried IP attribution information are located.

* `ips` - (Required, String) Specifies the list of IP addresses to be queried.  
  The maximum number of IPs that can be queried is 20, and multiple IPs are separated by commas (,).

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the IP attribution
  information belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `information` - The list of IP attribution information that matched filter parameters.  
  The [information](#cdn_information) structure is documented below.

<a name="cdn_information"></a>
The `information` block supports:

* `ip` - The IP address to be queried.

* `belongs` - Whether the IP belongs to CDN nodes.  
  The valid values are as follows:
  + **true**
  + **false**

* `region` - The province where the IP is located.  
  "Unknown" indicates that the attribution is unknown.

* `isp` - The ISP name.  
  If the IP attribution is unknown, this field returns null.

* `platform` - The platform name.  
  If the IP attribution is unknown, this field returns null.
