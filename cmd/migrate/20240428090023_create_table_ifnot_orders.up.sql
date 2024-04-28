create table if not exists orders(
    id int not null auto_increment,
    total int not null,
    status enum("pending","sukses","cancel") not null default "pending",
    userid int not null ,
    foreign key (userid) references users(id),
    primary key(id)
)engine=innodb;