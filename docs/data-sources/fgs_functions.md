---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_functions"
description: ""
---

# huaweicloud_fgs_functions

Use this data source to filter FGS functions within HuaweiCloud.

## Example Usage

### Obtain all public functions

```hcl
data "huaweicloud_fgs_functions" "test" {}
```

### Obtain specific public function by package name

```hcl
data "huaweicloud_fgs_functions" "test" {
  package_name = "default"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the functions.
  If omitted, the provider-level region will be used.

* `package_name` - (Optional, String) Specifies the package name of the function to which the functions belong.

* `urn` - (Optional, String) Specifies the function urn used to query specified function.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the function.

* `runtime` - (Optional, String) Specifies the dependent package runtime to match. Valid values: **Java8**,
  **Node.js6.10**, **Node.js8.10**, **Node.js10.16**, **Node.js12.13**, **Python2.7**, **Python3.6**, **Go1.8**,
  **Go1.x**, **C#(.NET Core 2.0)**, **C#(.NET Core 2.1)**, **C#(.NET Core 3.1)** and **PHP7.3**.

* `name` - (Optional, String) Specifies the name of the function.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A data source ID.

* `functions` - The filtered function. The object structure is documented below.

The functions block supports:

* `name` - The function name.

* `urn` - The function urn used to query specified function.

* `package` - The package name of the function to which the functions belong.

* `runtime` - The environment for executing the function.

* `timeout` - The timeout interval of the function, ranges from 3s to 900s.

* `code_type` - The function code type, which can be:
  + **inline**: inline code.
  + **zip**: ZIP file.
  + **jar**: JAR file or java functions.
  + **obs**: function code stored in an OBS bucket.

* `handler` - The entry point of the function.

* `memory_size` - The memory size(MB) allocated to the function.

* `code_url` - The code url.

* `code_filename` - The name of a function file.

* `user_data` - The Key/Value information defined for the function.

* `encrypted_user_data` - The key/value information defined to be encrypted for the function.

* `version` - The function version.

* `agency` - The agency.

* `app_agency` - An execution agency enables you to obtain a token or an AK/SK for
  accessing other cloud services.

* `description` - The description of the version alias.

* `vpc_id` - The VPC ID.

* `network_id` - The network ID of subnet.

* `max_instance_num` - The maximum number of instances of the function.  

* `initializer_handler` - The initializer of the function.

* `initializer_timeout` - The maximum duration the function can be initialized.

* `enterprise_project_id` - The enterprise project id of the function.

* `log_group_id` - The LTS log group ID.

* `log_stream_id` - The LTS log stream ID.

* `functiongraph_version` - The FunctionGraph version.
