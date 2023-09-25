#!/bin/sh

minikube image build -t konstellation/demo-cronjob-trigger-py . -f Dockerfile -p kai-local
