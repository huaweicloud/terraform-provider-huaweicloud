# Create a PV and PVC with an existing OBS bucket under the CCE environment

This example provides best practice code for using Terraform to create a Persistent Volume (PV) and
Persistent Volume Claim (PVC) using an existing OBS bucket and mount the PVC to a container under the
CCE (Cloud Container Engine) environment (a cluster with one node and public EIP access enabled) in HuaweiCloud.

For the Kubernetes provider authentication, this example demonstrates how to configure the Kubernetes provider to
authenticate with the CCE cluster using cluster certificates.
The authentication is automatically configured using the following CCE cluster certificate information:

* **Cluster CA Certificate**: Used for verifying the cluster's identity
* **Client Certificate**: Used for authenticating the Terraform client
* **Client Key**: Private key for client certificate authentication

The Kubernetes provider configuration automatically extracts and decodes these certificates from the CCE cluster
resource, ensuring secure and seamless authentication without manual certificate management.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `cluster_name` - The name of the CCE cluster
* `keypair_name` - The name of the key pair that is used to access the CCE nodes
* `node_name` - The name of the CCE node
* `bucket_name` - The name of the OBS bucket
* `secret_name` - The name of the Kubernetes secret for OBS authentication
* `secret_data` - The data of the Kubernetes secret (sensitive)
* `pv_name` - The name of the persistent volume
* `pvc_name` - The name of the persistent volume claim
* `deployment_name` - The name of the deployment
* `deployment_containers` - The container list for the deployment
  + `name` - The name of the container
  + `image` - The image of the container
  + `volume_mounts` - The volume mounts of the container
    - `mount_path` - The mount path of the volume

#### Optional Variables

* `availability_zone` - The availability zone where the resources will be created (default: "")
* `vpc_id` - The ID of the VPC (default: "")
* `subnet_id` - The ID of the subnet (default: "")
* `vpc_name` - The name of the VPC (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_name` - The name of the subnet (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `eip_address` - The EIP address of the CCE cluster (default: "")
* `eip_type` - The type of the EIP (default: "5_bgp")
* `bandwidth_name` - The name of the bandwidth (default: "")
* `bandwidth_size` - The size of the bandwidth (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")
* `cluster_flavor_id` - The flavor ID of the CCE cluster (default: "cce.s1.small")
* `cluster_version` - The version of the CCE cluster (default: null)
* `cluster_type` - The type of the CCE cluster (default: "VirtualMachine")
* `container_network_type` - The type of container network (default: "overlay_l2")
* `authentication_mode` - The mode of the CCE cluster (default: "rbac")
* `delete_all_resources_on_terminal` - Whether to delete all resources on terminal (default: true)
* `enterprise_project_id` - The enterprise project ID of the CCE cluster (default: "0")
* `node_flavor_id` - The flavor ID of the node (default: "")
* `node_performance_type` - The performance type of the node (default: "general")
* `node_cpu_core_count` - The CPU core count of the node (default: 4)
* `node_memory_size` - The memory size of the node (default: 8)
* `root_volume_type` - The type of the root volume (default: "SATA")
* `root_volume_size` - The size of the root volume (default: 40)
* `data_volumes_configuration` - The configuration of the data volumes (default: [])
* `bucket_multi_az` - Whether to enable multi-AZ for the OBS bucket (default: true)
* `namespace_name` - The name of the CCE namespace (default: "default")
* `secret_labels` - The labels of the Kubernetes secret (default: {"secret.kubernetes.io/used-by" = "csi"})
* `secret_type` - The type of the Kubernetes secret (default: "cfe/secure-opaque")
* `pv_csi_provisioner_identity` - The identity of the CSI provisioner (default: "everest-csi-provisioner")
* `pv_access_modes` - The access modes of the persistent volume (default: ["ReadWriteMany"])
* `pv_storage` - The storage of the persistent volume (default: "1Gi")  
  This is only for verification purposes, and the set value will not take effect when the PV is created.
* `pv_driver` - The driver of the persistent volume (default: "obs.csi.everest.io")
* `pv_fstype` - The instance type of the persistent volume (default: "s3fs")
* `pv_obs_volume_type` - The type of the OBS volume of the persistent volume (default: "standard")
* `pv_reclaim_policy` - The reclaim policy of the persistent volume (default: "Retain")
* `pv_storage_class_name` - The name of the storage class of the persistent volume (default: "csi-obs")
* `deployment_replicas` - The number of pods for the deployment (default: 2)
* `deployment_volume_name` - The name of the volume of the deployment (default: "pvc-obs-volume")
* `deployment_image_pull_secrets` - The image pull secrets of the deployment (default: [{"name" = "default-secret"}])

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  cluster_name = "tf_test_cluster"
  keypair_name = "tf_test_keypair"
  node_name    = "tf_test_node"
  bucket_name  = "tf_test_bucket"
  secret_name  = "tf_test_secret"
  secret_data = {
    "access.key" = "your_access_key"
    "secret.key" = "your_secret_key"
  }
  pv_name         = "tf_test_pv"
  pvc_name        = "tf_test_pvc"
  deployment_name = "tf_test_deployment"
  deployment_containers = [
    {
      name  = "nginx"
      image = "nginx:latest"
      volume_mounts = [
        {
          mount_path = "/data"
        }
      ]
    }
  ]
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Architecture

This example creates the following resources:

1. **VPC and Subnet** - Network infrastructure for the CCE cluster
2. **EIP and Bandwidth** - Public IP and bandwidth for the CCE cluster
3. **CCE Cluster** - Container cluster with specified configuration
4. **CCE Node** - Worker node in the cluster
5. **OBS Bucket** - Object storage bucket
6. **Kubernetes Secret** - Secret for OBS authentication
7. **Kubernetes PV** - Persistent volume manually configured for OBS
8. **Kubernetes PVC** - Persistent volume claim bound to the PV
9. **Kubernetes Deployment** - Application deployment with PVC mount

### Manual Storage Configuration

This example uses manual storage configuration where:

* **OBS Bucket** - An OBS bucket is manually created
* **Persistent Volume (PV)** - A Kubernetes PV resource is manually created and configured for the manually created OBS bucket
* **Persistent Volume Claim (PVC)** - A PVC is created and bound to the manually configured PV
* **Storage Binding** - The PV is bound to the PVC, making the OBS storage available to pods

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The CCE cluster is dependent on the VPC, subnet, and key pair
* The PV and PVC are manually configured to use the manually created OBS bucket
* The deployment mounts the PVC to provide persistent storage
* The example uses Ubuntu public images by default, but you can specify custom images
* The OBS bucket must be created manually before the PV and PVC are created

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.57.0 |
| kubernetes | >= 1.11.4 |
