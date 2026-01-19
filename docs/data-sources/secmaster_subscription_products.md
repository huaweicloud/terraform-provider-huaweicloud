---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_subscription_products"
description: |-
  Use this data source to get the list of SecMaster subscription products.
---

# huaweicloud_secmaster_subscription_products

Use this data source to get the list of SecMaster subscription products.

## Example Usage

```hcl
data "huaweicloud_secmaster_subscription_products" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `basic` - The basic product code object.

  The [basic](#code_object_struct) structure is documented below.

* `standard` - The standard product code object.

  The [standard](#code_object_struct) structure is documented below.

* `professional` - The professional product code object.

  The [professional](#code_object_struct) structure is documented below.

* `large_screen` - The large screen product code object.

  The [large_screen](#code_object_struct) structure is documented below.

* `log_collection` - The log collection product code object.

  The [log_collection](#code_object_struct) structure is documented below.

* `log_retention` - The log retention product code object.

  The [log_retention](#code_object_struct) structure is documented below.

* `log_analysis` - The log analysis product code object.

  The [log_analysis](#code_object_struct) structure is documented below.

* `soar` - The SOAR product code object.

  The [soar](#code_object_struct) structure is documented below.

<a name="code_object_struct"></a>
The `basic`, `standard`, `professional`, `large_screen`, `log_collection`, `log_retention`, `log_analysis`,
`soar` block supports:

* `cloud_service_type` - The primary service type for cloud service products. The default value is: **hws.service.type.sa**.

* `resource_type` - The resource type of the subscription product.

* `resource_spec_code` - The resource specification code of the subscription product.

* `resource_size_measure_id` - The resource size measurement unit ID of the subscription product.

* `usage_factor` - The usage factor of the subscription product. The value must match the usage factor in the call detail
  record (CDR). The correspondence between cloud services and usage factors is as follows:
  + **duration**: Time, primarily for major versions (basic, standard, professional)
  + **count**: Number of times, primarily for security orchestration
  + **flow**: Traffic, primarily for log analysis and collection
  + **retention**: Retention, primarily for log retention

* `usage_measure_id` - The usage measurement unit ID of the subscription product.
  For example, for hourly pricing, the usage value is `1`, and the unit is hours. Enumerated values ​​are as follows:
  + `4`: Hour
  + `10`: GB (Bandwidth used for traffic pricing)
  + `11`: MB (Bandwidth used for traffic pricing)
  + `13`: Byte (Bandwidth used for traffic pricing)

* `region_id` - The region ID of the subscription product. Value **null** means the current region's encoding.
