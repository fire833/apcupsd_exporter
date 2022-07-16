
VERSION		=		v0.3.0

.PHONY: binary
binary:
	go build cmd/apcupsd_exporter/main.go -o apcupsd_exporter

.PHONY: buildah
buildah:
	buildah bud -t apcupsd_exporter:${VERSION}
