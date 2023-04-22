-- Table: public.mst_bank_account

-- DROP TABLE IF EXISTS public.mst_bank_account;

CREATE TABLE IF NOT EXISTS public.mst_bank_account
(
    account_id serial PRIMARY KEY NOT NULL,
    bank_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    account_number character varying(20) COLLATE pg_catalog."default" NOT NULL,
    account_holder_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    username character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT bank_account_pkey PRIMARY KEY (account_id),
    CONSTRAINT fk_mst_bank_mst_users FOREIGN KEY (username)
        REFERENCES public.mst_users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.mst_bank_account
    OWNER to postgres;

    -- Table: public.mst_users

-- DROP TABLE IF EXISTS public.mst_users;

CREATE TABLE IF NOT EXISTS public.mst_users
(
    user_id integer NOT NULL DEFAULT 'nextval('users_user_id_seq'::regclass)',
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    email character varying(50) COLLATE pg_catalog."default" NOT NULL,
    password character varying COLLATE pg_catalog."default" NOT NULL,
    phone_number character varying(20) COLLATE pg_catalog."default" NOT NULL,
    address character varying(100) COLLATE pg_catalog."default" NOT NULL,
    balance integer NOT NULL,
    username character varying COLLATE pg_catalog."default" NOT NULL,
    point integer,
    role character varying COLLATE pg_catalog."default",
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT unique_username UNIQUE (username)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.mst_users
    OWNER to postgres;

    -- Table: public.mst_card

-- DROP TABLE IF EXISTS public.mst_card;

CREATE TABLE IF NOT EXISTS public.mst_card
(
    card_id integer NOT NULL DEFAULT 'nextval('card_card_id_seq'::regclass)',
    user_id integer NOT NULL,
    card_type character varying(10) COLLATE pg_catalog."default" NOT NULL,
    card_number character varying(20) COLLATE pg_catalog."default" NOT NULL,
    expiration_date date NOT NULL,
    cvv character varying(5) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT card_pkey PRIMARY KEY (card_id),
    CONSTRAINT card_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.mst_users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.mst_card
    OWNER to postgres;

    -- Table: public.tx_transaction

-- DROP TABLE IF EXISTS public.tx_transaction;

CREATE TABLE IF NOT EXISTS public.tx_transaction
(
    id integer NOT NULL DEFAULT 'nextval('transaction_id_seq'::regclass)',
    username character varying COLLATE pg_catalog."default",
    amount numeric(10,2),
    type character varying(20) COLLATE pg_catalog."default",
    created_at timestamp without time zone DEFAULT 'now()',
    CONSTRAINT transaction_pkey PRIMARY KEY (id),
    CONSTRAINT transaction_username_fkey FOREIGN KEY (username)
        REFERENCES public.mst_users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tx_transaction
    OWNER to postgres;