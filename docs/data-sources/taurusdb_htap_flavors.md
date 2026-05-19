---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_flavors"
description: |-
  Use this data source to query the list of TaurusDB HTAP instance flavors within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_flavors

Use this data source to query the list of TaurusDB HTAP instance flavors within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_taurusdb_htap_flavors" "test" {
  engine_name            = "star-rocks"
  availability_zone_mode = "single"
}
```

### Filter by spec code and version name

```hcl
variable "flavor_spec_code" {}

data "huaweicloud_taurusdb_htap_flavors" "test" {
  engine_name            = "star-rocks"
  availability_zone_mode = "single"
  spec_code              = var.flavor_spec_code
  version_name           = "3.1.6.0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP instance flavors.
  If omitted, the provider-level region will be used.

* `engine_name` - (Required, String) Specifies the HTAP engine type. Value options: **star-rocks**.

* `availability_zone_mode` - (Required, String) Specifies the AZ type of the HTAP instance flavor.
  Value options: **single**.

* `spec_code` - (Optional, String) Specifies the specification code of the flavor.

* `version_name` - (Optional, String) Specifies the HTAP database version of the flavor.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of HTAP instance flavors that matched filter parameters.
  The [flavors](#taurusdb_htap_flavors_attr) structure is documented below.

<a name="taurusdb_htap_flavors_attr"></a>
The `flavors` block supports:

* `id` - The ID of the HTAP instance flavor.

* `type` - The type of CPU architecture.
  The valid values are as follows:
  + **x86**: exclusive x86
  + **arm**: exclusive arm
  + **generalX86**: general-purpose x86

* `spec_code` - The specification code of the HTAP instance flavor.

* `version_name` - The version name of the HTAP database.

* `instance_mode` - The instance mode of the HTAP instance.

* `vcpus` - The number of vCPUs for the HTAP instance flavor.

* `ram` - The memory size in GB for the HTAP instance flavor.

* `az_status` - Map of availability zone names and their status for the HTAP instance flavor.
  The valid values are as follows:
  + **normal**: The specifications are available in the AZ.
  + **unsupported**: The specifications are not supported.
  + **sellout**: The specifications are sold out.
