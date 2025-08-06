.PHONY: ca.gen.%
ca.gen.%:
	$(eval CA := $(word 1,$(sub ., ,$*)))
	@echo "=============> Generating CA files for $(CA)"
	@${ROOT_DIR}/scripts/gencrets.sh generate-iam-cert ${OUTPUT_DIR}/cert ${CA}

.PHONY: ca.gen
ca.gen: $(addprefix ca.gen., ${CERTIFICATES})
