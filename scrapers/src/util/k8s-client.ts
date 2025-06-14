import * as k8sClient from "@kubernetes/client-node"
import { V1Job, V1Pod } from "@kubernetes/client-node"
import { k8sConfig } from "./k8s-config.ts"

const namespace = k8sConfig.podNamespace

// Create Kubernetes clients factory
const kc = new k8sClient.KubeConfig()
kc.loadFromDefault()

// Create Kubernetes API clients
export const coreAPI  = kc.makeApiClient(k8sClient.CoreV1Api)
export const batchAPI = kc.makeApiClient(k8sClient.BatchV1Api)

// Obtain current pod information
export const pod: V1Pod = await coreAPI.readNamespacedPod({ namespace, name: k8sConfig.podName })

// Obtain controlling job
export const job: V1Job = await batchAPI.readNamespacedJob({ namespace, name: k8sConfig.jobName })
