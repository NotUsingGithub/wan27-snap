# wan27-snap — Linux Snap Package for Wan 2.7 AI CLI

Official Linux Snap package for the **[Wan 2.7 AI](https://wan2-7.io/)** CLI (`wan27`). Generate 1080p/2K text-to-video, image-to-video, and AI images, or estimate credit costs directly from Ubuntu, Debian, Fedora, openSUSE, and Arch Linux terminals.

## Official Documentation & Links

- **Wan 2.7 Website & Studio**: [https://wan2-7.io/](https://wan2-7.io/)
- **Wan 2.7 MCP Server**: [https://wan2-7.io/mcp/](https://wan2-7.io/mcp/)
- **Wan 2.7 MCP Setup Guide**: [https://wan2-7.io/mcp/guide/](https://wan2-7.io/mcp/guide/)
- **Wan 2.7 CLI Reference**: [https://wan2-7.io/cli/](https://wan2-7.io/cli/)
- **Wan 2.7 Pricing & Credits**: [https://wan2-7.io/pricing/](https://wan2-7.io/pricing/)

## Quick Start

```bash
sudo snap install wan27
export WAN27_API_TOKEN="wan27_your_token_here"   # From https://wan2-7.io/mcp/guide/

# Check remaining credits
wan27 -credits

# Estimate credit cost before running
wan27 -estimate -model wan2.7-t2v -resolution 1080p "A cinematic drone shot over neon Tokyo"

# Generate video or image
wan27 -model wan2.7-t2v -resolution 1080p "A cinematic drone shot over neon Tokyo"
wan27 -model wan2.7-image-pro -aspect-ratio 1:1 "Studio product photography of a mechanical watch"
```

## Building & Publishing

`./build.sh` cross-compiles static `amd64` and `arm64` Linux binaries and packs `build/wan27_0.1.0_amd64.snap` and `build/wan27_0.1.0_arm64.snap` using `mksquashfs` (no LXD or Multipass VM required). GitHub Actions automatically builds and attaches both `.snap` packages to each GitHub Release.

## License

MIT © [Wan27](https://wan2-7.io/)
