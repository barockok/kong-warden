workdir:=$(PWD)
outdir:=$(workdir)/out

all: gsvc kplugin
kplugin:
	cd cmd/kong-plugin && go build -buildmode=default && mv kong-plugin $(outdir)/warden
gsvc:
	cd cmd/svc && go build && mv svc $(outdir)/
clean:
	rm out/*