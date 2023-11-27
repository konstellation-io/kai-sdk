#!/bin/sh

minikube image build -t konstellation/demo-metrics-py . -f Dockerfile -p kai-local
