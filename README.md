# ols-load-generator
Load generator tool for openshift lightspeed service (OLS).

## **Prerequisites**
### **Running on openshift cluster**
OLS deployed on openshift cluster. Please refer to the instructions [here](https://github.com/openshift/lightspeed-operator?tab=readme-ov-file#running-on-the-cluster).
### **Running on local machine (Optional)**
A running instance of OLS (to test). Please refer to the instructions [here](https://github.com/openshift/lightspeed-service?tab=readme-ov-file#installation).

## **Installation**
```
make build; make install
```
> **NOTE**: You might want to add **sudo** to the install command as it involves creating `ols-load-generator` binary in your $PATH.   

If running on openshift simply run `make` which will build and push the image to specified image registry.

## **Usage**
### **Usage on openshift platform**
Once we have OLS deployed and running on openshift cluster, in order to trigger the load test, deploy [assets/ols-load-generator.yaml](https://github.com/vishnuchalla/ols-load-generator/blob/perf-testing/assets/ols-load-generator.yaml) replaced with your corresponding envs values onto the cluster.

### **Envs**
* `OLS_TEST_HOST` - String indicating OLS endpoint to perform load testing.
* `OLS_TEST_RUNID`(Optional) - String specifying an unique ID. Will be helpful while comparing two runs or while looking at a specific run results.
* `OLS_TEST_AUTH_TOKEN` - OLS auth token string.
* `OLS_TEST_HIT_SIZE` - Indicates the total amount of requests to be sent as a part of the load testing.
* `OLS_TEST_RPS` - Requests per second to be sent to the OLS endpoint in parallel.

### **Example Usage**
```
oc apply -f ~/assets/ols-load-generator.yaml
```
Once applied it will create a job in the specified namespace and will start running the tests with above mentioned values. We can tail the logs in order to look at the benchmark results.

### **Usage on Local Machine**
```
NAME:
   ols-load-generator - A command-line tool to load test openshift lightspeed service (OLS).

USAGE:
   ols-load-generator [global options] command [command options] [arguments...]

DESCRIPTION:
   A command-line tool to load test openshift lightspeed service (OLS).

COMMANDS:
   attack   ols-load-test attack
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -D          print debugging logs (default: false)
   -W          quieter log output (default: false)
   --help, -h  show help

```
### Attack
```
NAME:
   ols-load-generator attack - ols-load-test attack

USAGE:
   ols-load-generator attack [command options] [arguments...]

DESCRIPTION:
   perform attack on ols endpoints

OPTIONS:
   --host value       --host localhost:6060/ (default: "http://localhost:6060/") [$OLS_TEST_HOST]
   --runid value      --runid f519d9b2-aa62-44ab-9ce8-4156b712f6d2 (default: "d016ff7f-8986-4e67-8b96-fb2254673687") [$OLS_TEST_RUNID]
   --authtoken value  --authtoken authtoken [$OLS_TEST_AUTH_TOKEN]
   --hitsize value    --hitsize 100 (default: 25) [$OLS_TEST_HIT_SIZE]
   --rps value        --rps 50 (default: 10) [$OLS_TEST_RPS]
   --help, -h         show help
```
### Example Usage
```
ols-load-generator attack --host https://127.0.0.1:9001 --runid random-uuid --authtoken 'auth-token' --hitsize 5 --rps 1
```