---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_scripts"
description: |-
  Use this data source to query DataArts Factory scripts within HuaweiCloud.
---

# huaweicloud_dataarts_factory_scripts

Use this data source to query DataArts Factory scripts within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_factory_scripts" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the scripts are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the scripts belong.

* `script_name` - (Optional, String) Specifies the name of script to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scripts` - The list of scripts that matched filter parameters.  
  The [scripts](#dataarts_factory_scripts) structure is documented below.

<a name="dataarts_factory_scripts"></a>
The `scripts` block supports:

* `id` - The script ID.

* `name` - The script name.

* `type` - The script type.

* `directory` - The directory path where the script is located.

* `create_user` - The user who created the script.

* `connection_name` - The connection name associated with the script.

* `database` - The database associated with the script.

* `queue_name` - The DLI queue name associated with the script.

* `configuration` - The user-defined configuration parameters associated with the script.

* `description` - The description of the script.

* `modify_time` - The last modification time of the script, in RFC3339 format.

* `owner` - The owner of the script.

* `version` - The version number of the script.
