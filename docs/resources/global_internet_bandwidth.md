---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_internet_bandwidth"
description: ""
---

# huaweicloud_global_internet_bandwidth

Manages a global internet bandwidth resource within HuaweiCloud.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "bandwidth_name" {}

data "huaweicloud_global_eip_pools" "all" {}

resource "huaweicloud_global_internet_bandwidth" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  charge_mode           = "95peak_guar"
  enterprise_project_id = var.enterprise_project_id
  size                  = 300
  isp                   = data.huaweicloud_global_eip_pools.all.geip_pools[0].isp
  name                  = var.bandwidth_name
  type                  = data.huaweicloud_global_eip_pools.all.geip_pools[0].allowed_bandwidth_types[0].type

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `access_site` - (Required, String, ForceNew) Specifies the access site name.
  Changing this creates a new resource.

* `charge_mode` - (Required, String) Specifies the charge mode.
  Value can be as follows:
  + **bandwidth**: billed by bandwidth
  + **traffic**: billed by traffic
  + **95peak_plus_1000**: billed by enhanced 95
  + **95peak_bidirection**: billed by bidirectional traditional 95
  + **95peak_guar**: billed by traditional 95 with guaranteed minimum
  + **95peak_avr**: 95 peak average billing

* `isp` - (Required, String, ForceNew) Specifies the internet service provider of the global internet bandwidth.
  Changing this creates a new resource.

* `size` - (Required, Int) Specifies the size of the global internet bandwidth.
  The value ranges from **300 Mbit/s** to **5000 Mbit/s** in normal.

* `ingress_size` - (Optional, Int) Specifies the ingress size of the global internet bandwidth.
  It's **not** used for charge mode **95peak_guar**.

* `type` - (Optional, String) Specifies the type of the global internet bandwidth.

* `tags` - (Optional, Map) Specifies the tags of the global internet bandwidth.

* `name` - (Optional, String) Specifies the name of the global internet bandwidth.

* `description` - (Optional, String) Specifies the description of the global internet bandwidth.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id to which the global
  internet bandwidth belongs. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `frozen_info` - The frozen info of the global internet bandwidth.

* `ratio_95peak` - The enhanced 95% guaranteed rate of the global internet bandwidth.

* `status` - The status of the global internet bandwidth.

* `created_at` - The create time of the global internet bandwidth.

* `updated_at` - The update time of the global internet bandwidth.

## Import

The global internet bandwidth can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_global_internet_bandwidth.test <id>
```
