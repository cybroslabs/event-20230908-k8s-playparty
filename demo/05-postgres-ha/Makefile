.PHONY: install-postgres
install-postgres:
	helm upgrade --install postgres oci://registry-1.docker.io/bitnamicharts/postgresql-ha \
		--version 11.7.6 \
		--set postgresql.password=admin \
		--set postgresql.repmgrPassword=admin \
		--set pgpool.adminPassword=admin \
		-f ./values.yaml