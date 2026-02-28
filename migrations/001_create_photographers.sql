CREATE TABLE photographers (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    specialty  TEXT NOT NULL,
    city       TEXT NOT NULL,
    bio        TEXT NOT NULL DEFAULT '',
    email      TEXT NOT NULL DEFAULT '',
    whatsapp   TEXT NOT NULL DEFAULT '',
    website    TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO photographers (name, specialty, city, email, whatsapp) VALUES
    ('Amir Hamzah', 'Wedding', 'Kuala Lumpur', 'amir@email.com', '+60123456789');

INSERT INTO photographers (name, specialty, city, email, website) VALUES
    ('Sara Lim', 'Street', 'Penang', 'sara@email.com', 'saralim.com');

INSERT INTO photographers (name, specialty, city, whatsapp) VALUES
    ('Razif Osman', 'Commercial', 'Johor Bahru', '+60198765432');