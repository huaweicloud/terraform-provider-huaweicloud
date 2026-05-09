---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_notebooks"
description: |-
  Use this data source to get the list of ModelArts Notebook instances within HuaweiCloud.
---

# huaweicloud_modelarts_notebooks

Use this data source to get the list of ModelArts Notebook instances within HuaweiCloud.

## Example Usage

### Query all notebooks without any filter

```hcl
data "huaweicloud_modelarts_notebooks" "test" {}
```

### Query the notebooks and using notebook name to filter

```hcl
variable "notebook_name" {}

data "huaweicloud_modelarts_notebooks" "test" {
  name = var.notebook_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the notebooks are located.  
  If omitted, the provider-level region will be used.

* `feature` - (Optional, String) Specifies the feature type of the notebooks to be queried.

* `notebook_id` - (Optional, String) Specifies the notebook instance ID to be queried.

* `name` - (Optional, String) Specifies the name of the notebooks to be queried. Fuzzy match is supported.

* `status` - (Optional, String) Specifies the status of the notebooks to be queried.

* `workspace_id` - (Optional, String) Specifies the workspace ID of the notebooks to be queried.

* `flavor_id` - (Optional, String) Specifies the flavor of the notebooks to be queried.

* `image_id` - (Optional, String) Specifies the image ID of the notebooks to be queried.

* `billing` - (Optional, String) Specifies the billing type of the notebooks to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `notebooks` - The list of notebooks that match the filter parameters.  
  The [notebooks](#modelarts_notebooks_attr) structure is documented below.

<a name="modelarts_notebooks_attr"></a>
The `notebooks` block supports:

* `id` - The ID of the notebook instance.

* `name` - The name of the notebook instance.

* `description` - The description of the notebook instance.

* `status` - The status of the notebook instance.

* `feature` - The feature type of the notebook instance.

* `flavor_id` - The flavor of the notebook instance.

* `workspace_id` - The workspace ID to which the notebook belongs.

* `pool_id` - The dedicated resource pool ID of the notebook instance.

* `pool_name` - The dedicated resource pool name of the notebook instance.

* `image_id` - The image ID of the notebook instance.

* `image_name` - The image name of the notebook instance.

* `image_type` - The image type of the notebook instance.

* `image_swr_path` - The SWR path of the image of the notebook instance.

* `url` - The access URL of the notebook instance.

* `fail_reason` - The failure reason of the notebook instance.

* `auto_stop_enabled` - Whether the auto stop feature is enabled for the notebook instance.

* `lease_duration` - The lease duration of the notebook instance, in milliseconds.

* `lease_type` - The auto stop type of the notebook instance.

* `key_pair` - The SSH key pair name when SSH access is configured.

* `ssh_uri` - The SSH access URI when SSH access is configured.

* `allowed_access_ips` - The allowed access IP list for SSH when SSH access is configured.

* `user_id` - The user ID of the notebook instance.

* `ip` - The IP address of the node where the notebook instance is located.

* `jupyter_version` - The JupyterLab version of the notebook instance.

* `billing_items` - The billing resource types of the notebook instance.

* `custom_spec` - The custom specification of the notebook instance, in JSON format.

* `user_vpc` - The user VPC configuration of the notebook instance, in JSON format.

* `volume` - The system volume configuration of the notebook instance.  
  The [volume](#modelarts_notebooks_volume_attr) structure is documented below.

* `data_volumes` - The extended storage list of the notebook instance.  
  The [data_volumes](#modelarts_notebooks_data_volumes_attr) structure is documented below.

* `action_progress` - The initialization progress of the notebook instance.  
  The [action_progress](#modelarts_notebooks_action_progress_attr) structure is documented below.

* `tags` - The tags of the notebook instance.  
  The [tags](#modelarts_notebooks_tags_attr) structure is documented below.

* `created_at` - The creation time of the notebook instance.

* `updated_at` - The last update time of the notebook instance.

<a name="modelarts_notebooks_volume_attr"></a>
The `volume` block supports:

* `type` - The storage category of the notebook instance.

* `ownership` - The ownership type of the storage.

* `size` - The storage capacity.

* `uri` - The storage URI.

* `id` - The storage ID.

* `mount_path` - The mount path of the storage.

<a name="modelarts_notebooks_data_volumes_attr"></a>
The `data_volumes` block supports:

* `type` - The extended storage category.

* `mount_path` - The mount path of the extended storage.

* `path` - The source path of the extended storage.

* `status` - The status of the extended storage.

* `mount_type` - The mount type of the extended storage.

<a name="modelarts_notebooks_action_progress_attr"></a>
The `action_progress` block supports:

* `status` - The status of the step.

* `step` - The step number.

* `description` - The description of the step.

<a name="modelarts_notebooks_tags_attr"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
