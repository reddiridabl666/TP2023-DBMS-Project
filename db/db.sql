create table if not exists Users(
    id serial primary key,
    nickname varchar not null,
    fullname varchar not null,
    about text,
    email varchar(256) not null
);

create table if not exists Forum(
    id serial primary key,
    author_id integer references Users not null,
    title varchar unique not null,
    slug varchar unique not null,
    post_num integer default 0 check (post_num >= 0),
    thread_num integer default 0 check (thread_num >= 0)
);

create table if not exists Thread(
    id serial primary key,
    forum_id integer references Forum not null,
    author_id integer references Users not null,
    title varchar not null,
    message varchar not null,
    rating integer default 0,
    created_at timestamp default now()
);

create table if not exists Post(
    id serial primary key,
    thread_id integer references Thread not null,
    author_id integer references Users not null,
    parent_id integer references Post,
    message varchar not null,
    is_edited boolean default false,
    created_at timestamp default now()
);

create table if not exists Vote(
    author_id integer references Users,
    thread_id integer references Thread,
    value smallint check (value = 1 or value = -1),
    primary key(thread_id, author_id)
);

create unique index on users (lower(nickname));

create unique index on users (lower(email));

create or replace function update_vote_count() returns trigger as $$
    begin
        if (tg_op = 'INSERT') then
            update Thread set votes = votes + NEW.value where id = NEW.thread_id;
            return NEW;
        elsif (tg_op = "UPDATE") then
            update Thread set votes = votes - OLD.value + NEW.value where id = NEW.thread_id;
            return NEW;
        elsif (tg_op = "DELETE") then
            update Thread set votes = votes - OLD.value where id = OLD.thread_id;
            return OLD;
        end if;
        return NULL;
    end;
$$ language plpgsql;

create or replace function update_post_count() returns trigger as $$
    begin
        if (tg_op = 'INSERT') then
            update Forum set post_num = post_num + 1 where id = (select forum_id from Thread where id = NEW.thread_id);
            return NEW;
        elsif (tg_op = "DELETE") then
            update Forum set post_num = post_num - 1 where id = (select forum_id from Thread where id = OLD.thread_id);
            return OLD;
        end if;
        return NULL;
    end;
$$ language plpgsql;

create or replace function update_thread_count() returns trigger as $$
    begin
        if (tg_op = 'INSERT') then
            update Forum set thread_num = thread_num + 1 where id = NEW.forum_id;
            return NEW;
        elsif (tg_op = "DELETE") then
            update Forum set thread_num = thread_num - 1 where id = OLD.forum_id;
            return OLD;
        end if;
        return NULL;
    end;
$$ language plpgsql;

create trigger on_vote
after insert or update or delete on Vote
    for each row execute procedure update_vote_count();

create trigger on_post
after insert or delete on Post
    for each row execute procedure update_post_count();

create trigger on_thread
after insert or delete on Thread
    for each row execute procedure update_thread_count();
