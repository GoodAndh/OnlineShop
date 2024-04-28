create table if not exists users(
    id int not null auto_increment,
    name varchar(155) not null,
    password varchar(155) not null,
    email varchar(155) not null,
    primary key (id),
    unique key  (email)
)engine=innodb;