-- Create invites table
CREATE TABLE invites (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('guest', 'customer', 'driver', 'admin')),
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    invited_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
