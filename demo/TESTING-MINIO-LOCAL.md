# Testing manually with a minio docker image

- Run minio docker with `docker run     -p 9000:9000     -p 9090:9090     --name minio     -v ~/minio/data:/data     -e "MINIO_ROOT_USER=minio_user"     -e "MINIO_ROOT_PASSWORD=minio_password"     quay.io/minio/minio server /data --console-address ":9090"`

Another option is deleting the minio server initialization code in the persistent_storage file inside each sdk, but then the go.mod/pyproject.toml files need to be updated to point to local:
- In go this can be done with `replace github.com/konstellation-io/kai-sdk/go-sdk => ../../go-sdk` inside the demo process folder
- In python this can be done by replacing `runner = {develop = true, path = "../../py-sdk/runner"}` inside the demo process folder and then replacing `sdk = {path = "../sdk", develop = true}` in the `py-sdk/runner` folder

A third option is removing the authentication code in model registry and persistent storage and changing the credentials when creating the client to 

```
access_key=v.get_string("minio.client_user"),
secret_key=v.get_string("minio.client_password"),
```

and then creating in the demo or here a bucket

```
self.minio_client.make_bucket(self.minio_bucket_name, location="us-east-1", object_lock=True)
```