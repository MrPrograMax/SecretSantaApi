CREATE TABLE groups
(
    id          serial       not null unique,
    name        varchar(255) not null,
    description varchar(255)
);

CREATE TABLE participants
(
    id           serial       not null unique,
    name         varchar(255) not null,
    wish         varchar(255),
    recipient_id int          references participants (id) on delete set null unique
);

CREATE TABLE groups_list
(
    id             serial                                                 not null unique,
    group_id       int references groups (id) on delete cascade           not null,
    participant_id int references participants (id) on delete set default not null
);








