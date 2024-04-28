---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_access_sites"
description: ""
---

# huaweicloud_global_eip_access_sites

Use this data source to get a list of global EIP access sites.

## Example Usage

### Get all global EIP access sites

```hcl
data "huaweicloud_global_eip_access_sites" "all" {}
```

### Get specific global EIP access sites through proxy region

```hcl
data "huaweicloud_global_eip_access_sites" "test" {
  proxy_region = "cn-south-1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the name of the access sites.

* `proxy_region` - (Optional, String) Specifies the region ID where the pop site is hosted.

* `iec_az_code` - (Optional, String) Specifies the availability zone code of the edge site.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `access_sites` - The access sites list.
  The [access_sites](#attrblock--access_sites) structure is documented below.

<a name="attrblock--access_sites"></a>
The `access_sites` block supports:

* `id` - The ID of the access site.

* `proxy_region` - The region ID where the pop site is hosted.

* `iec_az_code` - The availability zone code of the edge site

* `name` - The name of the access site.

* `cn_name` - The Chinese name of the access site.

* `en_name` - The English name of the access site.

* `created_at` - The create time of the access site.

* `updated_at` - The update time of the access site.
