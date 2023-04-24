CREATE TABLE mst_users
(
    user_id serial NOT NULL PRIMARY KEY,
    name character varying(50) NOT NULL,
    email character varying(50) NOT NULL,
    password character varying NOT NULL,
    phone_number character varying(20) NOT NULL,
    address character varying(100) NOT NULL,
    balance integer NOT NULL,
    username character varying NOT NULL,
    point integer,
    role character varying,
    CONSTRAINT unique_username UNIQUE(username)
);

CREATE TABLE mst_bank_account
(
    account_id serial NOT NULL PRIMARY KEY,
    user_id integer NOT NULL,
    bank_name character varying(50)  NOT NULL,
    account_number character varying(20)  NOT NULL,
    account_holder_name character varying(50) NOT NULL,
    CONSTRAINT bankAcc_userID_fkey FOREIGN KEY(user_id)
    REFERENCES mst_users(user_id)
);

CREATE TABLE mst_card
(
    card_id serial NOT NULL PRIMARY KEY,
    user_id integer NOT NULL,
    card_type character varying(10) NOT NULL,
    card_number character varying(20) NOT NULL,
    expiration_date date NOT NULL,
    cvv character varying(5) NOT NULL,
    CONSTRAINT card_userID_fkey FOREIGN KEY(user_id)
    REFERENCES mst_users(user_id)
);

-- Table: public.tx_transaction

-- DROP TABLE IF EXISTS public.tx_transaction;

CREATE TABLE IF NOT EXISTS public.tx_transaction
(
    tx_id integer NOT NULL DEFAULT 'nextval('transaction_id_seq'::regclass)',
    amount integer,
    transaction_type character varying(20) COLLATE pg_catalog."default",
    "timestamp" timestamp without time zone DEFAULT 'now()',
    sender_id integer,
    recipient_id integer,
    bank_account_id integer,
    card_id integer,
    pe_id integer,
    point integer,
    CONSTRAINT transaction_pkey PRIMARY KEY (tx_id),
    CONSTRAINT fk_transaction_bank_acc FOREIGN KEY (bank_account_id)
        REFERENCES public.mst_bank_account (account_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fk_transaction_card FOREIGN KEY (card_id)
        REFERENCES public.mst_card (card_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fk_transaction_recipient_user_id FOREIGN KEY (recipient_id)
        REFERENCES public.mst_users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT fk_user_id FOREIGN KEY (sender_id)
        REFERENCES public.mst_users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tx_transaction
    OWNER to postgres;

CREATE TABLE mst_photo_url
(
    photo_id serial NOT NULL PRIMARY KEY,
    url_photo character varying,
    user_id integer,
    CONSTRAINT photo_userID_fkey FOREIGN KEY(user_id)
    REFERENCES mst_users(user_id)
);

CREATE TABLE mst_point_exchange
(
    pe_id serial NOT NULL PRIMARY KEY,
    reward character varying(100) NOT NULL,
    price_reward integer NOT NULL
);
