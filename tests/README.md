# Examples Acceptance Tests

## Overview

This folder contains the acceptance tests for the examples that you can run as part of the Kyma testing process. The tests are written in Go. Run them as standard Go tests. Each example has a separate folder such as `http-db-example`.

## Usage

This section provides information on how to build and version the Docker image, as well as how to configure Kyma.

### Configure Kyma

After building and pushing the Docker image, set the proper tag `acceptanceTest.imageTag` in the `resources/core/values.yaml` file.
