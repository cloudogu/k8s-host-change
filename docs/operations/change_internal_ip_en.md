# Configuration of an internal IP for a split DNS configuration.

The `k8s-ecosystem` offers the possibility to operate the system in a split DNS configuration.
In this case, an external reverse proxy can be operated in front of the `k8s-ecosystem`, which can be accessed under the
FQDN configured in the CES. The use of an internal IP prevents internal requests from the Dogus being routed via 
this external proxy.

This job offers the possibility to configure the internal IP while a CES instance is running.
The registry values `config/_global/use_internal_ip` and `config/_global/internal_ip` will be evaluated
and then the deployments of all dogus are adjusted.

Attention this change causes all dogus to be restarted.

Download and apply the job resources:

<!-- markdown-link-check-disable -->
[YAML resources](https://dogu.cloudogu.com/api/v1/k8s/k8s/k8s-host-change)

```bash
kubectl apply -f <fileName>.yaml --namespace ecosystem
```
