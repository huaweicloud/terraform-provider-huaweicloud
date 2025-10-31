---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_landing_zone"
description: |-
  Manages an RGC landing zone resource within HuaweiCloud.
---

# huaweicloud_rgc_landing_zone

Manages an RGC landing zone resource within HuaweiCloud.

## Example Usage

### SetupLandingZone with existed account

```hcl
variable "home_region" {}
variable "organizational_unit_name" {}
variable "logging_account_name" {}
variable "logging_account_id" {}
variable "audit_account_name" {}
variable "audit_account_id" {}
variable "audit_account_email" {}
variable "manage_account_email" {}

resource "huaweicloud_rgc_landing_zone" "test" {
  home_region            = var.home_region
  identity_center_status = "ENABLE"
  identity_store_email   = var.manage_account_email
  cloud_trail_type       = true
  region_configuration_list {
    region                      = var.home_region
    region_configuration_status = "ENABLED"
  }
  organization_structure {
    organizational_unit_type = "CORE"
    organizational_unit_name = var.organizational_unit_name
    accounts {
      account_name = var.logging_account_name
      account_type = "LOGGING"
      account_id   = var.logging_account_id
    }
    accounts {
      account_name  = var.audit_account_name
      account_type  = "SECURITY"
      account_id    = var.audit_account_id
      account_email = var.audit_account_email
    }
  }
  logging_configuration {
    logging_bucket {
      retention_days = 365
    }
    access_logging_bucket {
      retention_days = 3650
    }
  }
}
```

### SetupLandingZone to create new account

```hcl
variable "home_region" {}
variable "organizational_unit_name" {}
variable "logging_account_name" {}
variable "audit_account_name" {}
variable "audit_account_email" {}
variable "manage_account_email" {}

resource "huaweicloud_rgc_landing_zone" "test" {
  home_region            = var.home_region
  identity_center_status = "ENABLE"
  identity_store_email   = var.manage_account_email
  cloud_trail_type       = true
  region_configuration_list {
    region                      = var.home_region
    region_configuration_status = "ENABLED"
  }
  organization_structure {
    organizational_unit_type = "CORE"
    organizational_unit_name = var.organizational_unit_name
    accounts {
      account_name = var.logging_account_name
      account_type = "LOGGING"
    }
    accounts {
      account_name  = var.audit_account_name
      account_type  = "SECURITY"
      account_email = var.audit_account_email
    }
  }
  logging_configuration {
    logging_bucket {
      retention_days = 365
    }
    access_logging_bucket {
      retention_days = 3650
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `home_region` - (Required, String, NonUpdatable) Specifies the ID of home region.

* `organization_structure_type` - (Required, String, NonUpdatable) Specifies the type of the organization structure.
  If creating a core organizational unit, this field should be set to STANDARD;
  however, if not creating it, the field should be set to NON_STANDARD, with the default value being STANDARD.

* `identity_store_email` - (Optional, String) Specifies the email of idc account, if identity_center_status set true,
  this field will also be set.

* `identity_center_status` - (Optional, String) Specifies status of identity center.

* `deny_ungoverned_regions` - (Optional, Bool, NonUpdatable) Specifies the access of multi-region,
  default value is false

* `cloud_trail_type` - (Optional, Bool) Specifies the access of organizational aggregation.

* `kms_key_id` - (Required, String, NonUpdatable) Specifies the ID of kms key.

* `baseline_version` - (Optional, String) Specifies the excepted landing zone version for updating landing zone.

* `region_configuration_list` - (Required, List, NonUpdatable) Specifies the status of enrolled region.

  The [region_configuration_list](#region_configuration_list) structure is documented below.

<a name="region_configuration_list"></a>
The `region_configuration_list` block supports:

* `region` - (Required, String) Specifies the ID of enrolled region.

* `region_configuration_status` - (Required, String) Specifies the status of enrolled region.

* `organization_structure` - (Required, List, NonUpdatable) Specifies the organizational structure.

  The [organization_structure](#organization_structure) structure is documented below.

<a name="organization_structure"></a>
The `organization_structure` block supports:

* `organizational_unit_type` - (Required, String) Specifies the type of registered organizational unit,
  only support CORE or CUSTOM.

* `organizational_unit_name` - (Optional, String) Specifies the name of created organizational unit,
  leaving this field means that the organizational unit will not be created.

* `accounts` - (Optional, List) Specifies the account structure.

  The [accounts](#accounts) structure is documented below.

<a name="accounts"></a>
The `accounts` block supports:

* `account_name` - (Optional, String) Specifies the name of created account or enrolled account.

* `account_type` - (Required, String) Specifies the name of created account or enrolled account,
  only support LOGGING or SECURITY.

* `account_id` - (Optional, String) Specifies the id of enrolled account.

* `phone` - (Optional, String) Specifies the name of created account, domestic account is required.

* `account_email` - (Optional, String) Specifies the name of account, international account or security
  account is required.

* `logging_configuration` - (Required, List) Specifies the logging bucket configuration.

  The [logging_configuration](#logging_configuration) structure is documented below.

<a name="logging_configuration"></a>
The `logging_configuration` block supports:

* `logging_bucket_name` -  (Optional, String) Specifies the name of logging bucket, set this field means
  using existed bucket, RGC will not create bucket in setup landing zone progress.

* `logging_bucket` - (Optional, List) Specifies the logging bucket structure.

  The [logging_bucket](#logging_bucket) structure is documented below.

* `access_logging_bucket` - (Optional, List) Specifies the access logging bucket structure.

  The [access_logging_bucket](#access_logging_bucket) structure is documented below.

<a name="logging_bucket"></a>
The `logging_bucket` and  access_logging_bucket block supports:

* `logging_bucket_retention_days` - (Required, Int) Specifies retention days of logging bucket.

* `logging_bucket_enable_multi_az` - (Optional, Bool) Specifies multi-az storage access of logging bucket,
  default value is false.

<a name="access_logging_bucket"></a>
The `access_logging_bucket` block supports:

* `logging_bucket_retention_days` - (Required, Int) Specifies retention days of access logging bucket.

* `logging_bucket_enable_multi_az` - (Optional, Bool) Specifies multi-az storage access of access logging bucket,
  default value is false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `landing_zone_status` - Indicates the configuration status of landing zone.

* `deployed_version` Indicates the configuration deployed version of landing zone.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.

* `update` - Default is 60 minutes.

* `delete` - Default is 60 minutes.

## Import

The RGC landing zone can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rgc_landing_zone.test <id>
```
