#!/bin/sh

minikube image build -t konstellation/demo-redis-py . -f Dockerfile -p kai-local
