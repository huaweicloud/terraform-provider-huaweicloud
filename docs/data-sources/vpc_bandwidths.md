---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_bandwidths"
description: ""
---

# huaweicloud_vpc_bandwidths

Use this data source to get a list of shared bandwidths.

## Example Usage

### Example Usage of getting all bandwidths

```hcl
data "huaweicloud_vpc_bandwidths" "all" {}
```

### Example Usage to filter specific bandwidths

```hcl
variable "bandwidth_name" {}

data "huaweicloud_vpc_bandwidths" "filter_by_name" {
  name = var.bandwidth_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the bandwidths.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the bandwidth.

* `size` - (Optional, Int) Specifies the size of the bandwidth.

* `bandwidth_id` - (Optional, String) Specifies the ID of the bandwidth.

* `charge_mode` - (Optional, String) Specifies the charge mode of the bandwidth.
  Possible values can be **bandwidth** and **95peak_plus**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `bandwidths` - The filtered bandwidths.
  The [bandwidths](#attrblock--bandwidths) structure is documented below.

<a name="attrblock--bandwidths"></a>
The `bandwidths` block supports:

* `id` - Indicates the ID of the bandwidth.

* `bandwidth_type` - Indicates the bandwidth type.

* `charge_mode` - Indicates the charge mode of the bandwidth.

* `enterprise_project_id` - Indicates the enterprise project id the bandwidth belongs to.

* `name` - Indicates the name of the bandwidth.

* `publicips` - An array of EIPs that use the bandwidth. The object includes the following:
  The [publicips](#attrblock--bandwidths--publicips) structure is documented below.

* `share_type` - Indicates whether the bandwidth is shared or dedicated.

* `size` - Indicates the size of the bandwidth.

* `status` - Indicates the status of the bandwidth.

* `created_at` - Indicates the create time of the bandwidth.

* `updated_at` - Indicates the update time of the bandwidth.

<a name="attrblock--bandwidths--publicips"></a>
The `publicips` block supports:

* `id` - The ID of the EIP or IPv6 port that uses the bandwidth.

* `ip_address` - The IPv4 or IPv6 address.

* `ip_version` - The IP version.

* `type` - The EIP type.
