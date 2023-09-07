#!/bin/sh

minikube image build -t konstellation/demo-task-py . -f Dockerfile -p kai-local
