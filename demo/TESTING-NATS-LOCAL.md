# Testing manually with nats CLI and a docker image

- Install [nats CLI Tool](https://docs.nats.io/nats-concepts/what-is-nats/walkthrough_setup)
- Run nats docker with docker run -d --name NATS --network host -p 4222:4222 nats -js
- Create key value stores as needed with `nats kv add <NAME>`
- Create object stores as needed with `nats object add <NAME>`
- Create streams as needed with TODO
- Run the main for each node 