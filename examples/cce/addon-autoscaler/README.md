# Create a CCE Addon for Cluster Autoscaler

This example provides best practice code for using Terraform to install and configure the Cluster Autoscaler add-on for
HuaweiCloud CCE (Cloud Container Engine).

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing CCE cluster (Make sure the cluster have at least one non-default node group and with at least two nodes)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CCE cluster is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `addon_version` - The version of the autoscaler addon template

#### Optional Variables

* `cluster_id` - The ID of the CCE cluster (default: "")
  If not specified, cluster_name must be provided
* `cluster_name` - The name of the CCE cluster (default: "")
  If not specified, cluster_id must be provided
* `addon_template_name` - The name of the CCE addon template (default: "autoscaler")
* `project_id` - The ID of the project (default: "")
  If not specified, will be automatically detected from the region

## Usage

1. Copy this example script to your `main.tf`.

2. Create a `terraform.tfvars` file and fill in the required variables:

   ```hcl
   # CCE Configuration
   cluster_id    = "your_cce_cluster_id"
   addon_version = "1.32.5" # The version of cluster is v1.32
   ```

3. Initialize Terraform:

   ```bash
   $ terraform init
   ```

4. Review the Terraform plan:

   ```bash
   $ terraform plan
   ```

5. Apply the configuration:

   ```bash
   $ terraform apply
   ```

6. To clean up the resources:

   ```bash
   $ terraform destroy
   ```

## Configuration Details

### Cluster Autoscaler Parameters

The example automatically configures the Cluster Autoscaler with the following key parameters:

* **Scaling Configuration**:

  - `scaleUpCpuUtilizationThreshold`: CPU utilization threshold for scale-up
  - `scaleUpMemUtilizationThreshold`: Memory utilization threshold for scale-up
  - `scaleDownUtilizationThreshold`: Utilization threshold for scale-down
  - `scaleDownEnabled`: Enable/disable scale-down functionality

* **Timing Configuration**:

  - `scaleDownDelayAfterAdd`: Delay before scale-down after adding nodes
  - `scaleDownDelayAfterDelete`: Delay before scale-down after deleting nodes
  - `scaleDownUnneededTime`: Time before considering a node for removal

* **Resource Limits**:

  - `maxNodesTotal`: Maximum total number of nodes
  - `coresTotal`: Maximum total CPU cores
  - `memoryTotal`: Maximum total memory

* **Multi-AZ Support**:

  - `multiAZEnabled`: Enable multi-AZ balancing
  - `multiAZBalance`: Balance nodes across availability zones

### Custom Configuration

The example preserves all original template parameters while adding required fields:

* `cluster_id`: Automatically set to the target cluster
* `tenant_id`: Automatically set to the project ID

## Best Practices

1. **Version Management**: Always specify a specific addon version to ensure consistent deployments across environments.

2. **Cluster Identification**: Use either `cluster_id` or `cluster_name` to identify the target cluster.

3. **Resource Limits**: Configure appropriate resource limits (`maxNodesTotal`, `coresTotal`, `memoryTotal`) based on
   your cluster requirements and budget.

4. **Scaling Thresholds**: Adjust CPU and memory utilization thresholds based on your workload characteristics:
   - Higher thresholds for CPU-intensive workloads
   - Lower thresholds for memory-intensive workloads

5. **Scale-down Configuration**: Configure appropriate delays to prevent rapid scaling cycles:
   - `scaleDownDelayAfterAdd`: Prevents immediate scale-down after scale-up
   - `scaleDownUnneededTime`: Ensures nodes are truly unused before removal

6. **Multi-AZ Deployment**: Enable multi-AZ balancing for high availability and better resource distribution.

## Troubleshooting

* **Cluster Not Found**: Ensure the cluster exists and you have appropriate permissions
* **Addon Installation Failed**: Check the addon version compatibility with your cluster version
* **Scaling Issues**: Verify the autoscaler configuration and node group settings

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The autoscaler addon will automatically manage node scaling based on pod scheduling requirements
* Monitor cluster resource usage after installation to optimize autoscaler parameters

## Version Compatibility

Please refer to the [official Huawei Cloud documentation](https://support.huaweicloud.com/intl/en-us/usermanual-cce/cce_10_0154.html#section6)
for detailed version compatibility information between CCE cluster versions and Autoscaler addon versions.

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.37.0 |
