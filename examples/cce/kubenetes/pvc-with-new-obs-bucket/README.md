# Create a PVC with a new OBS bucket under the CCE environment

This example provides best practice code for using Terraform to create a Persistent Volume Claim (PVC) (via the
Kubernetes provider, which automatically creates a new OBS bucket) and mount the PVC to a container under the
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
* `secret_name` - The name of the Kubernetes secret for OBS authentication
* `secret_data` - The data of the Kubernetes secret (sensitive)
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
* `namespace_name` - The name of the CCE namespace (default: "default")
* `secret_labels` - The labels of the Kubernetes secret (default: {"secret.kubernetes.io/used-by" = "csi"})
* `secret_type` - The type of the Kubernetes secret (default: "cfe/secure-opaque")
* `pvc_obs_volume_type` - The type of the OBS volume of the PVC (default: "standard")
* `pvc_fstype` - The instance type of the PVC (default: "s3fs")
* `pvc_access_modes` - The access modes of the PVC (default: ["ReadWriteMany"])
* `pvc_storage` - The storage of the PVC (default: "1Gi")
* `pvc_storage_class_name` - The name of the storage class of the PVC (default: "csi-obs")
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
  secret_name  = "tf_test_secret"
  secret_data  = {
    "access.key" = "your_access_key"
    "secret.key" = "your_secret_key"
  }

  pvc_name              = "tf_test_pvc"
  deployment_name       = "tf_test_deployment"
  deployment_containers = [
    {
      name          = "nginx"
      image         = "nginx:latest"
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
5. **Kubernetes Secret** - Secret for OBS authentication
6. **Kubernetes PVC** - Persistent volume claim with OBS storage class
7. **Kubernetes PV** - Persistent volume (automatically created by PVC)
8. **OBS Bucket** - Object storage bucket (automatically created by PV)
9. **Kubernetes Deployment** - Application deployment with PVC mount

### Dynamic Storage Provisioning

When the PVC is created with the `csi-obs` storage class, the following resources are automatically provisioned:

* **Persistent Volume (PV)** - A Kubernetes PV resource is automatically created and bound to the PVC
* **OBS Bucket** - A new OBS bucket is automatically created to provide the actual storage backend
* **Storage Binding** - The PV is automatically bound to the PVC, making the storage available to pods

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The CCE cluster is dependent on the VPC, subnet, and key pair
* The PVC uses the CSI OBS driver for dynamic storage provisioning
* The deployment mounts the PVC to provide persistent storage
* The example uses Ubuntu public images by default, but you can specify custom images
* OBS bucket is created automatically when the PVC is created

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.57.0 |
| kubernetes | >= 1.6.2 |
