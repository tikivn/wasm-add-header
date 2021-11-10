package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// Other examples can be found at https://github.com/tetratelabs/proxy-wasm-go-sdk/tree/v0.0.8/examples

func main() {
	proxywasm.SetNewHttpContext(newHttpContext)
}

type httpHeaders struct {
	// we must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultHttpContext
	contextID uint32
}

func newHttpContext(rootContextID, contextID uint32) proxywasm.HttpContext {
	return &httpHeaders{contextID: contextID}
}

// override
func (ctx *httpHeaders) OnHttpRequestHeaders(int, bool) types.Action {
	cluster_path := []string{"node", "cluster"}
	cluster, err := proxywasm.GetProperty(cluster_path)
	if err != nil {
		proxywasm.LogInfof("Get cluster error: %v", err)
	}
	cluster_name := string(cluster)
	proxywasm.LogInfof("Cluster Name: %s", cluster_name)
	if err := proxywasm.SetHttpRequestHeader("tiki-source-app", cluster_name); err != nil {
		proxywasm.LogCriticalf("failed to set response cluster header: %v", err)
	}
	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	if err := proxywasm.SetHttpResponseHeader("hello", "wasm-world"); err != nil {
		proxywasm.LogCriticalf("failed to set response header: %v", err)
	}
	cluster_path := []string{"node", "cluster"}
	cluster, err := proxywasm.GetProperty(cluster_path)
	if err != nil {
		proxywasm.LogInfof("Get cluster error: %v", err)
	}
	cluster_name := string(cluster)
	proxywasm.LogInfof("Cluster Name: %s", cluster_name)
	if err := proxywasm.SetHttpResponseHeader("tiki-source-app", cluster_name); err != nil {
		proxywasm.LogCriticalf("failed to set response cluster header: %v", err)
	}
	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
