---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_instance"
description: ""
---

# huaweicloud_dsc_instance

Manages a DSC instance resource within HuaweiCloud.  

## Example Usage

```hcl
resource "huaweicloud_dsc_instance" "test" {
  edition                    = "base_standard"
  charging_mode              = "prePaid"
  period_unit                = "month"
  period                     = 1
  auto_renew                 = "false"
  obs_expansion_package      = 1
  database_expansion_package = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `charging_mode` - (Required, String, ForceNew) Billing mode.  
  The options are as follows:
    + **prePaid**: the yearly/monthly billing mode.

  Changing this parameter will create a new resource.

* `period_unit` - (Required, String, ForceNew) The charging period unit.  
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `period` - (Required, Int, ForceNew) The charging period.  
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `edition` - (Required, String, ForceNew) The edition of DSC.  
  By default, it supports 2 databases and 100GB of OBS storage
  The options are as follows:
    + **base_standard**: Standard Edition.
      It supports **Overview**, **Sensitive Data Identification** and **Data Usage Audit**.
    + **base_professional**: Professional Edition.
      It supports **Overview**, **Sensitive Data Identification**, **Data Usage Audit**, **Data Masking**,
      and **Watermark injection/extraction**

  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String, ForceNew) Whether auto renew is enabled. Valid values are **true** and **false**.  
  Defaults to **false**.  

  Changing this parameter will create a new resource.

* `obs_expansion_package` - (Optional, Int, ForceNew) The size of OBS expansion packages.  
  One expansion package offers 1 TB of OBS storage.

  Changing this parameter will create a new resource.

* `database_expansion_package` - (Optional, Int, ForceNew) The size of database expansion packages.  
  One expansion package offers one database.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The dsc instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dsc_instance.test 0ce123456a00f2591fabc00385ff1234
```
