---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_internet_bandwidths"
description: ""
---

# huaweicloud_global_internet_bandwidths

Use this data source to get a list of global internet bandwidths.

## Example Usage

### Get all global internet bandwidths

```hcl
data "huaweicloud_global_internet_bandwidths" "all" {}
```

### Get specific global internet bandwidths

```hcl
data "huaweicloud_global_internet_bandwidths" "test" {
  access_site = "cn-south-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_id` - (Optional, String) Specifies the global internet bandwidth ID.

* `name` - (Optional, String) Specifies the global internet bandwidth name.

* `size` - (Optional, String) Specifies the global internet bandwidth size.

* `access_site` - (Optional, String) Specifies the access site to which the global internet bandwidth belongs.

* `type` - (Optional, String) Specifies the global internet bandwidth type.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the global internet bandwidth.

* `status` - (Optional, String) Specifies the global internet bandwidth status. Valid values are **freezed** and **normal**.

* `tags` - (Optional, Map) Specifies the global internet bandwidth tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `internet_bandwidths` - The global internet bandwidths list.
  The [internet_bandwidths](#attrblock--internet_bandwidths) structure is documented below.

<a name="attrblock--internet_bandwidths"></a>
The `internet_bandwidths` block supports:

* `id` - The global internet bandwidth ID.

* `access_site` - The access site of the global internet bandwidth.

* `isp` - The internet service provider of the global internet bandwidth.

* `type` - The global internet bandwidth type.

* `charge_mode` - The charge mode of the global internet bandwidth.

* `size` - The global internet bandwidth size.

* `name` - The global internet bandwidth name.

* `enterprise_project_id` - The enterprise project ID of the global internet bandwidth.

* `ingress_size` - The ingress size of the global internet bandwidth.

* `ratio_95peak` - The enhanced 95% guaranteed rate of the global internet bandwidth.

* `frozen_info` - The frozen info of the global internet bandwidth.

* `status` - The status of the global internet bandwidth.

* `tags` - The tags of the global internet bandwidth.

* `description` - The description of the global internet bandwidth.

* `created_at` - The create time of the global internet bandwidth.

* `updated_at` - The update time of the global internet bandwidth.
