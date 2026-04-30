# checker-ping

ICMP ping checker for [happyDomain](https://www.happydomain.org/).

Checks reachability, round-trip time, and packet loss for IP addresses associated with a domain's services.

## Usage

### Standalone HTTP server

```bash
# Build and run
make
./checker-ping -listen :8080

# With privileged ICMP (requires CAP_NET_RAW or root)
./checker-ping -listen :8080 -privileged
```

The server exposes:

- `GET /health`, health check
- `POST /collect`, collect ping observations (happyDomain external checker protocol)

### Docker

```bash
make docker
docker run -p 8080:8080 happydomain/checker-ping

# With privileged ICMP
docker run --cap-add NET_RAW -p 8080:8080 happydomain/checker-ping -privileged
```

### happyDomain plugin

```bash
make plugin
# produces checker-ping.so, loadable by happyDomain as a Go plugin
```

The plugin exposes a `NewCheckerPlugin` symbol returning the checker
definition and observation provider, which happyDomain registers in its
global registries at load time.

### Versioning

The binary, plugin, and Docker image embed a version string overridable
at build time:

```bash
make CHECKER_VERSION=1.2.3
make plugin CHECKER_VERSION=1.2.3
make docker CHECKER_VERSION=1.2.3
```

### happyDomain remote endpoint

Set the `endpoint` admin option for the ping checker to the URL of the running checker-ping server (e.g., `http://checker-ping:8080`). happyDomain will delegate observation collection to this endpoint.

### Deployment

The `/collect` endpoint has no built-in authentication and will issue
ICMP probes to whatever hostnames or IP addresses the caller supplies,
including private (RFC 1918), loopback, link-local, and cloud metadata
addresses (e.g. `169.254.169.254`). It is meant to run on a trusted
network, reachable only by the happyDomain instance that drives it.
Restrict access via a reverse proxy with authentication, a network ACL,
or by binding the listener to a private interface; do not expose it
directly to the public internet.

## Rules

| Code                  | Description                                                                                       | Severity            |
|-----------------------|---------------------------------------------------------------------------------------------------|---------------------|
| `ping.reachable`      | Verifies every target replied to at least one ICMP probe.                                         | CRITICAL            |
| `ping.rtt`            | Flags targets whose average round-trip time crosses the warning or critical threshold.            | CRITICAL            |
| `ping.ipv6_reachable` | Verifies that at least one IPv6 target replied to an ICMP probe.                                  | WARNING             |
| `ping.packet_loss`    | Flags targets whose packet-loss ratio crosses the warning or critical threshold.                  | CRITICAL            |

## License & licensing roadmap

This project is currently licensed under the **GNU Affero General Public
License v3.0** (see `LICENSE`), because it still imports
`happydns.ServiceMessage` and `abstract.Server` from the happyDomain
server module (`git.happydns.org/happyDomain/model` and
`git.happydns.org/happyDomain/services/abstract`), which are themselves
distributed under AGPL-3.0 and a commercial license.

The core checker types (`CheckerOptions`, `CheckerDefinition`,
`ObservationProvider`, `CheckRule`, …) have already been migrated to
[`checker-sdk-go`](https://git.happydns.org/checker-sdk-go); only the
service-message types remain on the AGPL side.

**Planned relicensing:** as soon as the remaining `ServiceMessage` /
`abstract.Server` dependency has been removed (moved into a dedicated
permissively licensed module), this project will be relicensed under the
**MIT License**, in line with the rest of the happyDomain checker
ecosystem (see `checker-dummy` for the target shape).

**Contributors notice:** by submitting a contribution to this repository,
you accept that your contribution will be relicensed from AGPL-3.0 to MIT
at the time of the relicensing described above. If you do not agree with
this, please do not submit contributions until the relicensing has taken
place.

The third-party Apache-2.0 attributions for `checker-sdk-go` are recorded
in `NOTICE` and must accompany any binary or source redistribution of this
project.
