function renderCVEs(cves) {
    if (!cves || cves.length === 0) {
        return '<p>No CVEs found for this image.</p>';
    }
    return `
        <ul class="detail-list">
            ${cves.map(cve => `
                <li>
                    <strong class="severity-${cve.severity}">${cve.cve_id}</strong>
                    <span>(${cve.severity}) - ${cve.fix_available ? 'Fix Available' : 'No Fix'}</span>
                </li>
            `).join('')}
        </ul>
    `;
}

function renderSBOMs(sboms) {
    if (!sboms || sboms.length === 0) {
        return '<p>No SBOMs found.</p>';
    }
    return `
        <ul class="detail-list">
            ${sboms.map(sbom => `
                <li>
                    <strong>${sbom.format}</strong>
                    <a href="${sbom.uri}" target="_blank" rel="noopener noreferrer">${sbom.uri}</a>
                </li>
            `).join('')}
        </ul>
    `;
}

function renderVerification(verification) {
    return `
        <ul class="detail-list">
            <li><strong>Signature Key ID</strong> <code>${verification.signature.keyId}</code></li>
            <li><strong>Rekor Entry</strong> <a href="${verification.rekorEntry}" target="_blank" rel="noopener noreferrer">View Entry</a></li>
        </ul>
        <h4>Attestations</h4>
        <ul class="detail-list">
            ${verification.attestations.map(att => `
                <li>
                    <strong>${att.type}</strong>
                    <a href="${att.uri}" target="_blank" rel="noopener noreferrer">${att.uri}</a>
                </li>
            `).join('')}
        </ul>
    `;
}

export function renderImageDetail({ image, sboms, cves, verification }) {
    return `
        <a href="#" class="back-link">&larr; Back to Image List</a>
        <div class="card">
            <h2>Image Details: <code style="font-size: 1.1rem;">${image.digest.substring(0, 20)}...</code></h2>
            <ul class="detail-list">
                <li><strong>ID</strong> ${image.id}</li>
                <li><strong>Tenant</strong> ${image.tenant_namespace}</li>
                <li><strong>Digest</strong> <code>${image.digest}</code></li>
                <li><strong>Tags</strong> ${image.tags ? image.tags.map(tag => `<span class="tag">${tag}</span>`).join(' ') : 'None'}</li>
                <li><strong>SLSA Level</strong> ${image.slsa_level}</li>
                <li><strong>Created At</strong> ${new Date(image.created_at).toLocaleString()}</li>
            </ul>
        </div>

        <div class="grid">
            <div class="card">
                <h2>CVE Findings</h2>
                ${renderCVEs(cves)}
            </div>
            <div class="card">
                <h2>SBOMs (Software Bill of Materials)</h2>
                ${renderSBOMs(sboms)}
            </div>
        </div>
        <div class="card">
            <h2>Verification & Attestations</h2>
            ${renderVerification(verification)}
        </div>
    `;
}
```
