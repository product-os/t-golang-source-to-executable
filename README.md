# Golang Build Transformer

A transformer that takes in [Golang source][golang-source] bundles and outputs [executables][executable]

[golang-source]: https://github.com/product-os/t-golang-source
[executable]: https://github.com/product-os/t-executable


## Notes for local development

you can set various env vars to change some behaviors:

- set `ARTIFACTPATH=assets` to change the output artifact path
- set `DEBUG=1` to get more verbose output
