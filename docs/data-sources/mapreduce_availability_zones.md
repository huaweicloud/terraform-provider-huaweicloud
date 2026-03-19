---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_availability_zones"
description: |-
  Use this data source to query the availability zones of MRS within HuaweiCloud.
---

# huaweicloud_mapreduce_availability_zones

Use this data source to query the availability zones of MRS within HuaweiCloud.

## Example Usage

### Query all availability zones

```hcl
data "huaweicloud_mapreduce_availability_zones" "test" {}
```

### Query availability zones by scope

```hcl
data "huaweicloud_mapreduce_availability_zones" "test" {
  scope = "Center"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the availability zones are located.
  If omitted, the provider-level region will be used.

* `scope` - (Optional, String) Specifies the availability zone scope.  
  The valid values are as follows:
  + **Center**
  + **Edge**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `default_az_code` - The default availability zone code.

* `support_physical_az_group` - Whether the physical availability zone grouping is supported.

* `available_zones` - The availability zone list that matched the filter parameters.  
  The [available_zones](#mapreduce_availability_zones_available_zones) structure is documented below.

<a name="mapreduce_availability_zones_available_zones"></a>
The `available_zones` block supports:

* `id` - The availability zone code.

* `az_id` - The availability zone ID.

* `az_code` - The availability zone code.

* `az_name` - The availability zone name.

* `status` - The availability zone status.

* `region_id` - The region ID.

* `az_group_id` - The availability zone group ID.

* `az_type` - The availability zone type.  
  The valid values are as follows:
  + **Core**
  + **Satellite**
  + **Dedicated**
  + **Virtual**
  + **Edge**
  + **EdgeCentral**

* `az_category` - The availability zone category.  
  The valid values are as follows:
  + **0**: Large cloud primary AZ.
  + **21**: Local AZ.
  + **41**: Edge AZ.

* `charge_policy` - The charge policy of the availability zone.  
  The valid values are as follows:
  + **charge**
  + **notCharge**

* `az_tags` - The availability zone tags.  
  The [az_tags](#mapreduce_availability_zones_az_tags) structure is documented below.

<a name="mapreduce_availability_zones_az_tags"></a>
The `az_tags` block supports:

* `mode` - The availability zone mode.  
  The valid values are as follows:
  + **dedicated**
  + **shared**

* `alias` - The alias of the availability zone.

* `public_border_group` - The public border group to which the availability zone belongs.
