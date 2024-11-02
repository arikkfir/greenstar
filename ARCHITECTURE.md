# Architecture

```mermaid
flowchart TB
    user(("User"))
    user --> gateway
    subgraph k8s ["Kubernetes Cluster"]
        gateway["Gateway Controller"]
    %% Frontend
        gateway --> frontendHTTPRoute
        frontendHTTPRoute["Frontend HTTP Route"]
        frontendService["Frontend Service"]
        frontendHTTPRoute --> frontendService
        frontendService --> frontendPod1
        frontendService --> frontendPod2
        subgraph frontendDeployment["Frontend Deployment"]
            frontendPod1["Frontend Pod"]
            frontendPod2["Frontend Pod"]
        end
    %% Backend
        gateway --> backendHTTPRoute
        backendHTTPRoute["Backend HTTP Route"]
        backendService["Backend Service"]
        backendHTTPRoute --> backendService
        backendService --> backendPod1
        backendService --> backendPod2
        subgraph backendDeployment["Backend Deployment"]
            backendPod1["Backend Pod"]
            backendPod2["Backend Pod"]
        end
    %% Database
        neo4j[(Neo4j)]
        backendPod1 --> neo4j
        backendPod2 --> neo4j
        redis[(Redis)]
        backendPod1 --> redis
        backendPod2 --> redis
    end
```

## Open questions

- [ ] Verify Neo4j session semantics (are nested calls to `Service.getNeo4jSession(...)` coordinating in a single
  transaction?)
