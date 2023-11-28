# Testing manually with a minio docker image

- Run minio docker with `docker run     -p 9000:9000     -p 9090:9090     --name minio     -v ~/minio/data:/data     -e "MINIO_ROOT_USER=root"     -e "MINIO_ROOT_PASSWORD=password"     quay.io/minio/minio server /data --console-address ":9090"`

Another option is deleting the minio server initialization code in the persistent_storage file inside each sdk, but then the go.mod/pyproject.toml files need to be updated to point to local:
- In go this can be done with `replace github.com/konstellation-io/kai-sdk/go-sdk => ../../go-sdk` inside the demo process folder
- In python this can be done by replacing `runner = {develop = true, path = "../../py-sdk/runner"}` inside the demo process folder and then replacing `sdk = {path = "../sdk", develop = true}` in the `py-sdk/runner` folder