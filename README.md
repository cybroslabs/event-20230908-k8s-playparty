# Kubernetes Play Party

Pátek, 8. září, 2023

## Otázky

Ptejte se pořád, hned, na nic nečekejte! No možná se přihlásit, ať si neskáčeme do řeči... 😄

## Agenda

- Slovníček pojmů
- Architektura Kubernetes cluster
  - Control plane (master)
  - etcd
  - worker node
    - kubelet
    - kube-proxy
    - container runtime
  - dns v clusteru  
- Život clusteru
  - high availability
    - mastery
    - etcd
  - přidávání nodů
  - odebírání nodů
    - bezpečný postup odebírání nodů
  - selhání nodu (crash 💥)
  - monitoring nodů
    - co řeší kubernetes
    - co bych _měl_ řešit já
      - Prometheus
      - node exporter
      - node problem detector
- Běžné k8s objekty, se kterými se potkáme
  - _Node_
  - _Pod_
    - health checks
      - startup probe
      - liveness probe
      - readiness probe
  - _Deployment_
  - _Service_
  - _Ingress_
  - _StatefulSet_
  - _Job_ a _CronJob_
- Helm  
- Demo time, aneb pojďme něco rozbít
  - HTTP věc
  - PostgreSQL
    - persistent data storage
      - PVC
      - PV
    - database replication
    - replication strategies
    - failover
    - failover domains
    - 2 DC setup with two replicas within each DC
  - RabbitMQ
    - different data replication strategy than from Postgres
    - durable queues - replication of data across replicas

## Slovníček pojmů

- machine/stroj: virtuální server / fyzický hardware (baremetal) s OS
- node: stroj, který je člen k8s clusteru
- cluster: skupina propojených strojů (nodů)
- ingress: směr dovnitr |<-
- engress: směr ven |->
- docker/container image: obraz aplikace a jejích závislostí
- container: běžící image (proces)

## Architektura Kubernetes clusteru

![](/img/components-of-kubernetes.png)

### Control plane / master

- Kube API
- Cloud Controller Manager
- Controller Manager
- Scheduler
- etcd

### etcd

Mozek Kubernetes, drží stav cluster.

```sh
key=value
```

- distribuované - napříč N stroji
- musí se shnodnout většina -> 1, 3, 5,... prostě vždy lichý počet replik

### Worker node

Kde běží naše aplikace, může být zároveň master nodem, pokud zde běží komponenty z control plane.

#### Kubelet

Běží na nodu, statrá se o komunikaci s kube api, scheduler mu říká, co má spustit, co má vypnout...

Komunikuje s container runtime (containerd), skrz CRI (specifikace).

#### container runtime

To co reálně spouští (a stahuje) image.

- containerd
- cri-o

Dřív i docker (docker-shim).

### DNS v clusteru

Kubernetes používá interní DNS servery, abychom nemuseli řešit networking na úrovni IP adress, ale máme DNS.

Defaultní doména clusteru: `cluster.local`

Doména Servicy: `<name>.<namespace>` nebo `<name>.<namespace>.svc.cluster.local`.

## Život clusteru

Jak žít a pracovat s nody a clusterem potom, co jej vytvořím. Tzv. "day two operations".

### High availability

Neběží mi jen jeden proces...

Ideálně 3 a víc instancí, protože když jedna vypadne, load se distribuuje napříč zbývajícími instancemi.

#### Master

Stateless, takže 3+ a škáluju dle potřeby.

Leader election (scheduler) se řeší automaticky a my to neřešíme.

#### etcd

3+ replik, vždy lichý počet!!! Jinak hrozí, tzv. split-brain - protože nemám quorum (neshodnu se).

### Přidávání nodů

Nový stroj, který má nainstalované komponenty a má credentials, kterými se autentizuje do clusteru (kube api).

### Odebírání nodů

TL;DR: Vypnu a samžu z clusteru.

#### Bezpečné odebírání nodů

1. cordon - nový workload se přestane plánovat pro konkrétní nody
2. drain - vypneme workload na nodu (scheduler jej přesune jinam, pokud je kapacita)
3. graceful period - počkám, až se všechno vypne
4. node shutdown - bezpečně vypínám node
5. odeberu node z kube api

### Selhání nodu (crash 💥)

Situace: Node (stroj) náhle přestane běžet, přestane na něm běžet workload a je na clusteru, aby se s tím vyrovnal.

Selže node health check na kubelet -> kube-api dá vědět scheduleru, že node je K.O -> scheduler bere workload jako mrtvý (crashed) a přesune jej jinam.

### Monitoring nodů

#### Co řeší Kubernetes

Kube API komunikuje s kubelet-em na každém nodu a dělá health checky.

Pokud healthcheck selže (n krát), bere node jako mrtvý.

#### Co bych _měl_ řešit já

Monitoring vytížení - plánování a případně znám přířiny crashů. Počet běžících podů...

##### Prometheus

Prometheus je monitorovací systém, který umí sledovat metriky napříč procesy (sbírá je přes HTTP, pull-based).

Zároveň vyhodnocuje pravidla a pokud je metrika poruší, notifikuje alertmanager, aby alertoval do příslušných kanálů.

##### Node Exporter

Prometheus exporter.

Exportuje informace o stroji.

Využití RAM, využití CPU, disk, slow disk, disk I/O, síťový provoz, network I/O,...

-> alertovací pravidla

##### Node Problem Detector

[GitHub: Node Problem Detector](https://github.com/kubernetes/node-problem-detector)

> - Infrastructure daemon issues: ntp service down;
> - Hardware issues: Bad CPU, memory or disk;
> - Kernel issues: Kernel deadlock, corrupted file system;
> - Container runtime issues: Unresponsive runtime daemon;
> - ...

Prometheus exporter -> alerting pravidla

## Běžné k8s objekty, se kterými se potkáme

### Node

Kubernetes si drží objekt nodu (stroje), má k němu přiřazané labely, které můžeme využít při schedulingu (affinity a jiné).

Také jde nodu přidat "taint" a omezovat scheduling podů jinak - další cestou.

### Pod

Nejmenší jednotka, kterou v Kubernetes dokážu spustit.

Scheduler plánuje distribuci podů s jejich omezeními (compute resources, network, topology constraints, affinity, priority classes,...).

Nejsou persistentní, smazaný pod = smazaný pod. Sám se nerestartuje.

#### Health checks

##### Startu probe

"Než ta Java najede"

##### Liveness probe

Žije ten proces nebo je třeba jej restartovat?

##### Readiness probe

Proces sice žije, ale možná je aplikace přetížená -> není ready přijímat další traffic.

Vyřadím tedy Pod z rotace, kam load balancer posílá požadavky (řeší Service) a nechá Pod se sebrat.

Má limity, pokud se nespraví sám za nějaký čas -> restart.

### Service

Abstrakce a zároveň jednotný endpoint, jak po síti komunikovat s Pody.

### Ingress

Vystavení Service do veřejného internetu, s nějakými pravidly pro routing v rámci clusteru, HTTPS,...

### Deployment

Wrapper Podů, slouží k tomu, že říkám kolik replik chci (kolik Podů), Pody nezanikají.

Deployment > ReplicaSet > Pod

### StatefulSet

Stejně jako Deployment se stará o počet replik a to, že tvoří nové Pody po jejich smazání.

Zároveň má v jistých věcech jiné chování, než Deployment, které lépe svědší aplikacím, které drží stav a jednotlivé repliky nejsou mezi sebou zaměnitelné.

Například databáze nebo fronty. Nebo pokud mám striktně definovanou posloupnost, jak mohu zapínat/vypínat repliky (přidání, odebírání).

StatefulSet > Pod

### Job a CronJob

Jednorázové spuštění, umí konkurenci (víc vedle sebe), restart policy.

Job > Pod
CronJob > Job > Pod

## Helm

- Go templates (`{{ . }}`)

"parametrizovatelný YAML" s pár věcmi navíc

- for each (range)
- if & if else
- proměnné
- filtry/makra