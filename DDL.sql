CREATE TABLE mst_users(
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
    tx_count integer,
    CONSTRAINT unique_username UNIQUE(username)
);

CREATE TABLE mst_bank_account(
    account_id serial NOT NULL PRIMARY KEY,
    user_id integer NOT NULL,
    bank_name character varying(50)  NOT NULL,
    account_number character varying(20)  NOT NULL,
    account_holder_name character varying(50) NOT NULL,
    CONSTRAINT bankAcc_userID_fkey FOREIGN KEY(user_id)
    REFERENCES mst_users(user_id)
);

CREATE TABLE mst_card(
    card_id serial NOT NULL PRIMARY KEY,
    user_id integer NOT NULL,
    card_type character varying(10) NOT NULL,
    card_number character varying(20) NOT NULL,
    expiration_date date NOT NULL,
    cvv character varying(5) NOT NULL,
    CONSTRAINT card_userID_fkey FOREIGN KEY(user_id)
    REFERENCES mst_users(user_id)
);

<<<<<<< HEAD
CREATE TABLE tx_transaction
(
    tx_id serial NOT NULL PRIMARY KEY,
    user_id integer,
    amount numeric(10,2),
    type character varying(20),
    created_at timestamp without time zone DEFAULT 'now()',
    CONSTRAINT transaction_userID_fkey FOREIGN KEY(user_id)
    REFERENCES mst_users(user_id)
);

CREATE TABLE mst_photo_url
(

    photo_id serial NOT NULL PRIMARY KEY,
    url_photo character varying,
    user_id integer,
    CONSTRAINT photo_userID_fkey FOREIGN KEY(user_id)
    REFERENCES mst_users(user_id)
);

CREATE TABLE mst_point_exchange(
    pe_id serial NOT NULL PRIMARY KEY,
    reward character varying(100) NOT NULL,
    price_reward integer NOT NULL
);

CREATE TABLE tx_transaction(
  tx_id serial NOT NULL PRIMARY KEY,
  amount integer,
  transaction_type character varying(20),
  transaction_date date,
  sender_id integer,
  recipient_id integer,
  bank_account_id integer,
  card_id integer,
  pe_id integer,
  point integer,
  CONSTRAINT tx_bankAcc_fkey FOREIGN KEY(bank_account_id)
  REFERENCES mst_bank_account(account_id),
  CONSTRAINT tx_cardID_fkey FOREIGN KEY(card_id)
  REFERENCES mst_card(card_id),
  CONSTRAINT tx_senderID_fkey FOREIGN KEY(sender_id)
  REFERENCES mst_users(user_id),
  CONSTRAINT tx_recipientID_fkey FOREIGN KEY(recipient_id)
  REFERENCES mst_users(user_id)
);
