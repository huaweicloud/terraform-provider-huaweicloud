---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_metering"
description: |-
  Use this data source to get the metering information for a specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_job_metering

Use this data source to get the metering information for a specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" { 
  type = string 
}

data "huaweicloud_drs_job_metering" "test" { 
  job_id = var.job_id 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `product_info_list` - The metering information list.

  The [product_info_list](#product_info_list_struct) structure is documented below.

<a name="product_info_list_struct"></a>
The `product_info_list` block supports:

* `id` - The ID identifier.

* `cloud_service_type` - The cloud service type of the user-purchased cloud service product.

* `resource_type` - The resource type of the user-purchased cloud service product.
  The valid values are as follows:
  + **hws.resource.type.drs.instance**: Instance.
  + **hws.resource.type.drs.vm**: Virtual machine.
  + **hws.resource.type.drs.volume**: Volume.
  + **dbs.instanceName**: Instance name.
  + **hws.resource.type.drs.flow**: DRS flow fee.
  + **dbs.tag**: User tag.
  + **dbs.enterpriseProjectId**: Enterprise project.

* `resource_spec_code` - The resource specification of the user-purchased cloud service product.

* `resource_size` - The resource capacity measurement identifier.

* `resource_size_measure_id` - The resource capacity size, for example, the purchased volume size or bandwidth size.

* `usage_factor` - The usage factor.
  The valid values are as follows:
  + **Duration**: Cloud server.
  + **flow**: Traffic.

* `usage_value` - The usage value.

* `usage_measure_id` - The usage unit identifier.
  The valid values are as follows:
  + `4`: Hour.
  + `10`: GB.
  + `11`: MB.
  + `13`: Byte.
  + `17`: FLOW_GB.
