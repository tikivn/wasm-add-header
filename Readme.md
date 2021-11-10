# description
this module wasm demo how to use wasm module to add source workload name to header of outbound request.
every request make from workload will add header: tiki-source-app: <workload_name>.<namespace> to request then send to dest service
# deploy
deploy.md