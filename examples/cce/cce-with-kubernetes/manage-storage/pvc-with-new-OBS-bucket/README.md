# cce-with-kubernetes

In this example, we will create a CCE cluster with EIP,
then we will use the Kubernetes provider to connect and manage the cluster.
We authenticate to the cluster with certificate.
We create a PVC, which will create a new OBS automatically.

For more information about using an OBS bucket through a dynamic PV,
[please see document](https://support.huaweicloud.com/intl/en-us/usermanual-cce/cce_10_0630.html).
For more information about the Kubernetes PVC resource,
[please see document](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/persistent_volume_claim).
For more information about the Kubernetes PV resource,
[please see document](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/persistent_volume).
