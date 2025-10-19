export function renderImageTable(images) {
    if (!images || images.length === 0) {
        return '<div class="card"><p>No images found.</p></div>';
    }

    const rows = images.map(image => `
        <tr>
            <td>${image.id}</td>
            <td>${image.tenant_namespace}</td>
            <td><code title="${image.digest}">${image.digest.substring(0, 20)}...</code></td>
            <td>${image.tags ? image.tags.map(tag => `<span class="tag">${tag}</span>`).join(' ') : ''}</td>
            <td>SLSA ${image.slsa_level}</td>
            <td><a href="#/images/${image.id}" class="btn">View Details</a></td>
        </tr>
    `).join('');

    return `
        <div class="card">
            <h2>Image Explorer</h2>
            <table class="image-table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Tenant</th>
                        <th>Digest</th>
                        <th>Tags</th>
                        <th>SLSA Level</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    ${rows}
                </tbody>
            </table>
        </div>
    `;
}

