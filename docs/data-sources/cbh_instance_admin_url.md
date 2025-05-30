---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_instance_admin_url"
description: |-
  Use this data source to get the URL for logging in to a CBH instance as user admin within HuaweiCloud.
---

# huaweicloud_cbh_instance_admin_url

Use this data source to get the URL for logging in to a CBH instance as user admin within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_cbh_instance_admin_url" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the CBH instance ID, in UUID format.  
  For details about how to obtain the value, see "Viewing CBH Instance Details" in Cloud Bastion Host User Guide.  
  [reference](https://support.huaweicloud.com/intl/en-us/usermanual-cbh/cbh_02_0043.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, same as `server_id`.

* `admin_url` - The URL for logging in to a CBH instance as user admin.
