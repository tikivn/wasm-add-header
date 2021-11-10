k -n platform create secret generic go-envoy-wasm --from-file main.wasm

sidecar.istio.io/userVolume=[{"name": "wasm","secret": {"defaultMode": 256,"optional": false,"secretName": "go-envoy-wasm"}}]
sidecar.istio.io/userVolumeMount=[{"mountPath":"/tmp/wasm","name":"wasm"}]

echo "
---
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: wasm-console-app
  namespace: cicd
spec:
  configPatches:
  - applyTo: HTTP_FILTER
    match:
      context: SIDECAR_INBOUND
      listener:
        filterChain:
          filter:
            name: envoy.http_connection_manager
            subFilter:
              name: envoy.router
    patch:
      operation: INSERT_BEFORE
      value:
        name: envoy.filters.http.wasm
        typedConfig:
          '@type': type.googleapis.com/udpa.type.v1.TypedStruct
          typeUrl: type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm
          value:
            config:
              configuration:
                '@type': type.googleapis.com/google.protobuf.StringValue
                value: ""
              name: wasm-console-app
              rootId: root_id
              vmConfig:
                code:
                  local:
                    filename: /tmp/wasm/main.wasm
                runtime: envoy.wasm.runtime.v8
                vmId: wasm-console-app
  workloadSelector:
    labels:
      workload.user.cattle.io/workloadselector: deployment-cicd-nginx
" | k apply -f -