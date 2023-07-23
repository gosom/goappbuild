CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    name TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id),
    UNIQUE(user_id, name)
);

CREATE TABLE collections (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   created_at TIMESTAMP WITH TIME ZONE NOT NULL,
   updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
   name TEXT NOT NULL,
   project_id UUID NOT NULL REFERENCES projects (id),
   attributes JSONB NOT NULL CHECK (jsonb_typeof(attributes) = 'object'),
   UNIQUE(project_id, name)
);
