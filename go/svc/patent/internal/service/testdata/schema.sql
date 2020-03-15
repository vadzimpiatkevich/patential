CREATE TABLE patents (
  id CHAR(36) NOT NULL CONSTRAINT patents_pkey PRIMARY KEY,
  application_number CHAR(36) NOT NULL,
  application_kind CHAR(36) NOT NULL,
  grant_date DATE NOT NULL,

  UNIQUE(application_number)
);
