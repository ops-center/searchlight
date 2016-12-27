# Test Icinga2 Custom Plugins

### Configure Kubernetes Client
Provide  kubernetes cluster host address and auth information

Write `config.ini` file in `pkg/client/k8s/` directory

Example:

> config.ini

    host=https://127.0.0.1:6443/
    username=admin@kubernetes-cluster.com
    password=123456789

Note:

* Do note use leading, trailing space with key and value
* Client will load `config.ini` from `"$GOPATH/src/github.com/appscode/searchlight/pkg/client/k8s/config.ini"`
* Set ENV `APPSCODE_ENV` to `dev` or make it empty


### Build `hyperalert`

    ./hack/make.py build_cmd hyperalert

 This will create a binary in `dist/hyperalert/` directory


### Run Test

Run following command to test

* __component_status__

        go test -v github.com/appscode/searchlight/test -run ^TestComponentStatus$

* __json_path__

        go test -v github.com/appscode/searchlight/test -run ^TestJsonPath$

* __node_count__

        go test -v github.com/appscode/searchlight/test -run ^TestNodeCount$

* __node_status__

        go test -v github.com/appscode/searchlight/test -run ^TestNodeStatus$

> Following two will create temporary kubernetes objects to test

* __pod_exists__

        go test -v github.com/appscode/searchlight/test -run ^TestPodExists$

* __pod_status__

        go test -v github.com/appscode/searchlight/test -run ^TestPodStatus$
