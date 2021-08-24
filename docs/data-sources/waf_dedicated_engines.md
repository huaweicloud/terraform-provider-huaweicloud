---
subcategory: "Web Application Firewall (WAF)"
---

# huaweicloud_waf_dedicated_engines

Use this data source to get a list of WAF dedicated engine.

## Example Usage

```hcl
variable engine_name {}

data "huaweicloud_waf_dedicated_engines" "engine" {
  name = var.engine_name
}
```

### 

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the WAF dedicated engine.
  If omitted, the provider-level region will be used. 

* `name` - (Optional, String) The name of WAF dedicated engine. 


## Attributes Reference

The following attributes are exported:

* `engines` -  A list of WAF dedicated engines.

The `engines` block supports:

* `id` -  The id of the engine.

* `name` - The name of WAF dedicated engine. 

* `available_zone` - The available zone names for the WAF dedicated engines.

* `specification_code` - The specification code of WAF dedicated engine. Different specifications have different throughput. 

* `cpu_architecture` - The ECS cpu architecture of WAF dedicated engine.

* `cpu_flavor` - The ECS specification, e.g. `c6.large.2`.

* `vpc_id` - The VPC id of WAF dedicated engine.

* `subnet_id` - The subnet id of WAF dedicated engine VPC.

* `security_group` - The security group of the engine. This is an array of security group ids.

* `server_id` - The service of the engine.

* `service_ip` - The service ip of the engine.

* `run_status` - The running status of the engine.. Values are:
  * `0` - Engine is creating.
  * `1` - Engine has created.
  * `2` - Engine is deleting.
  * `3` - Engine has deleted.
  * `4` - Engine create failed.

* `access_status` - The access status of the engine. `0`: inaccessible, `1`: accessible.

* `upgradable` - The engine is to support upgrades. `0`: Cannot be upgraded, `1`: Can be upgraded.

