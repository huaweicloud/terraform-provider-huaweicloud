---
subcategory: "Enterprise Project Management Service (EPS)"
---

# huaweicloud\_enterprise\_project

Use this data source to get an enterprise project from HuaweiCloud

## Example Usage

```hcl
data "huaweicloud_enterprise_project" "test" {
  name = "test"
}
```

## Resources Supported Currently:
Service Name | Resource Name
---- | ---
VPC | huaweicloud_vpc<br>huaweicloud_vpc_eip<br>huaweicloud_vpc_bandwidth<br>huaweicloud_networking_secgroup
ECS | huaweicloud_compute_instance
CCE | huaweicloud_cce_cluster
RDS | huaweicloud_rds_instance
OBS | hauweicloud_obs_bucket
SFS | hauweicloud_sfs_file_system
DCS | huaweicloud_dcs_instance
NAT | huaweicloud_nat_gateway
CDM | huaweicloud_cdm_cluster
CDN | huaweicloud_cdn_domain
GaussDB | huaweicloud_gaussdb_cassandra_instance<br>huaweicloud_gaussdb_mysql_instance<br>huaweicloud_gaussdb_opengauss_instance

## Argument Reference

* `name` - (Optional) Specifies the enterprise project name. Fuzzy search is supported.

* `id` - (Optional) Specifies the ID of an enterprise project. The value 0 indicates enterprise project default.

* `status` - (Optional) Specifies the status of an enterprise project.
    - 1 indicates Enabled.
    - 2 indicates Disabled.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `description` - Provides supplementary information about the enterprise project.

* `created_at` - Specifies the time (UTC) when the enterprise project was created. Example: 2018-05-18T06:49:06Z

* `updated_at` - Specifies the time (UTC) when the enterprise project was modified. Example: 2018-05-28T02:21:36Z

