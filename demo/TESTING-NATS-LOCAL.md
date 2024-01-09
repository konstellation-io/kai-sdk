# Testing manually with nats CLI and a docker image

- Install [nats CLI Tool](https://docs.nats.io/nats-concepts/what-is-nats/walkthrough_setup)
- Run nats docker with `docker run -d --name NATS --network host -p 4222:4222 nats -js`

For emulating KAI creation of nats resources as defined in `app.yaml` do the following:
- Create key value stores as needed in centralized configuration with `nats kv add <BUCKET>`
- Create object stores as needed with `nats object add <OBJECT_STORE>`
- Create streams as needed for messaging with `nats stream add <STREAM>`
    - This command will ask you next the subject defined in the documents as `<OUTPUT>`, then select memory and confirm the rest of default settings
    - You can check the subscription is working with `nats subscribe <OUTPUT>`
- Run the main for each node 