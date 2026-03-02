-- Messages: soft delete via deleted_at. Links chats, sender, optional reply thread.

CREATE TABLE messages (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id uuid NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    sender_id uuid NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    type text NOT NULL,
    text text,
    reply_to_id uuid REFERENCES messages (id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,
    CONSTRAINT ck_messages_type CHECK (type IN ('text', 'image', 'file', 'system'))
);

CREATE INDEX idx_messages_chat_created ON messages (chat_id, created_at DESC)
WHERE deleted_at IS NULL;

CREATE INDEX idx_messages_sender_id ON messages (sender_id);
CREATE INDEX idx_messages_reply_to_id ON messages (reply_to_id) WHERE reply_to_id IS NOT NULL;
CREATE INDEX idx_messages_deleted_at ON messages (deleted_at) WHERE deleted_at IS NOT NULL;

ALTER TABLE chats
    ADD CONSTRAINT fk_chats_last_message
    FOREIGN KEY (last_message_id) REFERENCES messages (id) ON DELETE SET NULL;

ALTER TABLE chat_members
    ADD CONSTRAINT fk_chat_members_last_read_message
    FOREIGN KEY (last_read_message_id) REFERENCES messages (id) ON DELETE SET NULL;
