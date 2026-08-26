-- Page hits tracking table
CREATE TABLE page_hits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path VARCHAR(500) NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    referrer TEXT,
    country VARCHAR(100),
    city VARCHAR(200),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_page_hits_path ON page_hits(path);
CREATE INDEX idx_page_hits_created_at ON page_hits(created_at DESC);
CREATE INDEX idx_page_hits_user_id ON page_hits(user_id);
