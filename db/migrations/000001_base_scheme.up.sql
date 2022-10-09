create table users (
                       id serial,
                       primary key(id),
                       login text not null unique,
                       first_name text not null,
                       last_name text not null
);

create table events (
                        id bigserial,
                        primary key(id),
                        name text not null,
                        author int not null,
                        constraint fk_author foreign key(author) references users,
                        repeat_options text,
                        repeatable boolean,
                        begin_time timestamp not null,
                        duration int,
                        end_time timestamp not null,
                        is_private bool not null,
                        details text
);

create type invitation_status as enum (
    'accepted',
    'declined',
    'not_answered'
);

create table users_events (
                              id bigserial,
                              primary key(id),
                              user_id int not null,
                              event_id bigint not null,
                              status invitation_status,
                              constraint fk_user foreign key(user_id) references users,
                              constraint fk_event foreign key(event_id) references events
);

create index on users_events(user_id);
