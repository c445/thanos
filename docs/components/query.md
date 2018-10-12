# Query

The query component implements the Prometheus HTTP v1 API to query data in a Thanos cluster via PromQL.
It joins a Thanos cluster mesh and can access all data exposed by store APIs within it. It is fully stateless and horizontally scalable.

The querier needs to be passed an initial list of peers to through which it can join the Thanos cluster.
This can either be a repeated list of peer addresses or a DNS name, which it will resolve for peer IP addresses.

```
$ thanos query \
    --http-address     "0.0.0.0:9090" \
    --cluster.peers    "thanos-cluster.example.org" \
```

The query layer can deduplicate series that were collected from high-availability pairs of data sources such as Prometheus.
A fixed replica label must be chosen for the entire cluster and can then be passed to query nodes on startup.
Two or more series that have that are only distinguished by the given replica label, will be merged into a single time series. This also hides gaps in collection of a single data source.


```
$ thanos query \
    --http-address        "0.0.0.0:9090" \
    --query.replica-label "replica" \
    --cluster.peers       "thanos-cluster.example.org" \
```

## Expose UI on a sub-path

It is possible to expose thanos-query UI and optinally API on a sub-path.
The sub-path can be defined either statically or dynamically via an HTTP header.
Static path prefix definition follows the pattern used in Prometheus,
where `web.route-prefix` option defines HTTP request path prefix (endpoints prefix)
and `web.external-prefix` prefixes the URLs in HTML code and the HTTP redirect responces.

Additionally, Thanos supports dynamic prefix configuration, which
[is not yet implemented by Prometheus](https://github.com/prometheus/prometheus/issues/3156).
Dynamic prefixing simplifies setup when `thanos query` is exposed on a sub-path behind
a reverse proxy, for example, via a Kubernetes ingress controller
[Traefik](https://docs.traefik.io/basics/#frontends)
or [nginx](https://github.com/kubernetes/ingress-nginx/pull/1805).
If `PathPrefixStrip: /some-path` option or `traefik.frontend.rule.type: PathPrefixStrip`
Kubernetes Ingress annotation is set, then `Traefik` writes the stripped prefix into X-Forwarded-Prefix header.
Then, `thanos query --web.prefix-header=X-Forwarded-Prefix` will serve correct HTTP redirects and links prefixed by the stripped path.

## Deployment

## Flags

[embedmd]:# (flags/query.txt $)
```$
usage: thanos query [<flags>]

query node exposing PromQL enabled Query API with data retrieved from multiple
store nodes

Flags:
  -h, --help                     Show context-sensitive help (also try
                                 --help-long and --help-man).
      --version                  Show application version.
      --log.level=info           Log filtering level.
      --gcloudtrace.project=GCLOUDTRACE.PROJECT  
                                 GCP project to send Google Cloud Trace tracings
                                 to. If empty, tracing will be disabled.
      --gcloudtrace.sample-factor=1  
                                 How often we send traces (1/<sample-factor>).
                                 If 0 no trace will be sent periodically, unless
                                 forced by baggage item. See
                                 `pkg/tracing/tracing.go` for details.
      --grpc-address="0.0.0.0:10901"  
                                 Listen ip:port address for gRPC endpoints
                                 (StoreAPI). Make sure this address is routable
                                 from other components if you use gossip,
                                 'grpc-advertise-address' is empty and you
                                 require cross-node connection.
      --grpc-advertise-address=GRPC-ADVERTISE-ADDRESS  
                                 Explicit (external) host:port address to
                                 advertise for gRPC StoreAPI in gossip cluster.
                                 If empty, 'grpc-address' will be used.
      --grpc-server-tls-cert=""  TLS Certificate for gRPC server, leave blank to
                                 disable TLS
      --grpc-server-tls-key=""   TLS Key for the gRPC server, leave blank to
                                 disable TLS
      --grpc-server-tls-client-ca=""  
                                 TLS CA to verify clients against. If no client
                                 CA is specified, there is no client
                                 verification on server side. (tls.NoClientCert)
      --http-address="0.0.0.0:10902"  
                                 Listen host:port for HTTP endpoints.
      --cluster.address="0.0.0.0:10900"  
                                 Listen ip:port address for gossip cluster.
      --cluster.advertise-address=CLUSTER.ADVERTISE-ADDRESS  
                                 Explicit (external) ip:port address to
                                 advertise for gossip in gossip cluster. Used
                                 internally for membership only.
      --cluster.peers=CLUSTER.PEERS ...  
                                 Initial peers to join the cluster. It can be
                                 either <ip:port>, or <domain:port>. A lookup
                                 resolution is done only at the startup.
      --cluster.gossip-interval=<gossip interval>  
                                 Interval between sending gossip messages. By
                                 lowering this value (more frequent) gossip
                                 messages are propagated across the cluster more
                                 quickly at the expense of increased bandwidth.
                                 Default is used from a specified network-type.
      --cluster.pushpull-interval=<push-pull interval>  
                                 Interval for gossip state syncs. Setting this
                                 interval lower (more frequent) will increase
                                 convergence speeds across larger clusters at
                                 the expense of increased bandwidth usage.
                                 Default is used from a specified network-type.
      --cluster.refresh-interval=1m  
                                 Interval for membership to refresh
                                 cluster.peers state, 0 disables refresh.
      --cluster.secret-key=CLUSTER.SECRET-KEY  
                                 Initial secret key to encrypt cluster gossip.
                                 Can be one of AES-128, AES-192, or AES-256 in
                                 hexadecimal format.
      --cluster.network-type=lan  
                                 Network type with predefined peers
                                 configurations. Sets of configurations
                                 accounting the latency differences between
                                 network types: local, lan, wan.
      --http-advertise-address=HTTP-ADVERTISE-ADDRESS  
                                 Explicit (external) host:port address to
                                 advertise for HTTP QueryAPI in gossip cluster.
                                 If empty, 'http-address' will be used.
      --grpc-client-tls-secure   Use TLS when talking to the gRPC server
      --grpc-client-tls-cert=""  TLS Certificates to use to identify this client
                                 to the server
      --grpc-client-tls-key=""   TLS Key for the client's certificate
      --grpc-client-tls-ca=""    TLS CA Certificates to use to verify gRPC
                                 servers
      --grpc-client-server-name=""  
                                 Server name to verify the hostname on the
                                 returned gRPC certificates. See
                                 https://tools.ietf.org/html/rfc4366#section-3.1
      --web.route-prefix=""      Prefix for API and UI endpoints. This allows
                                 thanos UI to be served on a sub-path. This
                                 option is analogous to --web.route-prefix of
                                 Promethus.
      --web.external-prefix=""   Static prefix for all HTML links and redirect
                                 URLs in the UI query web interface. Actual
                                 endpoints are still served on / or the
                                 web.route-prefix. This allows thanos UI to be
                                 served behind a reverse proxy that strips a URL
                                 sub-path.
      --web.prefix-header=""     Name of HTTP request header used for dynamic
                                 prefixing of UI links and redirects. This
                                 option is ignored if web.external-prefix
                                 argument is set. Security risk: enable this
                                 option only if a reverse proxy in front of
                                 thanos is resetting the header. The
                                 --web.prefix-header=X-Forwarded-Prefix option
                                 can be useful, for example, if Thanos UI is
                                 served via Traefik reverse proxy with
                                 PathPrefixStrip option enabled, which sends the
                                 stripped prefix value in X-Forwarded-Prefix
                                 header. This allows thanos UI to be served on a
                                 sub-path.
      --query.timeout=2m         Maximum time to process query by query node.
      --query.max-concurrent=20  Maximum number of queries processed
                                 concurrently by query node.
      --query.replica-label=QUERY.REPLICA-LABEL  
                                 Label to treat as a replica indicator along
                                 which data is deduplicated. Still you will be
                                 able to query without deduplication using
                                 'dedup=false' parameter.
      --selector-label=<name>="<value>" ...  
                                 Query selector labels that will be exposed in
                                 info endpoint (repeated).
      --store=<store> ...        Addresses of statically configured store API
                                 servers (repeatable).
      --store.sd-files=<path> ...  
                                 Path to files that contain addresses of store
                                 API servers. The path can be a glob pattern
                                 (repeatable).
      --store.sd-interval=5m     Refresh interval to re-read file SD files. It
                                 is used as a resync fallback.
      --query.auto-downsampling  Enable automatic adjustment (step / 5) to what
                                 source of data should be used in store gateways
                                 if no max_source_resolution param is specified.

```
