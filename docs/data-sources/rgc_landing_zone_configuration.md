---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_landing_zone_configuration"
description: |-
  Use this data source to get landing zone configuration in Resource Governance Center.
---

# huaweicloud_rgc_landing_zone_configuration

Use this data source to get landing zone configuration in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_landing_zone_configuration" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `common_configuration` - Information about the common configuration of the landing zone.

  The [common_configuration](#common_configuration) structure is documented below.

* `logging_configuration` - Information about the logging configuration of the landing zone.

  The [logging_configuration](#logging_configuration) structure is documented below.

* `organization_structure` - Information about the organization structure of the landing zone.

  The [organization_structure](#organization_structure) structure is documented below.

* `regions` - A list of regions and their configuration status.

  The [regions](#regions) structure is documented below.

<a name="common_configuration"></a>
The `common_configuration` block supports:

* `home_region` - The home region of the landing zone.

* `cloud_trail_type` - Indicates whether the cloud trail type is enabled.

* `identity_center_status` - The status of the identity center.

* `organization_structure_type` - The type of organization structure.

<a name="logging_configuration"></a>
The `logging_configuration` block supports:

* `logging_bucket_name` - The name of the logging bucket.

* `access_logging_bucket` - Configuration details for the access logging bucket.

* `logging_bucket` - Configuration details for the logging bucket.

<a name="access_logging_bucket"></a>
The `access_logging_bucket` block supports:

* `retention_days` - The number of days to retain logs for access logging bucket.

* `enable_multi_az` - Whether multi-AZ is enabled for the access logging bucket.

<a name="logging_bucket"></a>
The `logging_bucket` block supports:

* `retention_days` - The number of days to retain logs for logging bucket.

* `enable_multi_az` - Whether multi-AZ is enabled for the logging bucket.

<a name="organization_structure"></a>
The `organization_structure` block supports:

* `organizational_unit_name` - The name of the organizational unit.

* `organizational_unit_type` - The type of organizational unit.

* `accounts` - A list of accounts within the organizational unit.
  
The accounts is documented below.

<a name="accounts"></a>
The `accounts` block supports:

* `account_id` - The ID of the account.

* `account_type` - The type of the account.

<a name="regions"></a>
The `regions` block supports:

* `region` - The name of the regulated region.

* `region_configuration_status` - The configuration status of the regulated region.
