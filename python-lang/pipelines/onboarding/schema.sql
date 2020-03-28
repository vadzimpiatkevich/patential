CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TABLE public.patents (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  application_number VARCHAR (255) NOT NULL,
  application_kind VARCHAR (255) NOT NULL,
  grant_date DATE NOT NULL,

  UNIQUE (application_number)
);

CREATE UNIQUE INDEX index_patents_on_application_number
  ON public.patents
  USING btree(
    lower(application_number)
  );
