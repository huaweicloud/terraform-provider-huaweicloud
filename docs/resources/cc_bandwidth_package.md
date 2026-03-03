---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package"
description: |-
  Manages a bandwidth package resource of Cloud connect within HuaweiCloud.
---

# huaweicloud_cc_bandwidth_package

Manages a bandwidth package resource of Cloud connect within HuaweiCloud.

## Example Usage

## Postpaid Example Usage

```hcl
variable "name" {}

resource "huaweicloud_cc_bandwidth_package" "test" {
  name           = var.name
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 5
  description    = "This is a demo"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Prepaid Example Usage

```hcl
variable "name" {}

resource "huaweicloud_cc_bandwidth_package" "test" {
  name           = var.name
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 1
  bandwidth      = 5
  description    = "This is a demo"

  prepaid_options {
    period_type   = "month"
    period_num    = 1
    is_auto_renew = true
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The bandwidth package name.

* `local_area_id` - (Required, String, NonUpdatable) Specifies the local area ID. Valid values are:
  + **Chinese-Mainland**: Chinese mainland
  + **Asia-Pacific**: Asia Pacific
  + **Africa**: Africa
  + **Western-Latin-America**: Western Latin America
  + **Eastern-Latin-America**: Eastern Latin America
  + **Northern-Latin-America**: Northern Latin America

* `remote_area_id` - (Required, String, NonUpdatable) Specifies the remote area ID. Valid values are:
  + **Chinese-Mainland**: Chinese mainland
  + **Asia-Pacific**: Asia Pacific
  + **Africa**: Africa
  + **Western-Latin-America**: Western Latin America
  + **Eastern-Latin-America**: Eastern Latin America
  + **Northern-Latin-America**: Northern Latin America

* `charge_mode` - (Required, String, NonUpdatable) Specifies the billing option of the bandwidth package.
  Valid value is **bandwidth**.

* `billing_mode` - (Required, String) Specifies the billing mode of the bandwidth package. Valid values are:
  + **1**: yearly/monthly on the Chinese mainland website
  + **2**: yearly/monthly on the International website
  + **3**: pay-per-use on the Chinese mainland website
  + **4**: pay-per-use on the International website
  + **5**: 95th percentile bandwidth billing on the Chinese mainland website
  + **6**: 95th percentile bandwidth billing on the International website

  -> Below is a detailed usage guide for the `billing_mode` field:
  <br/>1. The value **​​1** and **2** indicate instances of the prepaid type.
  <br/>2. The value **​​3**, **4**, **5** and **6** indicate instances of the postpaid type.
  <br/>3. This field can only be edited when its value is **3** or **4**; otherwise, editing it will result in an error.
  <br/>4. The value of this field can only be edited to **1**, **2**, **5**, or **6**.
  <br/>5. When the value of this field is configured as **1** or **2**, the `prepaid_options` field must be configured.

* `bandwidth` - (Required, Int) Specifies the bandwidth capacity specified for the bandwidth package.

* `project_id` - (Optional, String, NonUpdatable) Specifies the project ID.
  If omitted, the provider-level project ID will be used.

* `interflow_mode` - (Optional, String, NonUpdatable) Specifies the bandwidth package applicability.
  Valid values are **Area** and **Region**, defaults to **Area**.

* `spec_code` - (Optional, String, NonUpdatable) Specifies the specification code of the bandwidth package.
  If the value of `interflow_mode` is **Area**, the values are as follows:
  + **bandwidth.aftoela**: Southern Africa-Eastern Latin America on both the Chinese Mainland website and International
    website.
  + **bandwidth.aftonla**: Southern Africa-Northern Latin America on both the Chinese Mainland website and International
    website.
  + **bandwidth.aftowla**: Southern Africa-Western Latin America on both the Chinese Mainland website and International
    website.
  + **bandwidth.aptoaf**: Asia Pacific-Southern Africa on the International website.
  + **bandwidth.aptoap**: Asia Pacific on the International website.
  + **bandwidth.aptoela**: Asia Pacific-Eastern Latin America on both the Chinese Mainland website and International
    website.
  + **bandwidth.aptonla**: Asia Pacific-Northern Latin America on both the Chinese Mainland website and International
    website.
  + **bandwidth.aptowla**: Asia Pacific-Western Latin America on both the Chinese Mainland website and International
    website.
  + **bandwidth.cmtoaf**: Chinese mainland-Southern Africa on the International website.
  + **bandwidth.cmtoap**: Chinese mainland-Asia Pacific on the International website.
  + **bandwidth.cmtocm**: Chinese mainland on the International website.
  + **bandwidth.cmtoela**: Chinese mainland-Eastern Latin America on both the Chinese Mainland website and International
    website.
  + **bandwidth.cmtonla**: Chinese mainland-Northern Latin America on both the Chinese Mainland website and
    International website.
  + **bandwidth.cmtowla**: Chinese mainland-Western Latin America on both the Chinese Mainland website and International
    website.
  + **bandwidth.elatoela**: Eastern Latin America on both the Chinese Mainland website and International website.
  + **bandwidth.elatonla**: Eastern Latin America–Northern Latin America on both the Chinese Mainland website and
    International website.
  + **bandwidth.wlatoela**: Western Latin America-Eastern Latin America on both the Chinese Mainland website and
    International website.
  + **bandwidth.wlatonla**: Western Latin America–Northern Latin America on both the Chinese Mainland website and
    International website.
  + **bandwidth.wlatowla**: Western Latin America on both the Chinese Mainland website and International website.

  If the value of `interflow_mode` is **Region**, the value depends on the specified interflow regions,
  e.g. **Beijing4toGuangzhou**.

* `description` - (Optional, String) Specifies the description about the bandwidth package.
  Angle brackets (`<>`) are not allowed.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project that the bandwidth package
  belongs to.

* `resource_id` - (Optional, String) Specifies the ID of the resource that the bandwidth package is bound to.

* `resource_type` - (Optional, String) Specifies the type of the resource that the bandwidth package is bound to.
   Valid value is **cloud_connection**.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the bandwidth package.

* `prepaid_options` - (Optional, List) Specifies the payment attributes.
  The [prepaid_options](#prepaid_options_struct) structure is documented below.

  -> Below is a detailed usage guide for the `prepaid_options` field:
  <br/>1. The field `prepaid_options` must be used in conjunction with the `billing_mode` field. Editing this field
  separately makes no sense.
  <br/>2. The field `prepaid_options` is required when the value of `billing_mode` is **1** or **2**.
  <br/>3. For prepaid type instances, this field does not support repeated editing. For example, when the value of
  `billing_mode` is already **1** or **2**, editing the value of `prepaid_options` has no effect or meaning.

<a name="prepaid_options_struct"></a>
The `prepaid_options` block supports:

* `period_type` - (Required, String) Specifies the unit of a subscription period. Valid values are:
  + **month**: The unit of the subscription period is month.
  + **year**: The unit of the subscription period is year.

* `period_num` - (Required, Int) Specifies the number of subscription periods.
  The value ranges from `1` to `9`, if `period_type` is set to **month**.
  The value ranges from `1` to `3`, if `period_type` is set to **year**.

* `is_auto_renew` - (Optional, Bool) Specifies whether to enable auto renewal.
  + **true**: Auto renewal is enabled.
  + **false**: Auto renewal is disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The bandwidth package status.
  The valid value are as follows:
  + **ACTIVE**: Bandwidth packages are available.

* `created_at` - The time when the resource was created.

* `updated_at` - The time when the resource was updated.

* `order_id` - The order ID of the bandwidth package.

* `product_id` - The product ID of the bandwidth package.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The bandwidth package can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cc_bandwidth_package.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `prepaid_options`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cc_bandwidth_package" "test" {
  ...

  lifecycle {
    ignore_changes = [
      prepaid_options,
    ]
  }
}
```
