# [MTB-PL1-BC-HI-3] Block use of host networking and ports

**Profile Applicability:**

Level 1

**Type:**

Behavioral

**Category:**

Host Isolation

**Description:**

Tenants should not be allowed to use host networking and host ports for their workloads.

**Rationale:**

Using `hostPort` and `hostNetwork` allows tenants workloads to share the host networking stack allowing potential snooping of network traffic across application pods

**Audit:**

Create a pod defining a container using a host port. The pod creation must fail.

Create a pod defining a container using a host network. The pod creation must fail.
