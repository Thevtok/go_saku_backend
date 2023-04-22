
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
);


CREATE TABLE IF NOT EXISTS public.mst_users
(
    user_id serial NOT NULL PRIMARY KEY,
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
);

CREATE TABLE IF NOT EXISTS public.mst_card
(
    card_id serial NOT NULL PRIMARY KEY,
    username character varying(100) NOT NULL,
    card_type character varying(10) COLLATE pg_catalog."default" NOT NULL,
    card_number character varying(20) COLLATE pg_catalog."default" NOT NULL,
    expiration_date date NOT NULL,
    cvv character varying(5) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT card_pkey PRIMARY KEY (card_id),
    CONSTRAINT card_username_fkey FOREIGN KEY (username)
        REFERENCES public.mst_users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);


CREATE TABLE IF NOT EXISTS public.tx_transaction
(
    tx_id serial NOT NULL PRIMARY KEY,
    username character varying COLLATE pg_catalog."default",
    amount numeric(10,2),
    type character varying(20) COLLATE pg_catalog."default",
    created_at timestamp without time zone DEFAULT 'now()',
    CONSTRAINT transaction_pkey PRIMARY KEY (tx_id),
    CONSTRAINT transaction_username_fkey FOREIGN KEY (username)
        REFERENCES public.mst_users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS public.mst_photo_url
(
    photo_id serial NOT NULL PRIMARY KEY,
    url_photo character varying COLLATE pg_catalog."default",
    username character varying COLLATE pg_catalog."default",
    CONSTRAINT photo_pkey PRIMARY KEY (photo_id),
    CONSTRAINT photo_username_fkey FOREIGN KEY (username)
        REFERENCES public.mst_users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS public.mst_point_exchange
(
    pe_id serial NOT NULL PRIMARY KEY,
    reward character varying(100) NOT NULL,
    price_reward integer NOT NULL,
);
