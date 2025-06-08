# Proof Program Registry

This Helm chart deploys the Proof Program Registry application along with its required IPFS dependency.

## Prerequisites

* Kubernetes 1.19+
* Helm 3.2+
* A default StorageClass available in your Kubernetes cluster for PersistentVolume provisioning. You can check this by running `kubectl get storageclass`.

## Installation



**Install the Chart**

To install the chart with the release name `proof-program-registry`, run the following command. This will deploy the registry and a dedicated IPFS node.

```bash
helm install proof-program-registry ./proof-program-registry/
```

After the installation, Kubernetes will begin provisioning the necessary resources. You can track the status of your deployment by running:

```bash
kubectl get pods -l app.kubernetes.io/name=proof-program-registry -w
```

## Accessing the Application

To interact with the registry from your local machine (for example, using Postman or `curl`), you need to forward the service port from the Kubernetes cluster to your local machine.

**Forward the port:**

The following command forwards the application's port (e.g., 8080) to your local machine.

```bash
kubectl port-forward svc/proof-program-registry 8080:8080
```
You can now send requests to `http://localhost:8080`.


## Uninstallation

To uninstall and delete all the Kubernetes components associated with the `proof-program-registry` release, run the following command:

```bash
helm uninstall proof-program-registry
```

**Important:** This command will not delete the PersistentVolumeClaim (PVC) by default, ensuring that your IPFS data is not lost. To delete the PVC and release the storage, you must manually delete it:

```bash
kubectl delete pvc -l app.kubernetes.io/name=proof-program-registry
