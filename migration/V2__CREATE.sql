CREATE TABLE IF NOT EXISTS goschema.person
(
    id uuid NOT null,
    "name" varchar(255) NOT null,
    age int NOT null,
    is_healthy boolean NOT NULL
);

CREATE TABLE IF NOT exists goschema.user
(   
    id uuid NOT null,
    login varchar(255) NOT null,
    password varchar(255) NOT null,
    role varchar(255) NOT null,
    refresh_token varchar(255)
);