---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_deployed_resources"
description: |-
  Use this data source to get the CCM SSL deployed resources.
---

# huaweicloud_ccm_deployed_resources

Use this data source to get the CCM SSL deployed resources.

## Example Usage

```hcl
variable "certificate_ids" {
  type = list(string)
}

variable "service_names" {
  type = list(string)
}

data "huaweicloud_ccm_deployed_resources" "test" {
  certificate_ids = var.certificate_ids
  service_names   = var.service_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `service_names` - (Required, List) Specifies the service name list. Valid values are:
  + **WAF**: Query WAF resources associated with a certificate.
  + **CDN**: Query the CDN resources associated with a certificate.
  + **ELB**: Query the resources associated with ELB (classic) of a certificate.
  + **ALL**: Query the resources of the preceding four services.

* `certificate_ids` - (Required, List) Specifies the certificate ID list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The request result list.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `certificate_id` - The certificate ID.

* `total_num` - The number of resources deployed for the certificate in the current service.

* `deployed_resources` - The list of resources deployed under the current certificate.

  The [deployed_resources](#results_deployed_resources_struct) structure is documented below.

<a name="results_deployed_resources_struct"></a>
The `deployed_resources` block supports:

* `region_resources` - The region resource list.

  The [region_resources](#deployed_resources_region_resources_struct) structure is documented below.

* `service` - The name of the resource service where the certificate has been deployed. Valid values are:
  + **WAF**: The certificate is associated with WAF resources.
  + **CDN**: The certificate is associated with the resources of the content delivery network.
  + **ELB**: The certificate is associated with ELB (classic) resources.

* `resource_num` - The number of resources deployed for the certificate in the current service.

* `resource_location` - The global service or region-level service.

<a name="deployed_resources_region_resources_struct"></a>
The `region_resources` block supports:

* `region_id` - The region ID. If the service is a global service, the value of this field is **global**.
  Other services are named based on the IAM.

* `is_error` - Whether an exception occurs in the response when the resource information of the current region
  is requested.
  + **true**: An exception occurs. The statistics of the current region are inaccurate.
  + **false**: No exception occurs. The statistics of the current region are correct.

* `resources` - The resource set. The identifier of each resource is in the format of `resource ID + : + resource name`.

  The [resources](#region_resources_resources_struct) structure is documented below.

<a name="region_resources_resources_struct"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name.
