--
-- PostgreSQL database dump
--

\restrict FzrQPY3KN2Q1ncKXsnaSQmz8oaNQEF8APTktPj63uTMmhm2NpxvoRERF4B2EyrP

-- Dumped from database version 15.15
-- Dumped by pg_dump version 15.15

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

--
-- Name: messages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.messages (
    id bigint NOT NULL,
    room_id character varying(255) NOT NULL,
    sender_id character varying(255) NOT NULL,
    sender_nickname character varying(50) NOT NULL,
    avatar character varying(50) NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.messages OWNER TO postgres;

--
-- Name: messages_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.messages_id_seq OWNER TO postgres;

--
-- Name: messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.messages_id_seq OWNED BY public.messages.id;


--
-- Name: messages id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.messages ALTER COLUMN id SET DEFAULT nextval('public.messages_id_seq'::regclass);


--
-- Data for Name: messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.messages (id, room_id, sender_id, sender_nickname, avatar, content, created_at) FROM stdin;
1	global	56c73fb4-5a92-4d7c-a952-3841a379ace5	고요한 사슴	avatar-13	안녕	2025-11-25 01:03:21.576598+00
2	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	안녕하세요	2025-11-25 01:03:52.357843+00
3	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	반갑습니다~	2025-11-25 01:04:00.016179+00
4	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	안녕하세요~~!	2025-11-25 01:04:19.583159+00
5	global	56c73fb4-5a92-4d7c-a952-3841a379ace5	고요한 사슴	avatar-13	누구세요	2025-11-25 01:05:10.531609+00
6	global	9f7a6837-e3d5-4956-b791-f3fcc3109b2b	성실한 강아지	avatar-10	안녕하세요~	2025-11-25 01:05:28.647968+00
7	global	9f7a6837-e3d5-4956-b791-f3fcc3109b2b	성실한 강아지	avatar-10	웹소켓 연결이 잘되나요?	2025-11-25 01:08:59.127918+00
8	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	ㅇㅇ	2025-11-25 01:12:25.165858+00
9	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	너 뭐야	2025-11-25 01:15:02.015852+00
10	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	넌 뭔데	2025-11-25 01:15:05.268999+00
11	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	나 ㅇㅈ이	2025-11-25 01:15:10.502785+00
12	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	..	2025-11-25 01:15:20.034663+00
13	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	넌 누군데	2025-11-25 01:15:20.986007+00
14	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	나 ㄷㅈ이	2025-11-25 01:15:27.171124+00
15	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	ㅌㅌ	2025-11-25 01:16:12.798791+00
16	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	ㅋㅋ	2025-11-25 01:17:22.683988+00
17	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	ㅋㅋ\\	2025-11-25 01:17:25.624197+00
18	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	몇초 지나면	2025-11-25 01:17:33.022969+00
19	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	자동으로	2025-11-25 01:17:34.205715+00
20	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	종료되나보다	2025-11-25 01:17:36.438531+00
21	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	그런가	2025-11-25 01:17:38.67858+00
22	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	왜그러냐	2025-11-25 01:17:41.375654+00
23	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	??	2025-11-25 01:17:43.623936+00
24	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	왜그러냐고	2025-11-25 01:17:44.817375+00
25	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	시작	2025-11-25 01:18:37.18729+00
26	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	3초	2025-11-25 01:18:41.672852+00
27	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	5초	2025-11-25 01:18:48.658992+00
28	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	7초	2025-11-25 01:18:59.165953+00
29	global	775fca69-8f50-43a0-9d61-fb628792e129	성실한 강아지	avatar-10	오 준규 쩐다	2025-11-25 01:20:38.874302+00
30	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	1초	2025-11-25 01:33:15.582627+00
31	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	2초	2025-11-25 01:33:17.332149+00
32	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	3초	2025-11-25 01:33:18.760489+00
33	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	4초	2025-11-25 01:33:19.960959+00
34	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	5초	2025-11-25 01:33:21.94253+00
35	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	6초	2025-11-25 01:33:23.533999+00
36	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	7초	2025-11-25 01:33:24.660908+00
37	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	8초	2025-11-25 01:33:26.326967+00
38	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	9초	2025-11-25 01:33:28.118969+00
39	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	10초	2025-11-25 01:33:29.213724+00
40	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	11초	2025-11-25 01:33:30.193807+00
41	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	12초	2025-11-25 01:33:31.111097+00
42	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	13초	2025-11-25 01:33:32.060035+00
43	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	15초	2025-11-25 01:34:44.158508+00
44	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	1초	2025-11-25 01:34:47.2401+00
45	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	2초	2025-11-25 01:34:48.23144+00
46	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	3초	2025-11-25 01:34:49.091452+00
47	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	4초	2025-11-25 01:34:50.55857+00
48	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	5초	2025-11-25 01:34:51.494128+00
49	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	6초	2025-11-25 01:34:52.609243+00
50	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	7초	2025-11-25 01:34:53.687598+00
51	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	8	2025-11-25 01:34:54.833359+00
52	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	9	2025-11-25 01:34:55.538376+00
53	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	10	2025-11-25 01:34:57.240737+00
54	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	111	2025-11-25 01:34:59.240676+00
55	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	12	2025-11-25 01:35:01.240886+00
56	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	13	2025-11-25 01:35:03.24053+00
57	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	14	2025-11-25 01:35:05.240978+00
58	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	15	2025-11-25 01:35:07.241129+00
59	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	16	2025-11-25 01:35:09.241321+00
60	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	----	2025-11-25 01:57:04.447082+00
61	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	1	2025-11-25 01:57:05.278032+00
62	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	2	2025-11-25 01:57:05.910997+00
63	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	3	2025-11-25 01:57:06.475976+00
64	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	4	2025-11-25 01:57:07.283864+00
65	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	5	2025-11-25 01:57:07.798928+00
66	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	6	2025-11-25 01:57:08.446468+00
67	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	7	2025-11-25 01:57:10.447948+00
68	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	8	2025-11-25 01:57:12.447758+00
69	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	9	2025-11-25 01:57:14.448092+00
70	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	10	2025-11-25 01:57:16.447426+00
71	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	11	2025-11-25 01:57:18.447164+00
72	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	12	2025-11-25 01:57:20.447859+00
73	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	13	2025-11-25 01:57:22.446978+00
74	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	14	2025-11-25 01:57:24.447348+00
75	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	15	2025-11-25 01:57:26.781868+00
76	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	16	2025-11-25 01:57:28.446944+00
77	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	17	2025-11-25 01:57:30.447329+00
78	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	18	2025-11-25 01:57:32.447423+00
79	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	19	2025-11-25 01:57:34.44788+00
80	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	20	2025-11-25 01:57:36.447699+00
81	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	21	2025-11-25 01:57:38.447871+00
82	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	22	2025-11-25 01:57:40.447323+00
83	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	23	2025-11-25 01:57:45.487064+00
84	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	24	2025-11-25 01:57:48.90881+00
85	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	25	2025-11-25 01:57:50.496237+00
86	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	asdf	2025-11-25 01:57:55.339791+00
87	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	zz	2025-11-25 01:57:57.009086+00
88	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	asdf	2025-11-25 01:59:13.624109+00
89	global	b99475a1-c50e-451e-a19b-9a2e7920f39c	배고픈 고슴도치	avatar-12	해결완료 ㅅㄱㅇ	2025-11-25 01:59:17.973281+00
90	global	775fca69-8f50-43a0-9d61-fb628792e129	성실한 강아지	avatar-10	안녕하세요	2025-11-25 02:12:46.851646+00
91	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	왜캐 시끄럽냐	2025-11-25 02:17:11.54355+00
92	global	ed725ff3-ec35-4d3a-867f-917255e8308c	유쾌한 고양이	avatar-7	어케 해결햇냐	2025-11-25 02:17:42.126015+00
93	global	586465c7-29de-4b82-9597-47c06defab9a	성실한 돌고래	avatar-5	유쾌한 고양이님 안녕하세요	2025-11-25 02:17:51.413156+00
94	global	9f7a6837-e3d5-4956-b791-f3fcc3109b2b	성실한 강아지	avatar-10	ㄴㄴ	2025-11-25 03:39:53.371963+00
95	global	1981ab21-ac12-42b8-836d-12ac161df10f	재빠른 기린	avatar-14	안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까안녕하십니까	2025-11-27 14:30:33.491258+00
96	global	c1b71e01-7adc-4bd5-bb0f-ae803da75fb7	총명한 강아지	avatar-10	ㅁㅇㄴㄹ	2025-11-27 21:20:51.35527+00
97	global	d580d3fc-95da-4f66-9546-30dab7aa27a3	총명한 다람쥐	avatar-16	안녕하세요	2025-12-02 07:09:30.788004+00
98	global	08ea2964-beb3-4a01-89bd-967f98b09ef2	친절한 고슴도치	avatar-12	Hello	2025-12-04 06:11:13.703457+00
\.


--
-- Name: messages_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.messages_id_seq', 98, true);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: idx_messages_room_created; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_messages_room_created ON public.messages USING btree (room_id, created_at DESC);


--
-- PostgreSQL database dump complete
--

\unrestrict FzrQPY3KN2Q1ncKXsnaSQmz8oaNQEF8APTktPj63uTMmhm2NpxvoRERF4B2EyrP

