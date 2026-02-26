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
  + **3**: pay-per-use on the Chinese mainland website
  + **4**: pay-per-use on the International website
  + **5**: 95th percentile bandwidth billing on the Chinese mainland website
  + **6**: 95th percentile bandwidth billing on the International website

  -> This argument can only be modified to **5** and **6**.

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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The bandwidth package status.
  The valid value are as follows:
  + **ACTIVE**: Bandwidth packages are available.

## Import

The bandwidth package can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cc_bandwidth_package.test <id>
```
