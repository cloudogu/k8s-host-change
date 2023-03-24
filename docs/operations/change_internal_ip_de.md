# Konfiguration einer internen IP für eine Split-DNS-Konfiguration

Das `k8s-ecosystem` bietet die Möglichkeit das System in einer Split-DNS-Konfiguration zu betreiben.
Dabei kann ein externer Reverse-Proxy vor dem `k8s-ecosystem` betrieben werden, der unter der im CES konfigurierten
FQDN erreichbar ist. Die Verwendung einer internen IP verhindert, dass interne Requests der Dogus nicht über diesen 
externen Proxy geleitet werden.

Dieser Job bietet die Möglichkeit die interne IP während einer laufenden CES-Instanz zu konfigurieren.
Dabei werden die Registry-Werte `config/_global/use_internal_ip` und `config/_global/internal_ip` ausgewertet
und anschließend die Deployments aller Dogus angepasst.

Achtung diese Änderung bewirkt, dass alle Dogus neu gestartet werden.

Download und Anwendung der Job-Ressourcen:

[YAML-Ressourcen](https://dogu.cloudogu.com/api/v1/k8s/k8s/k8s-host-change)

```bash
kubectl apply -f <fileName>.yaml --namespace ecosystem
```

