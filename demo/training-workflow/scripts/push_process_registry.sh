#!/bin/bash

kli process-registry register trigger github-trigger-mock --dockerfile ./src/github-trigger-mock/Dockerfile --product demo --src ./src/github-trigger-mock --version v1
kli process-registry register task splitter --dockerfile ./src/splitter/Dockerfile --product demo --src ./src/splitter --version v1
kli process-registry register task training-py --dockerfile ./src/training-py/Dockerfile --product demo --src ./src/training-py --version v1
kli process-registry register task training-go --dockerfile ./src/training-go/Dockerfile --product demo --src ./src/training-go --version v1
kli process-registry register task validation --dockerfile ./src/validation/Dockerfile --product demo --src ./src/validation --version v1
kli process-registry register exit exit --dockerfile ./src/exit/Dockerfile --product demo --src ./src/exit --version v1

echo "Done"
