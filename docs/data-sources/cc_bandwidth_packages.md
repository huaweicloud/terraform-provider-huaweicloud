---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_packages"
description: ""
---

# huaweicloud_cc_bandwidth_packages

Use this data source to get the list of CC bandwidth packages.

## Example Usage

```hcl
variable "bandwidth_package_id" {}

data "huaweicloud_cc_bandwidth_packages" "test" {
  bandwidth_package_id = var.bandwidth_package_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `bandwidth_package_id` - (Optional, String) Specifies the bandwidth package ID.

* `name` - (Optional, String) Specifies the bandwidth package name.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project that the bandwidth package
  belongs to.

* `status` - (Optional, String) Specifies the bandwidth package status.
  The valid value is as follows:
  + **ACTIVE**: Bandwidth packages are available.

* `billing_mode` - (Optional, String) Specifies the billing mode of the bandwidth package.
  The options are as follows:
  + **1**：pay by period for the Chinese Mainland website.
  + **2**: pay by period for the International website.
  + **3**: pay-per-use for the Chinese Mainland website.
  + **4**: pay-per-use for the International website.
  + **5**: 95th percentile bandwidth billing for the Chinese Mainland website.
  + **6**: 95th percentile bandwidth billing for the International website.

* `resource_id` - (Optional, String) Specifies the ID of the resource that the bandwidth package is bound to.

* `bandwidth` - (Optional, Int) Specifies the bandwidth range specified for the bandwidth package.

* `tags` - (Optional, Map) Specifies the bandwidth package tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidth_packages` - Bandwidth package list.

  The [bandwidth_packages](#bandwidth_packages_struct) structure is documented below.

<a name="bandwidth_packages_struct"></a>
The `bandwidth_packages` block supports:

* `id` - The bandwidth package ID.

* `name` - The bandwidth package name.

* `description` - The bandwidth package description.

* `domain_id` - The ID of the account that the bandwidth package belongs to.

* `enterprise_project_id` - The ID of the enterprise project that the bandwidth package belongs to.

* `project_id` - Project ID of the bandwidth package.

* `created_at` - Time when the resource was created.

* `updated_at` - Time when the resource was updated.

* `resource_id` - The ID of the resource that the bandwidth package is bound to.

* `resource_type` - Type of the resource that the bandwidth package is bound to.

* `local_area_id` - The ID of a local access point.

* `remote_area_id` - The ID of a remote access point.

* `spec_code` - Specification code of the bandwidth package.

* `billing_mode` - Billing mode of the bandwidth package.

* `tags` - The bandwidth package tags.

* `status` - Status of the bandwidth package.

* `order_id` - Order ID of the bandwidth package.

* `product_id` - Product ID of the bandwidth package.

* `charge_mode` - Billing option.

* `bandwidth` - Bandwidth range specified for the bandwidth package.

* `interflow_mode` - Interflow mode of the bandwidth package.
