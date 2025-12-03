CREATE TABLE IF NOT EXISTS tasks (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL, -- เจ้าของ task
    title       TEXT        NOT NULL,
    description TEXT        NULL,
    status      TEXT        NOT NULL DEFAULT 'todo', -- todo / in_progress / done
    due_date    TIMESTAMPTZ NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_user_id_status ON tasks(user_id, status);
