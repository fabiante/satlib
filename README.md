# Satlib

This is a library for implementing 'serverless-like' sidecar services to existing,
monolithic application's.

The target audience for this library are developers who must integrate
CLIs to an existing application but don't want or can't add those CLIs
to their existing container images.

A reasons for this may be that your software already has a
large container image and you don't want to bloat it further. By using
sidecar containers you keep your original image at the same size and you
you can more flexibly upgrade the sidecar or you main application. Think upgrading
your base Linux distro and then having _lots_ of unexpected package version changes
which may or may not lead to issues.

## Basic Idea

- You have some existing application which requires a CLI like `tesseract`
- You create a new Go project which imports this library and wraps the CLI
- You then create a new `Dockerfile` which essentially contains only the `tesseract` CLI and your compiled Go program
- This image can be integrated into your existing application's stack.
