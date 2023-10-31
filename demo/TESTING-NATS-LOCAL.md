# Testing manually with nats CLI and a docker image

- Install [nats CLI Tool](https://docs.nats.io/nats-concepts/what-is-nats/walkthrough_setup)
- Run `nats docker with docker run -d --name NATS --network host -p 4222:4222 nats -js`

For emulating KAI creation of nats resources as defined in `app.yaml` do the following:
- Create key value stores as needed in centralized configuration with `nats kv add <BUCKET>`
- Create object stores as needed with `nats object add <OBJECT_STORE>`
- Create streams as needed for messaging with `nats stream add <STREAM>`
    - This command will ask you next the subject defined in the documents as `<OUTPUT>`, then select memory and confirm the rest of default settings
    - You can check the subscription is working with `nats subscribe <OUTPUT>`
- Run the main file for each process 

# Testing manually with minio and a docker image

- Run 
```
mkdir -p ~/minio/data
docker run     -p 9000:9000     -p 9090:9090     --name minio     -v ~/minio/data:/data     -e "MINIO_ROOT_USER=minioadmin"     -e "MINIO_ROOT_PASSWORD=minioadmin"     quay.io/minio/minio server /data --console-address ":9090"```
