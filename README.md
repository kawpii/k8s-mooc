# K8s MOOC

## Exercises

### Chapter 2
- [1.1](https://github.com/kawpiya/k8s-mooc/tree/1.1)
- [1.2](https://github.com/kawpiya/k8s-mooc/tree/1.2)
- [1.3](https://github.com/kawpiya/k8s-mooc/tree/1.3)
- [1.4](https://github.com/kawpiya/k8s-mooc/tree/1.4)
- [1.5](https://github.com/kawpiya/k8s-mooc/tree/1.5)
- [1.6](https://github.com/kawpiya/k8s-mooc/tree/1.6)
- [1.7](https://github.com/kawpiya/k8s-mooc/tree/1.7)
- [1.8](https://github.com/kawpiya/k8s-mooc/tree/1.8)
- [1.9](https://github.com/kawpiya/k8s-mooc/tree/1.9)
- [1.10](https://github.com/kawpiya/k8s-mooc/tree/1.10)
- [1.11](https://github.com/kawpiya/k8s-mooc/tree/1.11)
- [1.12](https://github.com/kawpiya/k8s-mooc/tree/1.12)
- [1.13](https://github.com/kawpiya/k8s-mooc/tree/1.13)

### Chapter 3
- [2.1](https://github.com/kawpiya/k8s-mooc/tree/2.1)
- [2.2](https://github.com/kawpiya/k8s-mooc/tree/2.2)
- [2.3](https://github.com/kawpiya/k8s-mooc/tree/2.3)
- [2.4](https://github.com/kawpiya/k8s-mooc/tree/2.4)
- [2.5](https://github.com/kawpiya/k8s-mooc/tree/2.5)
- [2.6](https://github.com/kawpiya/k8s-mooc/tree/2.6)
- [2.7](https://github.com/kawpiya/k8s-mooc/tree/2.7)
- [2.8](https://github.com/kawpiya/k8s-mooc/tree/2.8)
- [2.9](https://github.com/kawpiya/k8s-mooc/tree/2.9)
- [2.10](https://github.com/kawpiya/k8s-mooc/tree/2.10)

### Chapter 4
- [3.1](https://github.com/kawpiya/k8s-mooc/tree/3.1)
- [3.2](https://github.com/kawpiya/k8s-mooc/tree/3.2)
- [3.3](https://github.com/kawpiya/k8s-mooc/tree/3.3)
- [3.4](https://github.com/kawpiya/k8s-mooc/tree/3.4)
- [3.5](https://github.com/kawpiya/k8s-mooc/tree/3.5)
- [3.6](https://github.com/kawpiya/k8s-mooc/tree/3.6)
- [3.7](https://github.com/kawpiya/k8s-mooc/tree/3.7)
- [3.8](https://github.com/kawpiya/k8s-mooc/tree/3.8)
- [3.9](https://github.com/kawpiya/k8s-mooc/tree/3.9)
- [3.10](https://github.com/kawpiya/k8s-mooc/tree/3.10)
    Uses workload identity
- [3.11](https://github.com/kawpiya/k8s-mooc/tree/3.11)
- [3.12](https://github.com/kawpiya/k8s-mooc/tree/3.12)

#### Postgres with PVs
##### Pros
 - Full control over the installation, including config.
 - Lower cost.

##### Cons
 - Have to manage ops (backups, recovery, failover etc.), so can be unreliable.
 - Storage is node and zone constrained.
 - Data loss with misconfigurations.

#### DBaaS
##### Pros
 - Fully manages no Ops.
 - Runs outside K8s.
 - Scalability and performance built-in by experts.
 - SLA guarantees.
 - Multiple AZ support.

##### Cons
 - Possible vendor lock in
 - High latency possible since it can be outside cluster network.

There is bit of middle ground to reduce ops burden by using a Postgres operator like CloudNativePg, Zalando/postgres-operator etc. But still there are cons to be aware of like handling node/disk failures, still have to monitor.

### Chapter 5

- [4.1](https://github.com/kawpiya/k8s-mooc/tree/4.1)
- [4.2](https://github.com/kawpiya/k8s-mooc/tree/4.2)
- [4.3](https://github.com/kawpiya/k8s-mooc/tree/4.3)

```
sum(kube_pod_info{created_by_kind="StatefulSet", namespace="prometheus"})
```

### Chapter 6
