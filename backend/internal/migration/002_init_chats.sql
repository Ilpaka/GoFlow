-- Chats + members. last_message_id / last_read_message_id FKs added in 003 after messages exist.

CREATE TABLE chats (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    type text NOT NULL,
    title varchar(255),
    avatar_url text,
    created_by uuid REFERENCES users (id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    last_message_id uuid,
    last_message_at timestamptz,
    is_deleted boolean NOT NULL DEFAULT false,
    direct_key varchar(64),
    CONSTRAINT ck_chats_type CHECK (type IN ('direct', 'group')),
    CONSTRAINT ck_chats_direct_key CHECK (
        (type = 'direct' AND direct_key IS NOT NULL)
        OR (type = 'group' AND direct_key IS NULL)
    )
);

CREATE UNIQUE INDEX uq_chats_direct_key ON chats (direct_key)
WHERE type = 'direct' AND direct_key IS NOT NULL AND is_deleted = false;

CREATE INDEX idx_chats_created_by ON chats (created_by) WHERE created_by IS NOT NULL;
CREATE INDEX idx_chats_last_message_at ON chats (last_message_at DESC) WHERE is_deleted = false;
CREATE INDEX idx_chats_type ON chats (type) WHERE is_deleted = false;
CREATE INDEX idx_chats_is_deleted ON chats (is_deleted) WHERE is_deleted = true;

CREATE TABLE chat_members (
    chat_id uuid NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role text NOT NULL,
    joined_at timestamptz NOT NULL DEFAULT now(),
    last_read_message_id uuid,
    last_read_at timestamptz,
    is_muted boolean NOT NULL DEFAULT false,
    is_archived boolean NOT NULL DEFAULT false,
    is_pinned boolean NOT NULL DEFAULT false,
    PRIMARY KEY (chat_id, user_id),
    CONSTRAINT ck_chat_members_role CHECK (role IN ('owner', 'admin', 'member'))
);

CREATE INDEX idx_chat_members_user_id ON chat_members (user_id);
CREATE INDEX idx_chat_members_user_archived ON chat_members (user_id, is_archived);
CREATE INDEX idx_chat_members_user_pinned ON chat_members (user_id, is_pinned) WHERE is_pinned = true;
