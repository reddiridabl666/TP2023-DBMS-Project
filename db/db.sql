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
    title varchar not null,
    slug varchar not null,
    post_num integer default 0 check (post_num >= 0),
    thread_num integer default 0 check (thread_num >= 0)
);

create table if not exists Thread(
    id serial primary key,
    forum_id integer references Forum not null,
    author_id integer references Users not null,
    title varchar not null,
    slug varchar,
    message varchar not null,
    rating integer default 0,
    created_at bigint not null
);

create table if not exists Post(
    id serial primary key,
    thread_id integer references Thread not null,
    author_id integer references Users not null,
    parent_id integer references Post,
    path integer[] not null,
    message varchar not null,
    is_edited boolean default false,
    created_at bigint not null
);

create table if not exists Vote(
    author_id integer references Users,
    thread_id integer references Thread,
    value smallint check (value = 1 or value = -1),
    primary key(thread_id, author_id)
);

create table if not exists User_Forum(
    user_id integer references Users,
    forum_id integer references Forum,
    primary key(user_id, forum_id)
);

create unique index on users (lower(nickname));

create unique index on users (lower(email));

create unique index on forum (lower(slug));

create unique index on thread (lower(slug));

create index on thread (forum_id);

create index on post (thread_id);

create index on post ((path[1]));

create index on post ((path[2:]));

create or replace function update_vote_count() returns trigger as $$
    begin
        if (tg_op = 'INSERT') then
            update Thread set rating = rating + NEW.value where id = NEW.thread_id;
            return NEW;
        elsif (tg_op = 'UPDATE') then
            update Thread set rating = rating - OLD.value + NEW.value where id = NEW.thread_id;
            return NEW;
        elsif (tg_op = 'DELETE') then
            update Thread set rating = rating - OLD.value where id = OLD.thread_id;
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
        elsif (tg_op = 'DELETE') then
            update Forum set thread_num = thread_num - 1 where id = OLD.forum_id;
            return OLD;
        end if;
        return NULL;
    end;
$$ language plpgsql;

create or replace function on_post_insertion() returns trigger as $$
    begin
        if (NEW.parent_id IS NULL) then
            NEW.path = array[NEW.id];
            return NEW;
        end if;

        if (NEW.thread_id != (SELECT thread_id FROM post WHERE id = NEW.parent_id)) then
            raise exception 'Thread id should match with parent`s' USING ERRCODE = '23000';
        end if;
        
        NEW.path = (SELECT path FROM post WHERE id = NEW.parent_id) || NEW.id;
        return NEW;
    end;
$$ language plpgsql;

create or replace function add_forum_link_from_thread() returns trigger as $$
    begin
        insert into user_forum(user_id, forum_id)
            values (NEW.author_id, NEW.forum_id)
            on conflict do nothing;
        return NEW;
    end;
$$ language plpgsql;

create trigger on_vote
after insert or update or delete on Vote
    for each row execute procedure update_vote_count();

create trigger on_thread
after insert or delete on Thread
    for each row execute procedure update_thread_count();

create trigger before_post_insert
before insert on Post
    for each row execute procedure on_post_insertion();

create trigger add_forum_link_on_thread
after insert on Thread
    for each row execute procedure add_forum_link_from_thread();
