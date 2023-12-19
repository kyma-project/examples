## Overview

To contribute to this project, follow the rules from the general [CONTRIBUTING.md](https://github.com/kyma-project/community/blob/main/CONTRIBUTING.md) document in the `community` repository.

For additional guidelines, see other sections of this document.

## Naming guidelines

In addition to the general [naming guidelines](https://github.com/kyma-project/community/blob/main/docs/guidelines/technical-guidelines/01-naming.md), the `examples` repository has these conventions to follow:

- Folders:
  - All folder names are lowercase and connected with a dash (`-`) when they are comprised of two or more words. For example, `http-db-service`.
  - Folder names are as short as possible while being clear and descriptive.

- Files:
  - All markdown files are suffixed with the `.md` extension. Any `README.md` or `CONTRIBUTING.md` document names are in uppercase.
  - Source code files follow their own programming language naming conventions.
  - Configuration and Deployment files follow their own platform naming conventions, such as `Dockerfile` or `Gopkg.toml`. If there is no convention specified, filenames must have a valid extension, be lowercase, and be connected with a dash (`-`) when they are comprised of two or more words.

## Documentation types

### README.md

This file type contains information about what an example illustrates, and the instructions on how to run it. Each example in this repository requires a `README.md` document.

To create new `README.md` documents for new examples, use the [template](https://github.com/kyma-project/template-repository/blob/main/README.md) provided.
Do not change the names or the order of the main sections in the `README.md` documents. However, you can create subsections to adjust each `README.md` document to the example's specific requirements. See the example of a [README.md](http-db-service/README.md) document.

Find all `README.md` documents listed in the **List of examples** section in the main [README.md](README.md) document. When you create a new `README.md` document, add a new entry to the **List of examples** table. Follow these steps:

* Go to the **List of examples** section in the main [README.md](README.md) document located in the root of the `examples` project.
* Locate the new `README.md` document in the repository structure.
* Add a new entry to the relevant place on the table. Each entry consists of these three columns:
  * A relative link to the corresponding `README.md` document. The alias is the name of the project that the `README.md` document describes.
  * A short description of the example.
  * A list of the technologies used.

## Contribution rules

To contribute to the project, see the [CONTRIBUTING.md](https://github.com/kyma-project/community/blob/main/CONTRIBUTING.md) document in the `community` repository.

In addition to these general rules, every contributor to the `examples` repository must follow the rules described in the subsections.

### Structure

Structure examples according to the feature or concept they illustrate. Optionally, if a feature or concept has examples in several languages or technologies, add a second level of nesting to separate the examples. Do not use more than two levels of nesting as it hinders navigation and readability.

To add a new example, follow this layout structure:

```
kubectl-container-describe/        # No subdivision
serverless-lambda/                 # Subdivision based on language
    |---js/
    |---go/
broker-gateway/                    # Subdivision based on provider
    |--azure/
    |--aws/
    |--gcp/
storage/                           # Subdivision based on technology
    |--cassandra/
    |--mysql
    |--rethinkdb
```

> **NOTE:** This is an example layout and does not reflect the actual structure. It is based on the [Kubernetes examples structure](https://github.com/kubernetes/examples).

### Build

For all examples, the main `README.md` documents include the build scripts and the instructions on how to run them. Additionally, provide the instructions on how to clean up or delete the examples if applicable.

### Deploy

All examples are deployable in the cloud. Follow these rules to achieve this and to keep the complexity to a minimum:

- All examples have only one `Dockerfile` which generates a single Docker image. The Docker image follows the Docker image [naming conventions](https://github.com/kyma-project/community/blob/main/docs/guidelines/technical-guidelines/01-naming.md).
- Only example images that are completely built are accessible in the registry, so the examples can be deployed to the cloud without the user having to build and push a new image.
- All `README.md` documents in the examples provide the scripts to build and deploy images locally, and the instructions on how to run them.
- Provide Deployment configurations and descriptors in a user-friendly format such as a `yaml` file, and never in bundled formats such as Helm charts, unless the example itself illustrates the usage of bundled formats.
- Deployed resources of an example have a label defined in the `example: {Example Name}` format.
- Deployed resources of an example do not specify a Namespace. This way, you can deploy to the Namespace of your choice. As an exception, you can include a Namespace in the example's Deployment if it is necessary to preserve the integrity of the cluster or the user's Namespace. For example, your example deploys a custom Event Bus that would crash the cluster's eventing if it were deployed into the default or user's Namespace.
- Provide the instructions for each example on how to clean up all the deployed resources based on the label of the example.

### Continuous integration

A CI/CD configuration file accompanies each example to ensure its validity as part of the automation pipeline.

### Releases

Each example is independently releasable by creating a `tag` in the repository with the `{example-name}/{version}` format. This triggers a release for the given example from its CI/CD configuration. Push released example Docker images to the `example` folder in the Kyma repository, under `eu.gcr.io/kyma-project/example`.
