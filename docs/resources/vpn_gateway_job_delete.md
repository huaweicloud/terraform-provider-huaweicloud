---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_gateway_job_delete"
description: |-
  Manages a VPN gateway job delete resource within HuaweiCloud.
---

# huaweicloud_vpn_gateway_job_delete

Manages a VPN gateway job delete resource within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_vpn_gateway_job_delete" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the gateway job ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
