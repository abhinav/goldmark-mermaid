# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0] - 2023-03-24
### Changed
- ClientRenderer: Use `<pre>` instead of `<div>` for diagram containers.

### Added
- Support changing the container tag with the `ContainerTag` option.
  This option is available on ClientRenderer, ServerRenderer, and Extender.

  [0.4.0]: https://github.com/abhinav/goldmark-mermaid/releases/tag/v0.4.0

## [0.3.0] - 2022-12-19
### Changed
- Change the module path to `go.abhg.dev/goldmark/mermaid`.

### Removed
- Delete previously deprecated Renderer type.
  Please use the ClientRenderer instead.

  [0.3.0]: https://github.com/abhinav/goldmark-mermaid/releases/tag/v0.3.0

## [0.2.0] - 2022-11-04

### Added

- ServerRenderer with support for rendering Mermaid diagrams
  into inline SVGs server-side.
  This is picked automatically if an 'mmdc' executable is found on PATH.
- Support opting out of the MermaidJS `<script>` tag.
  To use, set `Extender.NoScript` or `Transformer.NoScript` to true.
  Use this if the page you're rendering into already includes the tag
  elsewhere.

### Changed

- Deprecate Renderer in favor of ClientRenderer.
  Rendere has been aliased to the new type
  so existing code should continue to work unchanged.

  [0.2.0]: https://github.com/abhinav/goldmark-mermaid/releases/tag/v0.2.0

## [0.1.1] - 2021-11-03

### Fixed

- Fix handling of multiple mermaid blocks.

  [0.1.1]: https://github.com/abhinav/goldmark-mermaid/releases/tag/v0.1.1

## [0.1.0] - 2021-04-12

- Initial release.

  [0.1.0]: https://github.com/abhinav/goldmark-mermaid/releases/tag/v0.1.0
