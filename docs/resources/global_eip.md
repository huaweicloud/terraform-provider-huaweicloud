---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip"
description: ""
---

# huaweicloud_global_eip

Manages a global EIP resource within HuaweiCloud.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "bandwidth_name" {}
variable "eip_name" {}

data "huaweicloud_global_eip_pools" "all" {}

resource "huaweicloud_global_internet_bandwidth" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  charge_mode           = "95peak_guar"
  enterprise_project_id = var.enterprise_project_id
  size                  = 300
  isp                   = data.huaweicloud_global_eip_pools.all.geip_pools[0].isp
  name                  = var.bandwidth_name
  type                  = data.huaweicloud_global_eip_pools.all.geip_pools[0].allowed_bandwidth_types[0].type
}

resource "huaweicloud_global_eip" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  enterprise_project_id = var.enterprise_project_id
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  internet_bandwidth_id = huaweicloud_global_internet_bandwidth.test.id
  name                  = var.eip_name

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `access_site` - (Required, String, ForceNew) Specifies the access site name.
  Changing this creates a new resource.

* `geip_pool_name` - (Required, String, ForceNew) Specifies the global EIP pool name.
  Changing this creates a new resource.

* `internet_bandwidth_id` - (Required, String, ForceNew) Specifies the internet bandwidth id which the global EIP use.
  Changing this creates a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id to which the global EIP
  belongs. Changing this creates a new resource.

* `description` - (Optional, String) Specifies the description of the global EIP.

* `name` - (Optional, String) Specifies the name of the global EIP.

* `tags` - (Optional, Map) Specifies the tags of the global EIP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `ip_address` - The ip address of the global EIP.

* `ip_version` - The ip version of the global EIP.

* `isp` - The the internet service provider of the global EIP.

* `global_connection_bandwidth_id` - The ID of the global connection bandwidth.

* `global_connection_bandwidth_type` - The type of the global connection bandwidth.

* `associate_instance_region` - The region of the associate instance.

* `associate_instance_id` - The ID of the associate instance.

* `associate_instance_type` - The type of the associate instance.

* `frozen` - The global EIP is frozen or not.

* `frozen_info` - The frozen info of the global EIP.

* `polluted` - The global EIP is polluted or not.

* `status` - The status of the global EIP.

* `created_at` - The create time of the global EIP.

* `updated_at` - The update time of the global EIP.

## Import

The global EIP can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_global_eip.test <id>
```
