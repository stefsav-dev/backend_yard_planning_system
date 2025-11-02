--
-- PostgreSQL database dump
--

\restrict BsFZ5PoIeozMahLov2DIk8ubMv4ln2ez0bnq1BCX29seTtBYe4E9njnIsFMNfJE

-- Dumped from database version 18.0 (Ubuntu 18.0-1.pgdg25.10+3)
-- Dumped by pg_dump version 18.0 (Ubuntu 18.0-1.pgdg25.10+3)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: blocks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.blocks (
    id bigint NOT NULL,
    yard_id bigint NOT NULL,
    name text NOT NULL,
    max_slot bigint NOT NULL,
    max_row bigint NOT NULL,
    max_tier bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.blocks OWNER TO postgres;

--
-- Name: blocks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.blocks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.blocks_id_seq OWNER TO postgres;

--
-- Name: blocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.blocks_id_seq OWNED BY public.blocks.id;


--
-- Name: containers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.containers (
    id bigint NOT NULL,
    container_number text NOT NULL,
    block_id bigint NOT NULL,
    container_size bigint NOT NULL,
    container_height numeric NOT NULL,
    container_type text NOT NULL,
    slot bigint NOT NULL,
    "row" bigint NOT NULL,
    tier bigint NOT NULL,
    is_placed boolean DEFAULT true NOT NULL,
    placed_at timestamp with time zone,
    picked_up_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.containers OWNER TO postgres;

--
-- Name: containers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.containers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.containers_id_seq OWNER TO postgres;

--
-- Name: containers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.containers_id_seq OWNED BY public.containers.id;


--
-- Name: yard_plans; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.yard_plans (
    id bigint NOT NULL,
    block_id bigint NOT NULL,
    container_size bigint NOT NULL,
    container_height numeric NOT NULL,
    container_type text NOT NULL,
    start_slot bigint NOT NULL,
    end_slot bigint NOT NULL,
    start_row bigint NOT NULL,
    end_row bigint NOT NULL,
    priority_direction text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.yard_plans OWNER TO postgres;

--
-- Name: yard_plans_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.yard_plans_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.yard_plans_id_seq OWNER TO postgres;

--
-- Name: yard_plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.yard_plans_id_seq OWNED BY public.yard_plans.id;


--
-- Name: yards; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.yards (
    id bigint NOT NULL,
    name text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.yards OWNER TO postgres;

--
-- Name: yards_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.yards_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.yards_id_seq OWNER TO postgres;

--
-- Name: yards_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.yards_id_seq OWNED BY public.yards.id;


--
-- Name: blocks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blocks ALTER COLUMN id SET DEFAULT nextval('public.blocks_id_seq'::regclass);


--
-- Name: containers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.containers ALTER COLUMN id SET DEFAULT nextval('public.containers_id_seq'::regclass);


--
-- Name: yard_plans id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.yard_plans ALTER COLUMN id SET DEFAULT nextval('public.yard_plans_id_seq'::regclass);


--
-- Name: yards id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.yards ALTER COLUMN id SET DEFAULT nextval('public.yards_id_seq'::regclass);


--
-- Data for Name: blocks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.blocks (id, yard_id, name, max_slot, max_row, max_tier, created_at, updated_at) FROM stdin;
1	1	LC01	10	5	4	2025-11-01 19:47:34.999758+07	2025-11-01 19:47:34.999758+07
\.


--
-- Data for Name: containers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.containers (id, container_number, block_id, container_size, container_height, container_type, slot, "row", tier, is_placed, placed_at, picked_up_at, created_at, updated_at) FROM stdin;
1	ALFI000001	1	0	0		1	1	1	f	2025-11-01 20:11:48.1577+07	2025-11-01 20:16:54.521277+07	2025-11-01 20:11:48.157816+07	2025-11-01 20:16:54.521344+07
\.


--
-- Data for Name: yard_plans; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.yard_plans (id, block_id, container_size, container_height, container_type, start_slot, end_slot, start_row, end_row, priority_direction, created_at, updated_at) FROM stdin;
1	1	20	8.6	DRY	1	3	1	5	LEFT_TO_RIGHT	2025-11-01 19:47:35.00305+07	2025-11-01 19:47:35.00305+07
2	1	40	8.6	DRY	4	7	1	5	LEFT_TO_RIGHT	2025-11-01 19:47:35.006774+07	2025-11-01 19:47:35.006774+07
3	1	20	9.6	DRY	8	10	1	3	LEFT_TO_RIGHT	2025-11-01 19:47:35.009163+07	2025-11-01 19:47:35.009163+07
\.


--
-- Data for Name: yards; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.yards (id, name, created_at, updated_at) FROM stdin;
1	YRD1	2025-11-01 19:47:34.996929+07	2025-11-01 19:47:34.996929+07
\.


--
-- Name: blocks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.blocks_id_seq', 1, true);


--
-- Name: containers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.containers_id_seq', 1, true);


--
-- Name: yard_plans_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.yard_plans_id_seq', 3, true);


--
-- Name: yards_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.yards_id_seq', 1, true);


--
-- Name: blocks blocks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT blocks_pkey PRIMARY KEY (id);


--
-- Name: containers containers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.containers
    ADD CONSTRAINT containers_pkey PRIMARY KEY (id);


--
-- Name: yard_plans yard_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.yard_plans
    ADD CONSTRAINT yard_plans_pkey PRIMARY KEY (id);


--
-- Name: yards yards_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.yards
    ADD CONSTRAINT yards_pkey PRIMARY KEY (id);


--
-- Name: idx_containers_container_number; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_containers_container_number ON public.containers USING btree (container_number);


--
-- Name: idx_yards_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_yards_name ON public.yards USING btree (name);


--
-- Name: containers fk_blocks_containers; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.containers
    ADD CONSTRAINT fk_blocks_containers FOREIGN KEY (block_id) REFERENCES public.blocks(id);


--
-- Name: yard_plans fk_blocks_plans; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.yard_plans
    ADD CONSTRAINT fk_blocks_plans FOREIGN KEY (block_id) REFERENCES public.blocks(id);


--
-- Name: blocks fk_yards_blocks; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT fk_yards_blocks FOREIGN KEY (yard_id) REFERENCES public.yards(id);


--
-- PostgreSQL database dump complete
--

\unrestrict BsFZ5PoIeozMahLov2DIk8ubMv4ln2ez0bnq1BCX29seTtBYe4E9njnIsFMNfJE

