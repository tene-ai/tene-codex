#!/usr/bin/env python3
"""Generate a deterministic minimal SPDX 2.3 SBOM for a staged plugin."""
import hashlib, json, pathlib, sys
stage=pathlib.Path(sys.argv[1]);version=sys.argv[2]
files=[]
for path in sorted(p for p in stage.rglob("*") if p.is_file() and p.name not in {"checksums.txt","sbom.spdx.json"}):
    digest=hashlib.sha256(path.read_bytes()).hexdigest();rel=path.relative_to(stage).as_posix()
    files.append({"SPDXID":"SPDXRef-File-"+hashlib.sha256(rel.encode()).hexdigest()[:16],"fileName":"./"+rel,"checksums":[{"algorithm":"SHA256","checksumValue":digest}]})
doc={"spdxVersion":"SPDX-2.3","dataLicense":"CC0-1.0","SPDXID":"SPDXRef-DOCUMENT","name":"tene-codex-"+version,"documentNamespace":"https://github.com/tene-ai/tene-codex/sbom/"+version,"creationInfo":{"created":"2026-08-20T00:00:00Z","creators":["Organization: tene-ai","Tool: tene-codex-generate-sbom"]},"packages":[{"name":"tene-codex","SPDXID":"SPDXRef-Package","versionInfo":version,"downloadLocation":"https://github.com/tene-ai/tene-codex","filesAnalyzed":True,"licenseConcluded":"Apache-2.0","licenseDeclared":"Apache-2.0","copyrightText":"Copyright 2026 Kay Kim (kay@agentkay.it)"}],"files":files,"relationships":[{"spdxElementId":"SPDXRef-DOCUMENT","relationshipType":"DESCRIBES","relatedSpdxElement":"SPDXRef-Package"}]+[{"spdxElementId":"SPDXRef-Package","relationshipType":"CONTAINS","relatedSpdxElement":f["SPDXID"]} for f in files]}
(stage/"sbom.spdx.json").write_text(json.dumps(doc,indent=2)+"\n")
