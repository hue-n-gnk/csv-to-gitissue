CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "nulab_backlog_issue" (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    backlog_id VARCHAR(50),
    git_id VARCHAR(50),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    PRIMARY KEY (id)
);
