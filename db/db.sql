create table if not exists Users(
    id serial primary key,
    nickname varchar unique not null,
    fullname varchar not null,
    about text,
    email varchar(256) unique not null
);

create table if not exists Forum(
    id serial primary key,
    title varchar unique not null,
    slug varchar unique not null,
    post_num integer default 0 check (post_num >= 0),
    thread_num integer default 0 check (thread_num >= 0)
);

create table if not exists Thread(
    id serial primary key,
    forum_id integer references Forum,
    author_id integer references Users,
    title varchar not null,
    message varchar not null,
    created_at timestamp default now()
);

create table if not exists Post(
    id serial primary key,
    thread_id integer references Thread,
    author_id integer references Users,
    parent_id integer references Post,
    message varchar not null,
    rating integer default 0,
    created_at timestamp default now()
);

create table if not exists Vote(
    author_id integer references Users,
    thread_id integer references Thread,
    primary key(thread_id, user_id)
);
