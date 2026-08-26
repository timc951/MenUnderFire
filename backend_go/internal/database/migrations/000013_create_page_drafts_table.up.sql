-- Page drafts for preview functionality
CREATE TABLE page_drafts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id UUID REFERENCES site_pages(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    created_by_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Index for cleanup of expired drafts
CREATE INDEX idx_page_drafts_expires_at ON page_drafts(expires_at);

-- Index for finding drafts by page
CREATE INDEX idx_page_drafts_page_id ON page_drafts(page_id);
