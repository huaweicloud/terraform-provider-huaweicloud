---
subcategory: "Function Graph Service (FGS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_async_invoke_configurations"
description: |-
  Use this data source to get the list of async invoke configurations within HuaweiCloud.
---

# huaweicloud_fgs_async_invoke_configurations

Use this data source to get the list of async invoke configurations within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_fgs_async_invoke_configurations" "all" {
  function_urn = huaweicloud_fgs_function.test.urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the async invoke configurations are located.  
  If omitted, the provider-level region will be used.

* `function_urn` - (Required, String) Specifies the function URN to query async invoke configurations.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - The list of queried async invoke configurations.  
  The [configurations](#fgs_async_invoke_configurations_attr) structure is documented below.

<a name="fgs_async_invoke_configurations_attr"></a>
The `configurations` block supports:

* `func_urn` - The function URN.

* `max_async_event_age_in_seconds` - The maximum validity period of a message.

* `max_async_retry_attempts` - The maximum number of retry attempts to be made if asynchronous invocation fails.

* `destination_config` - The destination configuration for async invoke.  
  The [destination_config](#fgs_async_invoke_configuration_destination_config_attr) structure is documented below.

* `created_time` - The creation time of the async invoke configuration.

* `last_modified` - The last modification time of the async invoke configuration.

* `enable_async_status_log` - Whether async invoke status persistence is enabled.

<a name="fgs_async_invoke_configuration_destination_config_attr"></a>
The `destination_config` block supports:

* `on_success` - The target to be invoked when a function is successfully executed.  
  The [on_success](#fgs_async_invoke_configuration_destination_config_item_attr) structure is documented below.

* `on_failure` - The target to be invoked when a function fails to be executed.  
  The [on_failure](#fgs_async_invoke_configuration_destination_config_item_attr) structure is documented below.

<a name="fgs_async_invoke_configuration_destination_config_item_attr"></a>
The `on_success` and `on_failure` blocks support:

* `destination` - The target type.
  + **OBS**
  + **SMN**
  + **DIS**
  + **FunctionGraph**

* `param` - The parameters corresponding to the target service, in JSON format.
