CREATE TABLE IF NOT EXISTS posts (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    comments_allowed BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS comments (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    post_id BIGINT NOT NULL,
    reply_to_comment_id BIGINT,
    user_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_post FOREIGN KEY (post_id)
        REFERENCES posts(id) ON DELETE CASCADE,

    CONSTRAINT fk_reply_to_comment FOREIGN KEY (reply_to_comment_id)
        REFERENCES comments(id) ON DELETE CASCADE,

    CONSTRAINT comment_text_length_check
        CHECK (LENGTH(text) <= 2000)
);

CREATE INDEX IF NOT EXISTS idx_posts_created_at_id
ON posts(created_at, id);

CREATE INDEX IF NOT EXISTS idx_comments_post_reply_created_at_id
ON comments(post_id, reply_to_comment_id, created_at, id);