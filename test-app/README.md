# Test app

The test app is a pretty simple nodejs application using the express.js framework.
It counts up a number every 100ms and implements the following endpoints:

- `GET /state`
  Gets the current state in the following form:
```json
{
  "counter": 5
}
```
- `POST /state` Sets the current state. Takes:
```json
{
  "counter": 5
}
```
- `GET /kill` Kills the application
- `GET /restore` returns the current state and if the `POST /state` was ever called
```json
{
  "counter": 5,
  "restore": false
}
```
- `GET /health` Returns empty on success


## Usage

First, set an environment variable `IMAGE_REGISTRY` to point to your container image registry and replace the `IMAGEREGISTRY` in `test-app-deployment.yaml` with the value of that environment variable. You can use the following two commands to do so:

```
export IMAGE_REGISTRY=<your-container-image-registry>
sed -i "s/IMAGEREGISTRY/${IMAGE_REGISTRY}/g" test-app-deployment.yaml
```

Then, you can build the test app container image, push it to your image registry and deploy it to your Kubernetes cluster as follows:
```
make build push
make deploy
```