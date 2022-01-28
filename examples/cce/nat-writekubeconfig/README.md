# nat-writekubeconfig

This example is demo create a cce cluster,and this cluster will use nat to access internat,Node does not need to be assigned EIP.

Use local_file provider to write local kuberconfig file.When cce cluster created succesfull,you cao use Kubectl to connect clusters,no need to manually write kubeconfig file.

When use kubectl to connect CCE,You need to switch use-context depending on the network between your local and Apiserver.
```bash
kubectl config use-context internal
kubectl config use-context external
```