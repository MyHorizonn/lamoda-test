SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;


CREATE TABLE public.goods (
    uuid uuid NOT NULL,
    name character varying(48),
    size character varying(20)
);


ALTER TABLE public.goods OWNER TO postgres;


CREATE TABLE public.goods_in_store (
    id integer NOT NULL,
    store_id integer NOT NULL,
    goods_uuid uuid NOT NULL,
    amount bigint DEFAULT 0 NOT NULL,
    reserved bigint DEFAULT 0 NOT NULL
);

ALTER TABLE public.goods_in_store OWNER TO postgres;

CREATE SEQUENCE public.goods_in_store_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.goods_in_store_id_seq OWNER TO postgres;

ALTER SEQUENCE public.goods_in_store_id_seq OWNED BY public.goods_in_store.id;

CREATE TABLE public.store (
    id integer NOT NULL,
    name character varying(48),
    accessibility boolean DEFAULT true NOT NULL
);


ALTER TABLE public.store OWNER TO postgres;


CREATE SEQUENCE public.storage_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.storage_id_seq OWNER TO postgres;


ALTER SEQUENCE public.storage_id_seq OWNED BY public.store.id;

ALTER TABLE ONLY public.goods_in_store ALTER COLUMN id SET DEFAULT nextval('public.goods_in_store_id_seq'::regclass);

ALTER TABLE ONLY public.store ALTER COLUMN id SET DEFAULT nextval('public.storage_id_seq'::regclass);


ALTER TABLE ONLY public.goods_in_store
    ADD CONSTRAINT goods_in_store_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.goods
    ADD CONSTRAINT goods_pkey PRIMARY KEY (uuid);


ALTER TABLE ONLY public.store
    ADD CONSTRAINT storage_pkey PRIMARY KEY (id);


ALTER TABLE ONLY public.goods_in_store
    ADD CONSTRAINT goods_fk FOREIGN KEY (goods_uuid) REFERENCES public.goods(uuid);

ALTER TABLE ONLY public.goods_in_store
    ADD CONSTRAINT store_fk FOREIGN KEY (store_id) REFERENCES public.store(id);
