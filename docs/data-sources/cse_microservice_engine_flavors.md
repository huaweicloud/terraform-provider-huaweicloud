---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice_engine_flavors"
description: |-
  Use this data source to query available microservice engine flavors within HuaweiCloud.
---

# huaweicloud_cse_microservice_engine_flavors

Use this data source to query available microservice engine flavors within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cse_microservice_engine_flavors" "test" {
  version = "CSE2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the microservice engine flavors are located.  
  If omitted, the provider-level region will be used.

* `version` - (Optional, String) Specifies the version that used to filter the microservice engine flavors.
  + **CSE2**
  + **Nacos2**
  + **MicroGateway**

  Defaults to **CSE2**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - All microservice engine flavors that match the filter parameters.  
  The [flavors](#cse_microservice_engine_flavors) structure is documented below.

<a name="cse_microservice_engine_flavors"></a>
The `flavors` block supports:

* `id` - The ID of the microservice engine flavor.  
  
  -> For resource `huaweicloud_cse_microservice_engine`, the flavor of **Nacos2** need to declare the number of
     capacity units (1 unit = 50 microservice instances) when reference. The format is as follows:  
     **{flavor}.{capacity_units}**  
     If you need to create an engine with 500 microservice instances, the flavor corresponding to
     **cse.nacos2.c1.large** is **cse.nacos2.c1.large.10**.  
     The specifications of **CSE2** can be used directly without special processing.

* `spec` - The specification detail of the microservice engine flavor.  
  The [spec](#cse_microservice_engine_flavor_spec) structure is documented below.

<a name="cse_microservice_engine_flavor_spec"></a>
The `spec` block supports:

* `available_cpu_memory` - The CPU and memory combinations (each value separated by a hyphen) that the flavor is
  allowed, in string format and separated by the commas (,).
  For example, **2-4,4-8** means this flavors supports **2u4g** and **4u8g** instances create.

* `linear` - Whether the microservice engine flavor is a linear flavor, in string format.
  + **true**
  + **false**

* `available_prefix` - The flavor name prefix of the available node, e.g. **s,c,t**.
