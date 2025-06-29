import * as k8sClient from "@kubernetes/client-node"
import { config } from "../config.js"

const k8sNamespace = process.env.POD_NAMESPACE
if (!k8sNamespace) {
    throw new Error("POD_NAMESPACE env variable is not set")
}

// Create Kubernetes clients factory
const kc = new k8sClient.KubeConfig()
kc.loadFromDefault()

// Create Kubernetes API clients
const coreAPI  = kc.makeApiClient(k8sClient.CoreV1Api)
const appsAPI  = kc.makeApiClient(k8sClient.AppsV1Api)
const batchAPI = kc.makeApiClient(k8sClient.BatchV1Api)

// Obtain current pod information
const podName   = process.env.POD_NAME!
const namespace = process.env.POD_NAMESPACE!
const pod       = await coreAPI.readNamespacedPod({ namespace, name: podName })

// Obtain controlling replica-set
const replicaSetOwner = pod.metadata!.ownerReferences!.find(ref => ref.kind === "ReplicaSet")
if (!replicaSetOwner) {
    throw new Error("Pod is not owned by a ReplicaSet")
}
const replicaSet = await appsAPI.readNamespacedReplicaSet({ namespace, name: replicaSetOwner.name })

// Obtain controlling deployment
const deploymentOwner = replicaSet.metadata!.ownerReferences!.find(ref => ref.kind === "Deployment")
if (!deploymentOwner) {
    throw new Error("ReplicaSet is not owned by a Deployment")
}
const deployment = await appsAPI.readNamespacedDeployment({ namespace, name: deploymentOwner.name })

// Export APIs
export const k8s = {
    APIException: k8sClient.ApiException,
    namespace: config.k8s.namespace,
    api: {
        core: coreAPI,
        apps: appsAPI,
        batch: batchAPI,
    },
    pod,
    replicaSet,
    deployment,
}
