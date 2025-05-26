# AutoRestartPod

A Kubernetes operator that provides automated pod restart functionality based on cron schedules.

## Description

AutoRestartPod is a Kubernetes operator built with Kubebuilder that enables automated pod restarts based on cron schedules. It allows cluster administrators and application owners to define custom restart policies for their pods without manual intervention.

This controller is useful for scenarios where applications require periodic restarts to:
- Clear memory leaks
- Apply configuration changes that require restarts
- Refresh connections to external services
- Perform routine maintenance tasks
- Implement rolling restarts for stateless applications

The controller uses standard cron syntax for scheduling, supports time zone configuration, and targets pods using Kubernetes label selectors.

## Getting Started

### Prerequisites
- go version v1.24.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/autorestartpod:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/autorestartpod:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

## CRD Specification

The AutoRestartPod CRD has the following structure:

```yaml
apiVersion: stable.crazyfrank.com/v1
kind: AutoRestartPod
metadata:
  name: string
  namespace: string
spec:
  # Cron schedule in standard cron format (required)
  # Examples: 
  # - "0 3 * * * ?" (every day at 3:00 AM)
  # - "0 */6 * * * ?" (every 6 hours)
  # - "0 0 * * 0 ?" (every Sunday at midnight)
  schedule: string
  
  # Standard Kubernetes label selector (required)
  # Used to identify which pods to restart
  selector:
    matchLabels:
      key1: value1
      key2: value2
    matchExpressions:
      - {key: key3, operator: In, values: [value3, value4]}
  
  # Time zone for the schedule (optional, defaults to UTC)
  # Examples: "UTC", "America/New_York", "Asia/Shanghai"
  timeZone: string
  
status:
  # The last time pods were restarted by this controller
  lastRestartTime: timestamp
```

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

### Usage Examples

#### Basic AutoRestartPod Example

Create an AutoRestartPod resource to restart all nginx pods in the default namespace every day at 3:00 AM:

```yaml
apiVersion: stable.crazyfrank.com/v1
kind: AutoRestartPod
metadata:
  name: nginx-daily-restart
  namespace: default
spec:
  schedule: "0 3 * * * ?"  # Every day at 3:00 AM
  selector:
    matchLabels:
      app: nginx
  timeZone: "UTC"  # Optional timezone
```

#### Multiple Restart Schedules

You can create multiple AutoRestartPod resources for different restart schedules or target different pods:

```yaml
# Restart frontend pods on weekdays at 2:00 AM
apiVersion: stable.crazyfrank.com/v1
kind: AutoRestartPod
metadata:
  name: frontend-weekday-restart
  namespace: prod
spec:
  schedule: "0 2 * * 1-5 ?"  # Every weekday at 2:00 AM
  selector:
    matchLabels:
      component: frontend
  timeZone: "America/New_York"
```

```yaml
# Restart backend pods every 12 hours
apiVersion: stable.crazyfrank.com/v1
kind: AutoRestartPod
metadata:
  name: backend-12h-restart
  namespace: prod
spec:
  schedule: "0 */12 * * * ?"  # Every 12 hours
  selector:
    matchLabels:
      component: backend
  timeZone: "Asia/Shanghai"
```

#### Check Status

You can check the status of your AutoRestartPod resource:

```sh
kubectl get autorestartpod nginx-daily-restart -o yaml
```

The status section will show the last restart time:

```yaml
status:
  lastRestartTime: "2025-05-26T03:00:05Z"
```

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/autorestartpod:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/autorestartpod/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Contributing

Contributions to the AutoRestartPod project are welcome! Here's how you can contribute:

1. **Report Issues**: If you find bugs or have feature requests, please open an issue on the project repository.

2. **Submit Pull Requests**: 
   - Fork the repository
   - Create a feature branch (`git checkout -b feature/amazing-feature`)
   - Commit your changes (`git commit -m 'Add some amazing feature'`)
   - Push to the branch (`git push origin feature/amazing-feature`)
   - Open a Pull Request

3. **Documentation**: Help improve documentation, examples, or tutorials.

4. **Testing**: Add more test cases or improve existing tests.

### Development Environment Setup

1. Clone the repository:
   ```sh
   git clone https://github.com/crazyfrankie/autorestart-operator.git
   cd autorestart-operator
   ```

2. Install dependencies:
   ```sh
   go mod download
   ```

3. Run tests:
   ```sh
   make test
   ```

4. Run the controller locally:
   ```sh
   make install
   make run
   ```

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

