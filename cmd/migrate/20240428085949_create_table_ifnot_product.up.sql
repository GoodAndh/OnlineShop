create table if not exists product(
    id int not null auto_increment,
    name varchar(255) not null,
    deskripsi text not null,
    category enum ("electric","consumable","etc") not null default "etc",
    quantity int not null,
    price float not null,
    userid int not null ,
    primary key(id),
    foreign key(userid) references users(id)
)engine=innodb;