# Job and CronJob Support

Jobs and CronJobs are currently enabled as part of the ClowdApp spec. The
``jobs`` field contains a list of all currently defined jobs. The spec for a 
job is documented in the [Clowder API reference](https://redhatinsights.github.io/clowder/clowder/dev/api_reference.html#k8s-api-github-com-redhatinsights-clowder-apis-cloud-redhat-com-v1alpha1-clowdjobinvocation).

Jobs and CronJobs are split by a ``schedule`` field inside your job. If the job
has a ``schedule``, it is assumed to be a CronJob. If not, Clowder runs your 
job as a standard Job resource. Note that Jobs run as soon as they are applied. 

Jobs that need to be run at some arbitrary point in the future are run by a 
ClowdJobInvocation.

## Invoking Jobs via ClowdJobInvocation

Jobs can be triggered by applying a ``ClowdJobInvocation`` CRD to the cluster. 
Clowder will read the resource and run the specified Jobs.

Below is an example of a ClowdJobInvocation, or CJI for short. It is followed 
by a sample ClowdApp with the definition of the "curl" job. Notice that the 
``appName`` in the CJI matches the ``name`` of the ClowdApp. The ``jobs`` list
also corresponds to the ``jobs`` field in the ClowdApp. ``curl`` is a job 
hosted by the ClowdApp and the CJI will use that to find all the data it needs 
to invoke and apply the job. 

```yaml
apiVersion: cloud.redhat.com/v1alpha1
kind: ClowdJobInvocation
metadata:
  name: tester
spec:
  appName: sample-app
  jobs:
    - curl
```

```yaml
apiVersion: v1
kind: Template
metadata:
  name: sample-app
objects:
- apiVersion: cloud.redhat.com/v1alpha1
  kind: ClowdApp
  metadata:
    name: sample-app
  spec:
    envName: env-debugger
    deployments:
    - name: api
      minReplicas: 1
      webServices:
        public:
          enabled: true
      podSpec:
        image: quay.io/bholifie/simpleserver
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: 8000
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 5
    jobs:
      - name: hello-twenty
        schedule: "*/20 * * * *"
        parallelism: 2
        completions: 2
        podSpec:
          name: hello
          image: busybox
          args:
          - /bin/sh
          - -c
          - date; echo Hello from the Cron Job
      - name: curl
        podSpec:
          name: getter
          image: busybox
          args:
          - /bin/sh
          - -c
          - wget sample-app-api.debugger.svc:8000/
      - name: coming-soon
        podSpec:
          name: hello
          image: busybox
          args:
          - /bin/sh
          - -c
          - date; echo I'm ready!
```

To apply a CJI, run  ``oc apply -f cji.yml``

A CJI can then be checked by ``oc get cji``

## Running IQE Tests with ClowdJobs

Part of the mission for jobs was to empower developers to run the full suite
of testing on their local machine. Using ClowdJobInvocations, developers can
now run smoke tests locally and on a remote cluster. In order to get everything
setup correctly for the full smoke tests, we need to do the following:

1. Ensure your app's iqe plugin is configured to read from cdappconfig. Please feel 
free to use this [inventory MR as a reference](https://gitlab.cee.redhat.com/insights-qe/iqe-host-inventory-plugin/-/merge_requests/514/diffs). 
2. Use `bonfire` to deploy your app into an ephemeral namespace.
3. Use `bonfire deploy-iqe-cji` to deploy a CJI into the namespace.
