CREATE TABLE person
(
    id uuid NOT NULL
    name varchar(255) NOT NULL
    age int NOT NULL
    is_healthy boolean NOT NULL
) 

CREATE TABLE user
(   
    id uuid NOT NULL
    login varchar(255) NOT NULL
    password varchar(255) NOT NULL
    role varchar(255) NOT NULL
    refresh_token varchar(255)
)