## KUBEAPI
KubeAPI is used to generate client-go style APIs for crd resources.


### Install
```sh
go get github.com/seamounts/kubeapi
```

### Usage
Development kit for building Kubernetes extensions and tools.

Provides libraries and tools to create new projects, APIs and controllers.
Includes tools for packaging artifacts into an installer container.

Typical project lifecycle:

- initialize a project:
    ```sh
    kubeapi init --domain example.com --license apache2 --owner "The Kubernetes authors"
    ```
  

- create one or more a new resource APIs and add your code to them:
    ```sh
    kubeapi create api --group <group> --version <version> --kind <Kind>
    ```
After the scaffold is written, api will run make on the project.
