# Go CLI For Measure

## Setup
``` sh
asdf plugin add golang
asdf plugin add golangci-lint
asdf install
```

curl -G https://prometheus.state.api.cresplanex.org/api/v1/query_range \
  --data-urlencode "query=sum by(namespace) (rate(container_cpu_usage_seconds_total{container!='POD',container!='',namespace!='kube-system', image=~'docker.io/ablankz/.*',  image!='docker.io/ablankz/debezium:1.0.0', job='kubelet'}[3m]))" \
  --data-urlencode "start=$(date -u -d '30 minutes ago' +%s)" \
  --data-urlencode "end=$(date -u +%s)" \
  --data-urlencode "step=30"