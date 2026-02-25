CREATE TABLE IF NOT EXISTS projects (
  id          BIGSERIAL PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  description TEXT,
  owner_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_members (
  project_id  BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role        VARCHAR(50) NOT NULL DEFAULT 'member',
  joined_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  PRIMARY KEY (project_id, user_id)
);

CREATE TYPE task_status AS ENUM ('todo', 'in_progress', 'done');

CREATE TABLE IF NOT EXISTS tasks (
  id           BIGSERIAL PRIMARY KEY,
  project_id   BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  title        VARCHAR(255) NOT NULL,
  description  TEXT,
  status       task_status NOT NULL DEFAULT 'todo',
  assignee_id  BIGINT REFERENCES users(id) ON DELETE SET NULL,
  due_date     TIMESTAMP WITH TIME ZONE,
  created_by   BIGINT NOT NULL REFERENCES users(id),
  created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_task_project_id ON tasks(project_id);
CREATE INDEX idx_task_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_project_members_user_id ON project_members(user_id);
