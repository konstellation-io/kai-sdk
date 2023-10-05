#!/bin/sh

minikube image build -t konstellation/demo-exit-py . -f Dockerfile -p kai-local
