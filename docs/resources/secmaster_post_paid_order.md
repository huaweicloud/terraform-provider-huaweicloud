---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_post_paid_order"
description: |-
  Manages a SecMaster post paid order resource within HuaweiCloud.
---

# huaweicloud_secmaster_post_paid_order

Manages a SecMaster post paid order resource within HuaweiCloud.

-> Destroying this resource will not change the status of the post paid order resource.

## Example Usage

### Basic Example

```hcl
resource "huaweicloud_secmaster_post_paid_order" "test" {
  product_list {
    id                 = "CA696738-FD15-47C1-A389-CD0B34415055"
    product_id         = "OFFI908269345109094402"
    cloud_service_type = "hws.service.type.sa"
    resource_type      = "hws.resource.type.secmaster.typical"
    resource_spec_code = "secmaster.professional"
    usage_measure_id   = 4
    usage_value        = 1
    resource_size      = 4
    usage_factor       = "duration"    
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `operate_type` - (Optional, String, NonUpdatable) Specifies the operate type.
  The value can be: **create** and **addition**.

* `tags` - (Optional, Map, NonUpdatable) Specifies the tags of the resource in key/pair format.

* `product_list` - (Optional, List, NonUpdatable) Specifies the product list.
  The [object](#product_list) structure is documented below.

<a name="product_list"></a>
The `product_list` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the identifier, which must be unique.

* `product_id` - (Required, String, NonUpdatable) Specifies the offering ID,
  which is obtained from the CBC price inquiry.

* `cloud_service_type` - (Required, String, NonUpdatable) Specifies the cloud service type.
  The fixed value is **hws.service.type.sa**.

* `resource_type` - (Required, String, NonUpdatable) Specifies the resource type of the purchased product.
  For example, the resource type for typical scenarios in SecMaster is **hws.resource.type.secmaster.typical**.

* `resource_spec_code` - (Required, String, NonUpdatable) Specifies the resource specifications of the purchased
  product. For example, the resource specification for the basic edition in SecMaster is **secmaster.basic**.

* `usage_measure_id` - (Required, Int, NonUpdatable) Specifies the usage measurement unit.
  For example, the resources are billed by hour, the usage value is **1**, and the usage measurement unit is hour.
  The options are:

  + **4**: Hours;
  + **10**: GB, The bandwidth usage is measured by traffic (GB);
  + **11**: MB, The bandwidth usage is measured by traffic (MB);

* `usage_value` - (Required, Int, NonUpdatable) Specifies the usage value.

* `resource_size` - (Required, Int, NonUpdatable) Specifies the number of quotas.

* `usage_factor` - (Required, String, NonUpdatable) Specifies the usage factor.

* `resource_id` - (Optional, String, NonUpdatable) Specifies the resource ID,
  which is transferred only when the quota is added.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
