CREATE TABLE public.skills (
    id SERIAL PRIMARY KEY,
    name character varying(255) NOT NULL,
    level integer NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE TABLE public.contacts (
    id SERIAL PRIMARY KEY,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    full_name character varying(512) NOT NULL,
    address character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    mobile character varying(255) NOT NULL,
    password character varying(512) NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE TABLE public.contacts_skills (
    id SERIAL PRIMARY KEY,
    contact_id integer,
    skill_id integer
);

INSERT INTO contacts(first_name,last_name,full_name,email,address,mobile,password,created_at,updated_at) VALUES ('mei-chih','chang','mei-chih chang', 'mei@psi.ch', 'Imbisbühlstrasse 126, 8049 Zürich', '0765933632', '$2a$14$wVsaPvJnJJsomWArouWCtusem6S/.Gauq/GjOIEHpyh2DAMmso1wy', '2023-04-23 15:18:15', '2023-04-23 15:19:15');
