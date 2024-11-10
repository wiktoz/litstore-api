CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'lang_type') THEN
      CREATE TYPE lang_type AS ENUM ('pl', 'en', 'fr', 'de');
   END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'select_type') THEN
        CREATE TYPE select_type AS ENUM ('button', 'select');
    END IF;
END $$;


DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'unit_type') THEN
        CREATE TYPE unit_type AS ENUM ('pc.', 'l', 'kg', 'set');
    END IF;
END $$;